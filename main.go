package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/kitproj/kit/internal"
	"github.com/kitproj/kit/internal/types"
	"sigs.k8s.io/yaml"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

const kitDescription = "workflow engine for software development"

type options struct {
	help           bool
	printVersion   bool
	workingDir     string
	configFile     string
	configExplicit bool
	tasksToSkip    string
	port           int
	openBrowser    bool
	rewrite        bool
	completion     string
	taskNames      []string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if err := execute(args, stdout); err != nil {
		_, _ = fmt.Fprintln(stderr, err.Error())
		return 1
	}
	return 0
}

func execute(args []string, stdout io.Writer) error {
	opts, flagSet, err := parseOptions(args)
	if err != nil {
		return explainFlagError(err)
	}

	if opts.help {
		printUsage(flagSet, stdout)
		return nil
	}

	if opts.printVersion {
		info, _ := debug.ReadBuildInfo()
		_, _ = fmt.Fprintf(stdout, "%v\n", info.Main.Version)
		return nil
	}

	if opts.completion != "" {
		if err := printCompletion(opts.completion, resolveConfigFile(opts.configFile)); err != nil {
			return explainCompletionError(err)
		}
		return nil
	}

	previousDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(opts.workingDir); err != nil {
		return explainWorkingDirError(opts.workingDir, err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()

	configFile := resolveConfigFile(opts.configFile)
	configPath, err := filepath.Abs(configFile)
	if err != nil {
		configPath = configFile
	}
	configSource := configFileSource(opts.configExplicit)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wf := &types.Workflow{}

	in, err := os.ReadFile(configFile)
	if err != nil {
		return explainConfigReadError(configPath, configSource, err)
	}
	if err = yaml.UnmarshalStrict(in, wf); err != nil {
		return explainConfigParseError(configPath, configSource, err)
	}

	printStartup(stdout)
	printConfig(stdout, configPath, configSource)

	if opts.rewrite {
		out, err := yaml.Marshal(wf)
		if err != nil {
			return explainConfigWriteError(configPath, err)
		}
		if err := os.WriteFile(configFile, out, 0644); err != nil {
			return explainConfigWriteError(configPath, err)
		}
		return nil
	}

	port := opts.port
	if port == -1 {
		if wf.Port != nil {
			port = int(*wf.Port)
		} else {
			port = 3000 // default port
		}
	}

	// split the tasks on comma, but don't end up with a single entry of ""
	split := strings.Split(opts.tasksToSkip, ",")
	if len(split) == 1 && split[0] == "" {
		split = []string{}
	}

	taskNames := opts.taskNames
	if len(taskNames) == 0 {
		for taskName, task := range wf.Tasks {
			if task.Default {
				taskNames = []string{taskName}
				break
			}
		}
	}

	logger := log.New(stdout, "", 0)
	return explainRuntimeError(internal.RunSubgraph(
		ctx,
		cancel,
		port,
		opts.openBrowser,
		logger,
		wf,
		taskNames,
		split,
	), configPath)
}

func parseOptions(args []string) (*options, *flag.FlagSet, error) {
	opts := &options{workingDir: ".", port: -1}
	flagSet := flag.NewFlagSet("kit", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar(&opts.help, "h", false, "print help and exit")
	flagSet.BoolVar(&opts.printVersion, "v", false, "print version and exit")
	flagSet.StringVar(&opts.workingDir, "C", ".", "working directory (default current directory)")
	flagSet.StringVar(&opts.configFile, "f", "", "config file (default tasks.yaml)")
	flagSet.StringVar(&opts.tasksToSkip, "s", "", "tasks to skip (comma separated)")
	flagSet.IntVar(&opts.port, "p", opts.port, "port to start UI on (default 3000, zero disables)")
	flagSet.BoolVar(&opts.openBrowser, "b", false, "open the UI in the browser (default false)")
	flagSet.BoolVar(&opts.rewrite, "w", false, "rewrite the config file")
	flagSet.StringVar(&opts.completion, "completion", "", "generate shell completion script (bash, zsh, fish)")
	if err := flagSet.Parse(args); err != nil {
		return nil, flagSet, err
	}
	opts.taskNames = flagSet.Args()
	flagSet.Visit(func(f *flag.Flag) {
		if f.Name == "f" {
			opts.configExplicit = true
		}
	})
	return opts, flagSet, nil
}

func printUsage(flagSet *flag.FlagSet, stdout io.Writer) {
	previousOutput := flagSet.Output()
	flagSet.SetOutput(stdout)
	defer flagSet.SetOutput(previousOutput)
	flagSet.Usage()
}

func resolveConfigFile(configFile string) string {
	if configFile == "" {
		return "tasks.yaml"
	}
	return configFile
}

func configFileSource(explicit bool) string {
	if explicit {
		return "explicit"
	}
	return "default"
}

func printStartup(stdout io.Writer) {
	if version := buildVersion(); version != "" {
		_, _ = fmt.Fprintf(stdout, "kit: startup: %s; version=%s\n", kitDescription, version)
		return
	}
	_, _ = fmt.Fprintf(stdout, "kit: startup: %s\n", kitDescription)
}

func printConfig(stdout io.Writer, configPath, configSource string) {
	_, _ = fmt.Fprintf(stdout, "kit: config: path=%s source=%s\n", configPath, configSource)
}

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}

func explainFlagError(err error) error {
	return explainError(
		fmt.Sprintf("CLI argument parsing failed: %v", err),
		"an unsupported flag or invalid flag value was provided",
		"run `kit -h` to see the supported flags and usage",
	)
}

func explainCompletionError(err error) error {
	return explainError(
		fmt.Sprintf("completion generation failed: %v", err),
		"the requested shell is not supported",
		"pass `--completion bash`, `--completion zsh`, or `--completion fish`",
	)
}

func explainWorkingDirError(path string, err error) error {
	cause := "the directory does not exist or is not accessible"
	if errors.Is(err, os.ErrPermission) {
		cause = "kit does not have permission to access the directory"
	}
	return explainError(
		fmt.Sprintf("working directory change failed for %s: %v", path, err),
		cause,
		"pass a valid directory with `-C` or run kit from the project root",
	)
}

func explainConfigReadError(path, source string, err error) error {
	cause := "kit could not access the config file"
	next := fmt.Sprintf("verify %s and retry, or pass `-f /path/to/tasks.yaml`", path)
	if errors.Is(err, os.ErrNotExist) {
		cause = "the config file was not found"
		next = fmt.Sprintf("create %s or pass `-f /path/to/tasks.yaml`", path)
	} else if errors.Is(err, os.ErrPermission) {
		cause = "kit does not have permission to read the config file"
	}
	return explainError(
		fmt.Sprintf("config read failed for %s (source=%s): %v", path, source, err),
		cause,
		next,
	)
}

func explainConfigParseError(path, source string, err error) error {
	return explainError(
		fmt.Sprintf("config parse failed for %s (source=%s): %v", path, source, err),
		"the config file contains invalid YAML or unsupported fields",
		"fix the config file and retry",
	)
}

func explainConfigWriteError(path string, err error) error {
	return explainError(
		fmt.Sprintf("config rewrite failed for %s: %v", path, err),
		"kit could not serialize or write the config file",
		"check the file permissions and retry",
	)
}

func explainRuntimeError(err error, configPath string) error {
	if err == nil {
		return nil
	}
	message := err.Error()
	switch {
	case strings.HasPrefix(message, "task ") && strings.Contains(message, " not found in workflow"):
		return explainError(
			fmt.Sprintf("task selection failed: %s", message),
			"the requested task name is not defined in the loaded workflow",
			fmt.Sprintf("check the task names in %s and retry", configPath),
		)
	case strings.HasPrefix(message, "skipped task ") && strings.Contains(message, " not found in workflow"):
		return explainError(
			fmt.Sprintf("task skip selection failed: %s", message),
			"a task listed in `-s` is not defined in the loaded workflow",
			fmt.Sprintf("check the task names in %s and retry", configPath),
		)
	case strings.Contains(message, " is invalid:"):
		return explainError(
			fmt.Sprintf("workflow validation failed: %s", message),
			"one or more task definitions are invalid",
			fmt.Sprintf("fix the task definition in %s and retry", configPath),
		)
	case strings.HasPrefix(message, "failed tasks:"):
		return explainError(
			fmt.Sprintf("workflow run failed: %s", message),
			"one or more tasks exited with a non-zero status",
			"inspect the task output above or the logs/ directory and retry",
		)
	default:
		return explainError(
			fmt.Sprintf("workflow run failed: %s", message),
			"kit hit an unexpected runtime error while running the workflow",
			"inspect the task output above and retry",
		)
	}
}

func explainError(summary, cause, next string) error {
	return fmt.Errorf("kit: error: %s\nkit: cause: %s\nkit: next: %s", summary, cause, next)
}

func printCompletion(shell, configFile string) error {
	switch shell {
	case "bash":
		fmt.Print(bashCompletion(configFile))
	case "zsh":
		fmt.Print(zshCompletion(configFile))
	case "fish":
		fmt.Print(fishCompletion(configFile))
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
	}
	return nil
}

func bashCompletion(configFile string) string {
	return fmt.Sprintf(`_kit_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local tasks=""
    
    if [[ -f "%s" ]]; then
        tasks=$(grep -E '^  [a-zA-Z0-9_-]+:\s*$' "%s" 2>/dev/null | sed 's/:.*//' | tr -d ' ' | tr '\n' ' ')
    fi
    
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "-h -v -C -f -s -p -b -w --completion" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "${tasks}" -- "${cur}"))
    fi
}

complete -F _kit_completions kit
`, configFile, configFile)
}

func zshCompletion(configFile string) string {
	return fmt.Sprintf(`#compdef kit

_kit() {
    local -a tasks
    local -a opts
    
    opts=(
        '-h[print help and exit]'
        '-v[print version and exit]'
        '-C[working directory]:directory:_files -/'
        '-f[config file]:file:_files'
        '-s[tasks to skip]:tasks:'
        '-p[port to start UI on]:port:'
        '-b[open the UI in the browser]'
        '-w[rewrite the config file]'
        '--completion[generate shell completion script]:shell:(bash zsh fish)'
    )
    
    if [[ -f "%s" ]]; then
        tasks=(${(f)"$(grep -E '^  [a-zA-Z0-9_-]+:\s*$' "%s" 2>/dev/null | sed 's/:.*//' | tr -d ' ')"})
    fi
    
    _arguments -s $opts '*:task:'"(${tasks[*]})"
}

# don't run the completion function when being source-ed or eval-ed
if [ "$funcstack[1]" = "_kit" ]; then
    _kit
fi

# register the completion function (requires compinit to have been run)
if [[ -n ${_comps+x} ]]; then
    compdef _kit kit
fi
`, configFile, configFile)
}

func fishCompletion(configFile string) string {
	return fmt.Sprintf(`function __fish_kit_tasks
    if test -f "%s"
        grep -E '^  [a-zA-Z0-9_-]+:\s*$' "%s" 2>/dev/null | sed 's/:.*//' | string trim
    end
end

complete -c kit -f
complete -c kit -s h -d 'print help and exit'
complete -c kit -s v -d 'print version and exit'
complete -c kit -s C -d 'working directory' -r -a '(__fish_complete_directories)'
complete -c kit -s f -d 'config file' -r -F
complete -c kit -s s -d 'tasks to skip'
complete -c kit -s p -d 'port to start UI on'
complete -c kit -s b -d 'open the UI in the browser'
complete -c kit -s w -d 'rewrite the config file'
complete -c kit -l completion -d 'generate shell completion script' -r -a 'bash zsh fish'
complete -c kit -n 'not string match -q -- "-*" (commandline -ct)' -a '(__fish_kit_tasks)'
`, configFile, configFile)
}
