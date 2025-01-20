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
	configFile := ""
	tasksToSkip := ""
	port := 0
	openBrowser := false
	rewrite := false

	flag.BoolVar(&help, "h", false, "print help and exit")
	flag.BoolVar(&printVersion, "v", false, "print version and exit")
	flag.StringVar(&configFile, "f", "tasks.yaml", "config file (default tasks.yaml)")
	flag.StringVar(&tasksToSkip, "s", "", "tasks to skip (comma separated)")
	flag.IntVar(&port, "p", 3000, "port to start UI on (default 3000, zero disables)")
	flag.BoolVar(&openBrowser, "b", false, "open the UI in the browser (default false)")
	flag.BoolVar(&rewrite, "w", false, "rewrite the config file")
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

	err := func() error {

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		wf := &types.Workflow{}

		in, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}
		if err = yaml.UnmarshalStrict(in, wf); err != nil {
			return err
		}

		if rewrite {
			out, err := yaml.Marshal(wf)
			if err != nil {
				return err
			}
			return os.WriteFile(configFile, out, 0644)
		}

		// split the tasks on comma, but don't end up with a single entry of ""
		split := strings.Split(tasksToSkip, ",")
		if len(split) == 1 && split[0] == "" {
			split = []string{}
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
