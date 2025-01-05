package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kitproj/kit/internal"
	"github.com/kitproj/kit/internal/proc"
	"github.com/kitproj/kit/internal/types"
	"github.com/kitproj/kit/internal/util"
	"k8s.io/utils/strings/slices"
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

		pod := &types.Pod{}

		in, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}
		if err = yaml.UnmarshalStrict(in, pod); err != nil {
			return err
		}

		if rewrite {
			out, err := yaml.Marshal(pod)
			if err != nil {
				return err
			}
			return os.WriteFile(configFile, out, 0644)
		}

		dag := internal.NewDAG[bool]()
		for name, t := range pod.Tasks {
			dag.AddNode(name, true)
			for _, dependency := range t.Dependencies {
				dag.AddEdge(dependency, name)
			}
		}
		visited := dag.Subgraph(args)

		taskByName := pod.Tasks
		subgraph := internal.NewDAG[*taskNode]()
		for name := range visited {
			task := taskByName[name]
			subgraph.AddNode(name, &taskNode{task: task})
			for _, parent := range dag.Parents[name] {
				subgraph.AddEdge(parent, name)
			}
		}

		taskNames := make(chan string, len(subgraph.Nodes))

		// schedule the tasks in the subgraph that are ready to run , this is done by sending the task name to the events channel of any task that does not have any parents
		for taskName := range subgraph.Nodes {
			if len(subgraph.Parents[taskName]) == 0 {
				taskNames <- taskName
			}
		}

		// start a file watcher for each task
		for _, t := range subgraph.Nodes {

			// start watching files for changes
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}
			for _, source := range t.task.Watch {
				if err := watcher.Add(filepath.Join(t.task.WorkingDir, source)); err != nil {
					return err
				}
			}
			defer watcher.Close()

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case event := <-watcher.Events:
						if event.Op&fsnotify.Write == fsnotify.Write {
							log.Printf("file changed, re-running %s\n", t.name)
							taskNames <- t.name
						}
					}
				}
			}()
		}

		semaphores := util.NewSemaphores(pod.Semaphores)

		wg := sync.WaitGroup{}

		for {
			select {
			case <-ctx.Done():

				wg.Wait()

				// if any task failed, we will return an error
				var failures []string
				for _, node := range subgraph.Nodes {
					if node.phase == "failed" {
						failures = append(failures, fmt.Sprintf("%s %s", node.name, node.message))
					}
				}

				if len(failures) > 0 {
					return fmt.Errorf("failed tasks: %v", failures)
				}

				return nil
			case taskName := <-taskNames:

				// if we get the poison pill, we should see if all tasks are done and exit if so
				const PoisonPill = ""
				if taskName == PoisonPill {
					busy := false
					for _, node := range subgraph.Nodes {
						if node.busy() {
							busy = true
						}
					}
					if !busy {
						cancel()
					}
					continue
				}

				// we will only execute this task, if its parents are "succeeded" or "skipped" or ("running" and the task is a service)
				blocked := false
				for _, parentName := range subgraph.Parents[taskName] {
					parent := subgraph.Nodes[parentName]
					if parent.blocked() {
						log.Printf("task %q is blocked by %q (%s) %s\n", taskName, parentName, parent.phase, parent.message)
						blocked = true
					}
				}

				if blocked {
					continue
				}

				// we might already be waiting, starting or running this task, so we don't want to start it again
				node := subgraph.Nodes[taskName]
				if node.busy() {
					continue
				}

				// each task is executed in a separate goroutine
				wg.Add(1)

				node.phase = "waiting"
				node.message = ""

				go func(t types.Task) {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					defer func() { taskNames <- PoisonPill }()
					defer wg.Done()

					out := funcWriter(func(bytes []byte) (int, error) {
						prefix := fmt.Sprintf("%s[%s] (%s) ", internal.Color(node.name, t.IsService()), node.name, node.phase)
						// reset color and bold
						suffix := "\033[0m"

						// split on newlines
						lines := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
						for _, line := range lines {
							log.Println(prefix + line + suffix)
						}

						return len(bytes), nil
					})

					log := log.New(out, "", 0)

					queueChildren := func() {
						for _, child := range subgraph.Children[node.name] {
							// only queue tasks in the subgraph
							if _, ok := subgraph.Nodes[child]; ok {
								log.Printf("queuing %q\n", child)
								taskNames <- child
							}
						}
					}

					// if the task can be skipped, lets exit early
					if t.Skip() || slices.Contains(strings.Split(tasksToSkip, ","), node.name) {
						node.phase = "succeeded"
						node.message = "skipped"
						log.Println("skipping")
						queueChildren()
						return
					}

					// if the task needs a mutex, lets wait for it
					if t.Mutex != "" {
						mu := util.GetMutex(t.Mutex)
						node.phase = "waiting"
						node.message = "waiting for mutex"
						log.Println("waiting for mutex")
						mu.Lock()
						node.message = "acquired mutex"
						log.Println("acquired mutex")
						defer mu.Unlock()
					}

					// if the task needs a semaphore, lets wait for it
					if t.Semaphore != "" {
						sema := semaphores.Get(t.Semaphore)
						node.phase = "waiting"
						node.message = "waiting for semaphore"
						log.Println("waiting for semaphore")
						if err := sema.Acquire(ctx, 1); err != nil {
							node.phase = "failed"
							node.message = fmt.Sprintf("failed to acquire semaphore: %v", err)
							log.Printf("failed to acquire semaphore: %v", err)
							return
						}
						node.message = "acquired semaphore"
						defer sema.Release(1)
					}

					p := proc.New(t, log, types.PodSpec(*pod))

					if probe := t.GetLivenessProbe(); probe != nil {
						liveFunc := func(live bool, err error) {
							if !live {
								node.phase = "failed"
								node.message = fmt.Sprintf("liveness probe failed: %v", err)
								log.Printf("liveness probe failed: %v", err)
								cancel()
							}
						}
						go probeLoop(ctx, *probe, liveFunc)
					}
					if probe := t.GetReadinessProbe(); probe != nil {
						readyFunc := func(ready bool, err error) {
							if ready {
								node.phase = "running"
								node.message = fmt.Sprintf("readiness probe succeeded")
								log.Println("readiness probe succeeded")
								queueChildren()
							} else {
								node.phase = "failed"
								node.message = fmt.Sprintf("readiness probe failed: %v", err)
								log.Println("readiness probe failed")
								cancel()
							}
						}
						go probeLoop(ctx, *probe, readyFunc)
					}

					if t.IsService() {
						node.phase = "starting"
					} else {
						node.phase = "running"
					}

					restart := func() {
						select {
						case <-ctx.Done():
						case <-time.After(3 * time.Second):
							log.Println("restarting")
							taskNames <- node.name
						}
					}

					if t.Log != "" {
						out, err := os.Create(t.Log)
						if err != nil {
							node.phase = "failed"
							node.message = fmt.Sprintf("failed to create log file: %v", err)
							log.Printf("failed to create log file: %v", err)
							return
						}
						defer out.Close()
					}

					err = p.Run(ctx, out, out)
					if err != nil {
						node.phase = "failed"
						node.message = fmt.Sprint(err)
						log.Println(err)
						if t.GetRestartPolicy() != "Never" {
							restart()
						}
						return
					}

					node.phase = "succeeded"
					log.Println("succeeded")
					if t.GetRestartPolicy() == "Always" {
						restart()
					}
					queueChildren()

				}(taskByName[taskName])
			}
		}
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
