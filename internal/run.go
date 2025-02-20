package internal

import (
	"bytes"
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
	"github.com/pkg/browser"
	"k8s.io/utils/strings/slices"
)

var poisonPill = struct{}{}

func RunSubgraph(ctx context.Context, cancel context.CancelFunc, port int, openBrowser bool, logger *log.Logger, wf *types.Workflow, taskNames []string, tasksToSkip []string) error {

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

	// name is last part of pwd
	pwd := os.Getenv("PWD")
	name := filepath.Base(pwd)

	dag := NewDAG[bool](name)
	for name, t := range wf.Tasks {
		dag.AddNode(name, true)
		for _, dependency := range t.Dependencies {
			dag.AddEdge(dependency, name)
		}
	}
	visited := dag.Subgraph(taskNames)

	taskByName := wf.Tasks
	subgraph := NewDAG[*TaskNode](name)
	for name := range visited {
		task := taskByName[name]

		logFile := filepath.Join("logs", fmt.Sprintf("%s.log", name))
		if task.Log != "" {
			logFile = task.Log
		}

		subgraph.AddNode(name, &TaskNode{
			Name:    name,
			logFile: logFile,
			task:    task,
			Phase:   "pending",
			cancel:  func() {},
			mu:      &sync.Mutex{}})
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

	if len(subgraph.Nodes) == 0 {
		logger.Println("no tasks to run")
		return nil
	}

	// create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// start a file watcher for each task
	for _, node := range subgraph.Nodes {

		// start watching files for changes
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("failed to create watcher: %w", err)
		}
		for _, source := range node.task.Watch {
			if err := watcher.Add(filepath.Join(node.task.WorkingDir, source)); err != nil {
				return fmt.Errorf("failed to watch %q: %w", source, err)
			}
		}
		defer watcher.Close()

		go func() {
			debounceTimer := time.AfterFunc(0, func() {})
			defer debounceTimer.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						debounceTimer.Stop()
						debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
							logger.Printf("[%s] %s changed, re-running\n", node.Name, event.Name)
							events <- node.Name
						})
					}
				}
			}
		}()
	}

	semaphores := util.NewSemaphores(wf.Semaphores)

	wg := &sync.WaitGroup{}

	statusEvents := make(chan *TaskNode, 100)

	if port > 0 {
		go StartServer(ctx, port, wg, subgraph, statusEvents)
		if openBrowser {
			if err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", port)); err != nil {
				return fmt.Errorf("failed to open browser: %v", err)
			}
		}
	}

	stallTimers := map[string]*time.Timer{}
	for name, taskNode := range subgraph.Nodes {
		stalledTime := taskNode.task.GetStalledTimeout()
		stallTimers[name] = time.AfterFunc(stalledTime, func() {
			if taskNode.Phase == "starting" || taskNode.Phase == "running" {
				// we suffix the message with "starting" so we can differentiate between a task that is starting and one that is running, later on we can change the message to "output received"
				// and restore the phase to "running" or "starting"
				taskNode.Message = fmt.Sprintf("no output for %s or more while %s", stalledTime, taskNode.Phase)
				taskNode.Phase = "stalled"
				logger.Printf("[%s] %s\n", taskNode.Name, taskNode.Message)
				statusEvents <- taskNode
			}
		})
	}

	for {
		select {
		case <-ctx.Done():

			logger.Println("waiting for all tasks to complete")

			wg.Wait()

			// if any task failed, we will return an error
			var failures []string
			for _, node := range subgraph.Nodes {

				color := 30
				faint := 0
				switch node.Phase {
				case "failed":
					// red
					color = 31
					faint = 1
					failures = append(failures, node.Name)
				case "pending", "waiting":
					faint = 2
				}

				logger.Printf("\033[%d;%dm[%s] (%s) %s\033[0m\n", faint, color, node.Name, node.Phase, node.Message)
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
				// if all requests tasks are succeeded, we can exit
				{
					pendingTasks := map[string]bool{}
					for _, x := range taskNames {
						pendingTasks[x] = true
					}

					for _, node := range subgraph.Nodes {
						if (node.Phase == "succeeded" || node.Phase == "skipped") && node.task.GetRestartPolicy() != "Always" {
							delete(pendingTasks, node.Name)
						}
					}

					if len(pendingTasks) == 0 {
						logger.Println("exiting because all requested tasks completed and none should be restarted")
						cancel()
					}
				}

				// if a task that should not be restarted failed, we must exit
				for _, node := range subgraph.Nodes {
					if node.Phase == "failed" && node.task.GetRestartPolicy() == "Never" {
						logger.Printf("exiting because task  %q should not be restarted, and it failed", node.Name)
						cancel()
					}
				}

			// if the event is a string, it is the name of the task to run
			case string:
				taskName := x

				// we will only execute this task, if its parents are "succeeded" or "skipped" or ("running" and the task is a service)
				blocked := false
				for _, parentName := range subgraph.Parents[taskName] {
					parent := subgraph.Nodes[parentName]
					if parent.blocked() {
						logger.Printf("task %q is blocked by %q (%s): %s\n", taskName, parentName, parent.Phase, parent.Message)
						blocked = true
					}
				}

				if blocked {
					continue
				}

				// we might already be pending, waiting, starting or running this task, so we don't want to start it again
				node := subgraph.Nodes[taskName]

				node.cancel()

				// each task is executed in a separate goroutine
				wg.Add(1)

				go func(node *TaskNode) {

					// lock the task, so we do not run two instances of it at the same time
					node.mu.Lock()

					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					node.cancel = cancel

					// send a poison pill to indicate that we've finish and the main loop must check to see if we need to exit
					defer func() { events <- poisonPill }()
					defer wg.Done()
					defer node.mu.Unlock()

					t := node.task

					var out io.Writer = funcWriter(func(p []byte) (int, error) {
						prefix := fmt.Sprintf("%s[%s] (%s)  ", color(node.Name), node.Name, node.Phase)
						// reset color and bold
						suffix := "\033[0m"

						lines := bytes.Split(p, []byte("\n"))
						for i, line := range lines {
							if len(line) == 0 {
								continue
							}
							if i == len(lines)-1 {
								logger.Printf("%s%s%s", prefix, string(line), suffix)
							} else {
								logger.Printf("%s%s%s\n", prefix, string(line), suffix)
							}
						}

						return len(p), nil
					})

					logger := log.New(out, "", 0)

					setNodeStatus := func(node *TaskNode, phase string, message string) {
						node.Phase = phase
						node.Message = message
						stallTimers[node.Name].Reset(node.task.GetStalledTimeout())
						logger.Println(node.Message)
						statusEvents <- node
					}

					setNodeStatus(node, "waiting", "")

					queueChildren := func() {
						for _, child := range subgraph.Children[node.Name] {
							// only queue tasks in the subgraph
							if _, ok := subgraph.Nodes[child]; ok {
								logger.Printf("queuing %q\n", child)
								events <- child
							}
						}
					}

					// if the task can be skipped, lets exit early
					if t.Skip() || slices.Contains(tasksToSkip, node.Name) {
						setNodeStatus(node, "skipped", "")
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

					p := proc.New(taskName, t, logger, types.Spec(*wf))

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

					if t.GetType() == types.TaskTypeService {
						if t.Ports != nil {
							setNodeStatus(node, "starting", "service starting")
						} else {
							setNodeStatus(node, "running", "no ports to expose")
							queueChildren()
						}
					} else {
						// non a service, must be a job
						setNodeStatus(node, "running", "job running")
					}

					restart := func() {
						select {
						case <-ctx.Done():
						case <-time.After(3 * time.Second):
							logger.Println("restarting")
							cancel()
							events <- node.Name
						}
					}

					file, err := os.Create(node.logFile)
					if err != nil {
						setNodeStatus(node, "failed", fmt.Sprintf("failed to create log file: %v", err))
						return
					}
					defer file.Close()

					// if the task has a log file, we will write to that file, we sync after each write
					// so when we tail the log file, we see the output immediately
					buf := funcWriter(func(p []byte) (int, error) {
						stallTimers[node.Name].Reset(node.task.GetStalledTimeout())
						if node.Phase == "stalled" {
							if strings.HasSuffix(node.Message, "starting") {
								setNodeStatus(node, "starting", "output received")
							} else {
								setNodeStatus(node, "running", "output received")
							}
						}
						n, err := file.Write(p)
						if err != nil {
							return n, err
						}
						if err := file.Sync(); err != nil {
							return n, err
						}
						return n, nil
					})

					if t.Log != "" {
						out = buf
					} else {
						out = io.MultiWriter(out, buf)
					}

					err = p.Run(ctx, out, out)
					// if the task was cancelled, we don't want to restart it, this is normal exit
					if errors.Is(ctx.Err(), context.Canceled) {
						setNodeStatus(node, "cancelled", "")
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
