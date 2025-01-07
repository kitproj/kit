package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kitproj/kit/internal"
	"github.com/kitproj/kit/internal/types"
	"sigs.k8s.io/yaml"
)

//go:generate sh -c "git describe --tags > tag"
//go:embed tag
var tag string

// GitHub Actions
const defaultConfigFile = "tasks.yaml"

func init() {
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
	flag.StringVar(&configFile, "f", defaultConfigFile, "config file")
	flag.StringVar(&tasksToSkip, "s", "", "tasks to skip (comma separated)")
	flag.BoolVar(&rewrite, "w", false, "rewrite the config file")
	flag.Parse()
	args := flag.Args()

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

		dag := internal.NewDAG[bool]()
		for name, t := range wf.Tasks {
			dag.AddNode(name, true)
			for _, dependency := range t.Dependencies {
				dag.AddEdge(dependency, name)
			}
		}
		visited := dag.Subgraph(args)

		taskByName := wf.Tasks
		subgraph := internal.NewDAG[*internal.TaskNode]()
		for name := range visited {
			task := taskByName[name]
			subgraph.AddNode(name, &internal.TaskNode{Name: name, Task: task, Phase: "pending", Cancel: func() {}})
			for _, parent := range dag.Parents[name] {
				subgraph.AddEdge(parent, name)
			}
		}

		return internal.RunSubgraph(ctx, cancel, wf, subgraph, tasksToSkip)
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
