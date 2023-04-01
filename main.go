package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/kitproj/kit/internal/proc"
	"github.com/kitproj/kit/internal/types"
	"github.com/kitproj/kit/internal/util"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
	k8sstrings "k8s.io/utils/strings"
	"sigs.k8s.io/yaml"
)

//go:generate sh -c "git describe --tags > tag"
//go:embed tag
var tag string

// GitHub Actions
var isCI = os.Getenv("CI") != "" || // Travis CI, CircleCI, GitLab CI, AppVeyor, CodeShip, dsari
	os.Getenv("BUILD_ID") != "" || // Jenkins, TeamCity
	os.Getenv("RUN_ID") != "" || // TaskCluster, Codefresh
	os.Getenv("GITHUB_ACTIONS") == "true"
var isTerminal = isatty.IsTerminal(os.Stdin.Fd())
var muxOutput = isTerminal && !isCI
var faint = color.New(color.Faint).Sprint

func init() {
	if muxOutput {
		_ = os.Mkdir("logs", 0o777)
		f, err := os.Create("logs/kit.log")
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.Printf("tag=%v\n", tag)
	log.Printf("isTerminal=%v, isCI=%v\n", isTerminal, isCI)
}

const escape = "\x1b"

const defaultConfigFile = "tasks.yaml"

type message struct{ text, level string }
type taskStatus struct {
	reason  string
	message message
	stdout  io.Writer
	stderr  io.Writer
	recent  *util.LastNLinesWriter
	backoff backoff
}

func last(p string) string {
	parts := strings.Split(strings.TrimSpace(p), "\n")
	return parts[len(parts)-1]
}

func main() {
	help := false
	configFile := ""
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&configFile, "f", defaultConfigFile, "config file")
	flag.Parse()
	args := flag.Args()

	if help {
		fmt.Println("Usage: kit [-f tasks.yaml] [task1 task2 ...]")
		os.Exit(0)
	}

	started := time.Now()

	err := func() error {

		ctx, stopEverything := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		defer stopEverything()

		_ = os.Mkdir("logs", 0777)

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

		tasks := pod.Spec.Tasks.NeededFor(args)

		log.Printf("tasks: %v\n", tasks)

		statuses := sync.Map{}

		for _, task := range tasks {
			nLinesWriter := util.NewLastNLinesWriter(48)
			x := &taskStatus{
				reason:  "waiting",
				stdout:  io.MultiWriter(nLinesWriter, &prefixWriter{task.Name + ": ", os.Stdout}),
				stderr:  io.MultiWriter(nLinesWriter, &prefixWriter{task.Name + ": ", os.Stdout}),
				recent:  nLinesWriter,
				backoff: defaultBackoff,
			}
			if muxOutput {
				logFile, err := os.Create(filepath.Join("logs", task.Name+".log"))
				if err != nil {
					panic(err)
				}
				x.stdout = io.MultiWriter(nLinesWriter, logFile, funcWriter(func(p []byte) (n int, err error) {
					x.message = message{last(string(p)), "info"}
					return len(p), nil
				}))
				x.stderr = io.MultiWriter(nLinesWriter, logFile, io.MultiWriter(logFile, funcWriter(func(p []byte) (n int, err error) {
					x.message = message{last(string(p)), "error"}
					return len(p), nil
				})))
			}
			statuses.Store(task.Name, x)
		}
		terminating := false
		printTasks := func() {
			width, height, _ := terminal.GetSize(0)
			if width == 0 {
				width = 80
			}
			if height == 0 {
				height = 24
			}
			fmt.Printf("%s[2J", escape)   // clear screen
			fmt.Printf("%s[0;0H", escape) // move to 0,0
			// how many spaces left to print logs
			space := height - 2
			for _, t := range pod.Spec.Tasks {
				v, ok := statuses.Load(t.Name)
				if !ok {
					continue
				}
				status := v.(*taskStatus)
				if status.reason == "success" {
					continue
				}
				reason := status.reason
				const blackSquare = "■"
				icon := blackSquare
				switch reason {
				case "running":
					icon = color.BlueString(blackSquare)
				case "ready":
					icon = color.GreenString(blackSquare)
				case "error":
					icon = color.RedString(blackSquare)
				}
				prefix := fmt.Sprintf("%s %-10s %-8s", icon, k8sstrings.ShortenString(t.Name, 10), reason)
				if ports := t.GetHostPorts(); len(ports) > 0 {
					prefix = prefix + " " + faint(ports)
				}
				entry := status.message
				msgWidth := width - len(prefix) - 1
				msg := ""
				if msgWidth > 0 {
					msg = k8sstrings.ShortenString(entry.text, msgWidth)
					if entry.level == "error" {
						msg = color.RedString(msg)
					}
				}
				fmt.Println(prefix + " " + msg)
				space--
			}
			items := []string{
				strings.TrimSpace(tag),
				fmt.Sprint(time.Since(started).Truncate(time.Second)),
				"logs in ./logs",
			}
			if terminating {
				items = append(items, "terminating...")
			}
			fmt.Println(faint(strings.Join(items, "   ")))

			if terminating {
				return
			}

			var printLogs = func(test func(task types.Task, status *taskStatus) bool) {
				for _, t := range pod.Spec.Tasks {
					v, ok := statuses.Load(t.Name)
					if !ok {
						continue
					}
					status := v.(*taskStatus)
					if test(t, status) {
						for _, msg := range status.recent.Lines() {
							if space <= 0 {
								break
							}
							fmt.Println(k8sstrings.ShortenString(fmt.Sprintf("%-10s: %s", k8sstrings.ShortenString(t.Name, 10), msg), width))
							space--
						}
					}
				}
			}
			printLogs(func(task types.Task, status *taskStatus) bool { return status.reason == "error" })
			printLogs(func(task types.Task, status *taskStatus) bool {
				return status.reason == "running" && task.ReadinessProbe != nil
			})
		}

		if muxOutput {
			// every 1/2 second print the current status to the terminal
			go func() {
				defer handleCrash(stopEverything)
				for {
					printTasks()
					time.Sleep(time.Second / 10)
				}
			}()
		}

		work := make(chan types.Task)
		semaphores := util.NewSemaphores(pod.Spec.Semaphores)

		log.Printf("semaphores=%v\n", semaphores)

		go func() {
			defer handleCrash(stopEverything)
			for _, t := range tasks.GetLeaves() {
				work <- t
			}
		}()

		go func() {
			defer handleCrash(stopEverything)
			<-ctx.Done()
			terminating = true
			close(work)
		}()

		wg := sync.WaitGroup{}

		stop := &sync.Map{}

		maybeStartDownstream := func(name string) {
			select {
			case <-ctx.Done():
			default:
				log.Printf("%s: starting downstream tasks\n", name)
				for _, downstream := range tasks.GetDownstream(name) {
					fulfilled := true
					for _, upstream := range downstream.Dependencies {
						v, ok := statuses.Load(upstream)
						if ok {
							status := v.(*taskStatus)
							fulfilled = fulfilled && (status.reason == "success" || status.reason == "ready")
						} else {
							fulfilled = false
						}
					}
					if fulfilled {
						log.Printf("%s: starting: %s\n", name, downstream.Name)
						work <- downstream
					}
				}
			}
		}

		// stop everything if all tasks are complete/in error
		go func() {
			defer handleCrash(stopEverything)
			for {
				complete := true
				for _, task := range tasks {
					if v, ok := statuses.Load(task.Name); ok {
						status := v.(*taskStatus)
						complete = complete && !task.IsBackground() && (status.reason == "success" || status.reason == "error")
					} else {
						complete = false
					}
				}
				if complete {
					stopEverything()
				}
				time.Sleep(time.Second)
			}
		}()

		for t := range work {
			name := t.Name

			prc := proc.New(t, pod.Spec)
			v, _ := statuses.Load(t.Name)
			status := v.(*taskStatus)

			processCtx, stopProcess := context.WithCancel(ctx)
			defer stopProcess()

			// watch files for changes
			go func(t types.Task, stopProcess func()) {
				defer handleCrash(stopEverything)
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					panic(err)
				}
				defer watcher.Close()
				for _, w := range t.Watch {
					stat, err := os.Stat(w)
					if err != nil {
						panic(err)
					}
					if err := watcher.Add(w); err != nil {
						panic(err)
					}
					if stat.IsDir() {
						if err := filepath.WalkDir(w, func(path string, d fs.DirEntry, err error) error {
							if err != nil {
								return err
							}
							if d.IsDir() {
								log.Printf("%s: watching %q\n", t.Name, path)
								return watcher.Add(path)
							}
							return nil
						}); err != nil {
							panic(err)
						}
					}
				}

				timer := time.AfterFunc(100*365*24*time.Hour, func() {
					work <- t
				})
				defer timer.Stop()

				for {
					select {
					case <-processCtx.Done():
						return
					case e := <-watcher.Events:
						// ignore chmod events, these can be triggered by the editor, but we don't care
						if e.Op != fsnotify.Chmod {
							log.Printf("%s: %v changed\n", t.Name, e)
							timer.Reset(time.Second)
						}
					case err := <-watcher.Errors:
						panic(err)
					}
				}
			}(t, stopProcess)

			// run the process
			wg.Add(1)
			go func(t types.Task, status *taskStatus, stopProcess func()) {
				defer handleCrash(stopEverything)
				defer wg.Done()

				if f, ok := stop.Load(name); ok {
					log.Printf("%s: stopping process\n", t.Name)
					f.(func())()
				}

				stop.Store(name, stopProcess)

				if m := t.Mutex; m != "" {
					mutex := util.GetMutex(m)
					_, _ = fmt.Fprintf(status.stdout, "waiting for mutex %q\n", m)
					mutex.Lock()
					_, _ = fmt.Fprintf(status.stdout, "locked mutex %q\n", m)
					defer mutex.Unlock()
				}

				if s := t.Semaphore; s != "" {
					_, _ = fmt.Fprintf(status.stdout, "waiting for semaphore %q\n", s)
					semaphore := semaphores.Get(s)
					if err := semaphore.Acquire(ctx, 1); err != nil {
						return
					}
					_, _ = fmt.Fprintf(status.stdout, "acquired semaphore %q\n", s)
					defer semaphore.Release(1)
				}

				go func() {
					defer handleCrash(stopEverything)
					<-ctx.Done()
					stopProcess()
				}()

				var stdout = status.stdout
				var stderr = status.stderr
				for {
					select {
					case <-processCtx.Done():
						return
					default:
						log.Printf("%s: starting process\n", t.Name)
						err := func() error {
							runCtx, stopRun := context.WithCancel(processCtx)
							defer stopRun()
							log.Printf("%s: resetting process\n", t.Name)
							if err := prc.Reset(runCtx); err != nil {
								return err
							}
							for _, port := range t.GetHostPorts() {
								log.Printf("%s: waiting for port %d to be free\n", t.Name, port)
								if err := isPortFree(port); err != nil {
									return err
								}
							}
							status.reason = "running"
							if probe := t.GetLivenessProbe(); probe != nil {
								log.Printf("%s: liveness probe=%v\n", t.Name, probe)
								liveFunc := func(live bool, err error) {
									if !live {
										log.Printf("%s: is dead, stopping\n", t.Name)
										stopRun()
									}
								}
								go probeLoop(runCtx, stopEverything, *probe, liveFunc)
							}
							if probe := t.GetReadinessProbe(); probe != nil {
								log.Printf("%s: readiness probe=%v\n", t.Name, probe)
								readyFunc := func(ready bool, err error) {
									if ready {
										log.Printf("%s: is ready, starting downstream\n", t.Name)
										status.reason = "ready"
										maybeStartDownstream(name)
									} else {
										log.Printf("%s: is not ready\n", t.Name)
										status.reason = "running"
									}
								}
								go probeLoop(runCtx, stopEverything, *probe, readyFunc)
							}
							log.Printf("%s: running process\n", t.Name)
							return prc.Run(runCtx, stdout, stderr)
						}()
						if err != nil {
							if errors.Is(err, context.Canceled) {
								return
							}
							status.reason = "error"
							_, _ = fmt.Fprintln(stderr, err.Error())
							status.backoff = status.backoff.next()
						} else {
							status.reason = "success"
							maybeStartDownstream(name)
							if !t.IsBackground() && t.GetRestartPolicy() != "Always" {
								return
							}
						}
						if t.GetRestartPolicy() == "Never" {
							return
						}
					}
					_, _ = fmt.Fprintf(stdout, "backing off %s\n", status.backoff.Duration)
					time.Sleep(status.backoff.Duration)
				}
			}(t, status, stopProcess)
		}

		wg.Wait()

		for _, task := range tasks {
			if v, ok := statuses.Load(task.Name); ok && v.(*taskStatus).reason == "error" {
				return fmt.Errorf("%s failed", task.Name)
			}
		}
		return nil
	}()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func handleCrash(stop func()) {
	if r := recover(); r != nil {
		fmt.Println(r)
		debug.PrintStack()
		stop()
	}
}
