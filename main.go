package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kitproj/kit/internal"
	"github.com/kitproj/kit/internal/types"
	"sigs.k8s.io/yaml"
)

//go:generate sh -c "git describe --tags > tag"
//go:embed tag
var tag string

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	help := false
	printVersion := false
	configFile := ""
	tasksToSkip := ""
	rewrite := false

	flag.BoolVar(&help, "h", false, "print help and exit")
	flag.BoolVar(&printVersion, "v", false, "print version and exit")
	flag.StringVar(&configFile, "f", "tasks.yaml", "config file")
	flag.StringVar(&tasksToSkip, "s", "", "tasks to skip (comma separated)")
	flag.BoolVar(&rewrite, "w", false, "rewrite the config file")
	flag.Parse()
	taskNames := flag.Args()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(tag)
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

		return internal.RunSubgraph(
			ctx,
			cancel,
			log.Default(),
			wf,
			taskNames,
			strings.Split(tasksToSkip, ","))
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
