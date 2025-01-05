package main

import (
	"context"
	_ "embed"
	"errors"
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

	flag.BoolVar(&help, "h", false, "print help and exit")
	flag.BoolVar(&printVersion, "v", false, "print version and exit")
	flag.StringVar(&configFile, "f", defaultConfigFile, "config file")
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

		// make sure that the file is valid,
		// this helps Copilot to auto-complete the file,
		// no need to have any IDE plugin - welcome to the future
		if pod.ApiVersion != "kit/v1" {
			return errors.New("invalid apiVersion, must be 'kit/v1")
		}
		if pod.Kind != "Tasks" {
			return errors.New("invalid kind, must be 'Tasks'")
		}
		if pod.Metadata.Name == "" {
			return errors.New("metadata.name is required")
		}

		dag := internal.NewDAG[bool]()
		for _, t := range pod.Spec.Tasks {
			dag.AddNode(t.Name, true)
			for _, dependency := range t.Dependencies {
				dag.AddEdge(dependency, t.Name)
			}
		}
		visited := dag.Subgraph(args)

		taskByName := pod.Spec.Tasks.Map()
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
							log.Printf("file changed, re-running %s\n", t.task.Name)
							taskNames <- t.task.Name
						}
					}
				}
			}()
		}

		semaphores := util.NewSemaphores(pod.Spec.Semaphores)

		wg := sync.WaitGroup{}

		for {
			select {
			case <-ctx.Done():

				wg.Wait()

				// if any task failed, we will return an error code
				for _, node := range subgraph.Nodes {
					if node.status == "failed" {
						return errors.New("one or more tasks failed")
					}
				}

				return nil
			case taskName := <-taskNames:
				// we will only execute this task, if its parents are "succeeded" or "skipped" or ("running" and the task is a service)
				blocked := false
				for _, parentName := range subgraph.Parents[taskName] {
					parent := subgraph.Nodes[parentName]
					if parent.blocked() {
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

				subgraph.Nodes[taskName].status = "waiting"

				go func(t types.Task) {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					defer wg.Done()

					out := funcWriter(func(bytes []byte) (int, error) {
						prefix := fmt.Sprintf("%s[%s] (%s) ", internal.Color(t.Name, t.IsService()), t.Name, subgraph.Nodes[t.Name].status)

						// split on newlines
						lines := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
						for _, line := range lines {
							log.Println(prefix + line)
						}

						return len(bytes), nil
					})

					log := log.New(out, "", 0)

					queueChildren := func() {
						for _, child := range subgraph.Children[t.Name] {
							taskNames <- child
						}
					}

					// if the task can be skipped, lets exit early
					if t.Skip() {
						node.status = "skipped"
						log.Println("skipping")
						queueChildren()
						return
					}

					// if the task needs a mutex, lets wait for it
					if t.Mutex != "" {
						mu := util.GetMutex(t.Mutex)
						log.Println("waiting for mutex")
						mu.Lock()
						log.Println("acquired mutex")
						defer mu.Unlock()
					}

					// if the task needs a semaphore, lets wait for it
					if t.Semaphore != "" {
						sema := semaphores.Get(t.Semaphore)
						log.Println("waiting for semaphore")
						if err := sema.Acquire(ctx, 1); err != nil {
							log.Println("failed to acquire semaphore")
							return
						}
						defer sema.Release(1)
					}

					p := proc.New(t, log, pod.Spec)

					if probe := t.GetLivenessProbe(); probe != nil {
						liveFunc := func(live bool, err error) {
							if !live {
								node.status = "failed"
								cancel()
							}
						}
						go probeLoop(ctx, *probe, liveFunc)
					}
					if probe := t.GetReadinessProbe(); probe != nil {
						readyFunc := func(ready bool, err error) {
							if ready {
								node.status = "running"
								queueChildren()
							} else {
								node.status = "failed"
								cancel()
							}
						}
						go probeLoop(ctx, *probe, readyFunc)
					}

					if t.IsService() {
						node.status = "starting"
					} else {
						node.status = "running"
					}

					restart := func() {
						select {
						case <-ctx.Done():
						case <-time.After(3 * time.Second):
							log.Println("restarting")
							taskNames <- t.Name
						}
					}

					if t.Log != "" {
						out, err := os.Create(t.Log)
						if err != nil {
							log.Printf("failed to create log file: %v", err)
							return
						}
						defer out.Close()
					}

					err = p.Run(ctx, out, out)
					if err != nil {
						node.status = "failed"
						log.Println(err)
						if t.GetRestartPolicy() != "Never" {
							restart()
						}
						return
					}

					node.status = "succeeded"
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
