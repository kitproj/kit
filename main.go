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
			subgraph.AddNode(name, &taskNode{name: name, task: task, phase: "pending", cancel: func() {}})
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
		for _, node := range subgraph.Nodes {

			// start watching files for changes
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}
			for _, source := range node.task.Watch {
				if err := watcher.Add(filepath.Join(node.task.WorkingDir, source)); err != nil {
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
							log.Printf("file changed, re-running %s\n", node.name)
							taskNames <- node.name
							node.cancel()
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

				// if we get the poison pill, we should see if any job tasks are failed, if so we must exist
				// if all jobs are either succeeded or skipped, we can exit
				const PoisonPill = ""
				if taskName == PoisonPill {
					anyJobFailed := false
					numJobsSucceeded := 0
					for _, node := range subgraph.Nodes {
						if node.task.IsService() {
							continue
						}
						if node.phase == "failed" {
							anyJobFailed = true
						}
						if node.phase == "succeeded" {
							numJobsSucceeded++
						}
					}
					everyJobSucceeded := numJobsSucceeded == len(subgraph.Nodes)
					if everyJobSucceeded || anyJobFailed {
						cancel()
					}
					continue
				}

				// we will only execute this task, if its parents are "succeeded" or "skipped" or ("running" and the task is a service)
				blocked := false
				for _, parentName := range subgraph.Parents[taskName] {
					parent := subgraph.Nodes[parentName]
					if parent.blocked() {
						log.Printf("task %q is blocked by %q (%s): %s\n", taskName, parentName, parent.phase, parent.message)
						blocked = true
					}
				}

				if blocked {
					continue
				}

				// we might already be waiting, starting or running this task, so we don't want to start it again
				node := subgraph.Nodes[taskName]
				node.cancel()

				// each task is executed in a separate goroutine
				wg.Add(1)

				go func(node *taskNode) {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					node.cancel = cancel

					// send a poison pill to indicate that we've finish and the main loop must check to see if we need to exit
					defer func() { taskNames <- PoisonPill }()
					defer wg.Done()

					t := node.task

					out := funcWriter(func(bytes []byte) (int, error) {
						prefix := fmt.Sprintf("%s[%s] (%s) ", internal.Color(node.name, t.IsService()), node.name, node.phase)
						// reset color and bold
						suffix := "\033[0m"

						// split on newlines
						lines := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
						for _, line := range lines {
							fmt.Println(prefix + line + suffix)
						}

						return len(bytes), nil
					})

					log := log.New(out, "", 0)

					setNodeStatus := func(node *taskNode, phase string, message string) {
						node.phase = phase
						node.message = message
						log.Println(message)
					}

					setNodeStatus(node, "waiting", "")

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
						setNodeStatus(node, "succeeded", "skipped")
						queueChildren()
						return
					}

					// if the task needs a mutex, lets wait for it
					if t.Mutex != "" {
						mu := util.GetMutex(t.Mutex)
						setNodeStatus(node, "waiting", "waiting for mutex")
						mu.Lock()
						setNodeStatus(node, "waiting", "acquired mutex")
						defer mu.Unlock()
					}

					// if the task needs a semaphore, lets wait for it
					if t.Semaphore != "" {
						sema := semaphores.Get(t.Semaphore)
						setNodeStatus(node, "waiting", "waiting for semaphore")
						if err := sema.Acquire(ctx, 1); err != nil {
							setNodeStatus(node, "failed", fmt.Sprintf("failed to acquire semaphore: %v", err))
							return
						}
						setNodeStatus(node, "waiting", "acquired semaphore")
						defer sema.Release(1)
					}

					p := proc.New(t, log, types.PodSpec(*pod))

					if probe := t.GetLivenessProbe(); probe != nil {
						liveFunc := func(live bool, err error) {
							if !live {
								setNodeStatus(node, "failed", fmt.Sprintf("liveness probe failed: %v", err))
								cancel()
							}
						}
						go probeLoop(ctx, *probe, liveFunc)
					}
					if probe := t.GetReadinessProbe(); probe != nil {
						readyFunc := func(ready bool, err error) {
							if ready {
								setNodeStatus(node, "running", "readiness probe succeeded")
								queueChildren()
							} else {
								setNodeStatus(node, "failed", fmt.Sprintf("readiness probe failed: %v", err))
								cancel()
							}
						}
						go probeLoop(ctx, *probe, readyFunc)
					}

					if t.IsService() {
						setNodeStatus(node, "starting", "")
					} else {
						// non a service, must be a job
						setNodeStatus(node, "running", "")
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
							setNodeStatus(node, "failed", fmt.Sprintf("failed to create log file: %v", err))
							return
						}
						defer out.Close()
					}

					err = p.Run(ctx, out, out)
					// if the task was cancelled, we don't want to restart it, this is normal exit
					if errors.Is(ctx.Err(), context.Canceled) {
						return
					}

					if err != nil {
						setNodeStatus(node, "failed", fmt.Sprint(err))
						if t.GetRestartPolicy() != "Never" {
							restart()
						}
						return
					}

					setNodeStatus(node, "succeeded", "")
					if t.GetRestartPolicy() == "Always" {
						restart()
					}
					queueChildren()

				}(node)
			}
		}
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
