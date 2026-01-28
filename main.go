package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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

func main() {
	help := false
	printVersion := false
	workingDir := "."
	configFile := ""
	tasksToSkip := ""
	port := -1 // -1 means unspecified, 0 means disabled, >0 means specified
	openBrowser := false
	rewrite := false
	completion := ""

	flag.BoolVar(&help, "h", false, "print help and exit")
	flag.BoolVar(&printVersion, "v", false, "print version and exit")
	flag.StringVar(&workingDir, "C", ".", "working directory (default current directory)")
	flag.StringVar(&configFile, "f", "tasks.yaml", "config file (default tasks.yaml)")
	flag.StringVar(&tasksToSkip, "s", "", "tasks to skip (comma separated)")
	flag.IntVar(&port, "p", port, "port to start UI on (default 3000, zero disables)")
	flag.BoolVar(&openBrowser, "b", false, "open the UI in the browser (default false)")
	flag.BoolVar(&rewrite, "w", false, "rewrite the config file")
	flag.StringVar(&completion, "completion", "", "generate shell completion script (bash, zsh, fish)")
	flag.Parse()
	taskNames := flag.Args()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if printVersion {
		info, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", info.Main.Version)
		os.Exit(0)
	}

	if completion != "" {
		if err := printCompletion(completion, configFile); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := os.Chdir(workingDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to change to directory %s: %v\n", workingDir, err)
		os.Exit(1)
	}

	err := func() error {

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		wf := &types.Workflow{}

		in, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", configFile, err)
		}
		if err = yaml.UnmarshalStrict(in, wf); err != nil {
			return fmt.Errorf("failed to parse %s: %w", configFile, err)
		}

		if rewrite {
			out, err := yaml.Marshal(wf)
			if err != nil {
				return fmt.Errorf("failed to marshal %s: %w", configFile, err)
			}
			return os.WriteFile(configFile, out, 0644)
		}

		// if wf.Port is specified, use that, unless the user has specified a port on the command line
		if port == -1 {
			if wf.Port != nil {
				port = int(*wf.Port)
			} else {
				port = 3000 // default port
			}
		}

		// split the tasks on comma, but don't end up with a single entry of ""
		split := strings.Split(tasksToSkip, ",")
		if len(split) == 1 && split[0] == "" {
			split = []string{}
		}

		if len(taskNames) == 0 {
			for taskName, task := range wf.Tasks {
				if task.Default {
					taskNames = []string{taskName}
					break
				}
			}
		}

		return internal.RunSubgraph(
			ctx,
			cancel,
			port,
			openBrowser,
			log.Default(),
			wf,
			taskNames,
			split,
		)
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
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
