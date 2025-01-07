package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kitproj/kit/internal/proc"
	"github.com/kitproj/kit/internal/types"
	"github.com/kitproj/kit/internal/util"
	"k8s.io/utils/strings/slices"
)

var poisonPill = struct{}{}

func RunSubgraph(
	ctx context.Context,
	cancel context.CancelFunc,
	logger *log.Logger,
	wf *types.Workflow,
	taskNames []string,
	tasksToSkip []string,
) error {

	// check that the task names are valid
	for _, name := range taskNames {
		if _, ok := wf.Tasks[name]; !ok {
			return fmt.Errorf("task %q not found in workflow", name)
		}
	}

	// check skipped tasks are valid
	for _, name := range tasksToSkip {
		if _, ok := wf.Tasks[name]; !ok {
			return fmt.Errorf("skipped task %q not found in workflow", name)
		}
	}

	dag := NewDAG[bool]()
	for name, t := range wf.Tasks {
		dag.AddNode(name, true)
		for _, dependency := range t.Dependencies {
			dag.AddEdge(dependency, name)
		}
	}
	visited := dag.Subgraph(taskNames)

	taskByName := wf.Tasks
	subgraph := NewDAG[*TaskNode]()
	for name := range visited {
		task := taskByName[name]
		subgraph.AddNode(name, &TaskNode{name: name, task: task, phase: "pending", cancel: func() {}})
		for _, parent := range dag.Parents[name] {
			subgraph.AddEdge(parent, name)
		}
	}

	events := make(chan any, len(subgraph.Nodes)*2)

	// schedule the tasks in the subgraph that are ready to run , this is done by sending the task name to the events channel of any task that does not have any parents
	for taskName := range subgraph.Nodes {
		if len(subgraph.Parents[taskName]) == 0 {
			events <- taskName
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
						logger.Printf("file changed, re-running %s\n", node.name)
						events <- node.name
						node.cancel()
					}
				}
			}
		}()
	}

	semaphores := util.NewSemaphores(wf.Semaphores)

	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():

			logger.Println("waiting for all tasks to finish")

			wg.Wait()

			// if any task failed, we will return an error
			var failures []string
			for _, node := range subgraph.Nodes {
				logger.Printf("%s %s\n", node.name, node.phase)
				if node.phase == "failed" {
					failures = append(failures, fmt.Sprintf("%s: %s", node.name, node.message))
				}
			}

			if len(failures) > 0 {
				return fmt.Errorf("failed tasks: %v", failures)
			}

			return nil
		case event := <-events:
			switch x := event.(type) {
			// if we get the poison pill, we should see if any job tasks are failed, if so we must exist
			// if all jobs are either succeeded or skipped, we can exit
			case struct{}:
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
			// if the event is a string, it is the name of the task to run
			case string:
				taskName := x

				// we will only execute this task, if its parents are "succeeded" or "skipped" or ("running" and the task is a service)
				blocked := false
				for _, parentName := range subgraph.Parents[taskName] {
					parent := subgraph.Nodes[parentName]
					if parent.blocked() {
						logger.Printf("task %q is blocked by %q (%s): %s\n", taskName, parentName, parent.phase, parent.message)
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

				go func(node *TaskNode) {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					node.cancel = cancel

					// send a poison pill to indicate that we've finish and the main loop must check to see if we need to exit
					defer func() { events <- poisonPill }()
					defer wg.Done()

					t := node.task

					var out io.Writer = funcWriter(func(bytes []byte) (int, error) {
						prefix := fmt.Sprintf("%s[%s] (%s) ", Color(node.name, t.IsService()), node.name, node.phase)
						// reset color and bold
						suffix := "\033[0m"

						// split on newlines
						lines := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
						for _, line := range lines {
							logger.Println(prefix + line + suffix)
						}

						return len(bytes), nil
					})

					logger := log.New(out, "", 0)

					setNodeStatus := func(node *TaskNode, phase string, message string) {
						node.phase = phase
						node.message = message
						logger.Println(message)
					}

					setNodeStatus(node, "waiting", "")

					queueChildren := func() {
						for _, child := range subgraph.Children[node.name] {
							// only queue tasks in the subgraph
							if _, ok := subgraph.Nodes[child]; ok {
								logger.Printf("queuing %q\n", child)
								events <- child
							}
						}
					}

					// if the task can be skipped, lets exit early
					if t.Skip() || slices.Contains(tasksToSkip, node.name) {
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

					p := proc.New(t, logger, types.Spec(*wf))

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
							logger.Println("restarting")
							events <- node.name
						}
					}

					if t.Log != "" {
						file, err := os.Create(t.Log)
						if err != nil {
							setNodeStatus(node, "failed", fmt.Sprintf("failed to create log file: %v", err))
							return
						}
						out = file
						defer file.Close()
					}

					err := p.Run(ctx, out, out)
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
			default:
				panic(fmt.Sprintf("unexpected event: %v", event))
			}
		}
	}
}
