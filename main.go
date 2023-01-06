package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/alexec/kit/internal/proc"
	"github.com/alexec/kit/internal/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	k8sstrings "k8s.io/utils/strings"
	"sigs.k8s.io/yaml"
)

func init() {
	_ = os.Mkdir("logs", 0o777)
	f, err := os.OpenFile("logs/kit.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.Println(tag)
}

const escape = "\x1b"

const defaultConfigFile = "kit.yaml"

func main() {
	cmd := &cobra.Command{
		Use:   "kit [config_file]",
		Short: "Start-up processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := defaultConfigFile

			ctx, stopEverything := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
			defer stopEverything()

			_ = os.Mkdir("logs", 0777)

			pod := &types.Pod{}
			statuses := &types.Status{}
			logEntries := make(map[string]*types.LogEntry)

			go func() {
				defer handleCrash(stopEverything)
				for {
					width, _, _ := terminal.GetSize(0)
					if width == 0 {
						width = 80
					}
					fmt.Printf("%s[2J", escape)   // clear screen
					fmt.Printf("%s[0;0H", escape) // move to 0,0
					for _, t := range pod.Spec.Tasks {
						state := statuses.GetStatus(t.Name)
						if state == nil {
							continue
						}
						if state.IsSuccess() {
							continue
						}
						reason := state.GetReason()
						icon := "▓"
						switch reason {
						case "running":
							icon = color.BlueString("▓")
						case "ready":
							icon = color.GreenString("▓")
						case "error":
							icon = color.RedString("▓")
						}
						prefix := fmt.Sprintf("%s %-10s %-8s %v", icon, k8sstrings.ShortenString(state.Name, 10), reason, color.HiBlackString(fmt.Sprint(t.GetHostPorts())))
						entry := logEntries[t.Name]
						msg := k8sstrings.ShortenString(entry.Msg, width-len(prefix)-1)
						if entry.IsError() {
							msg = color.YellowString(msg)
						}
						fmt.Println(prefix + " " + msg)
					}
					time.Sleep(time.Second)
				}
			}()

			in, err := os.ReadFile(configFile)
			if err != nil {
				return err
			}
			if err = yaml.UnmarshalStrict(in, pod); err != nil {
				return err
			}
			data, err := yaml.Marshal(pod)
			if err != nil {
				return err
			}
			if err = os.WriteFile(configFile, data, 0o0644); err != nil {
				return err
			}

			tasks := make(chan types.Task)

			go func() {
				defer handleCrash(stopEverything)
				for _, t := range pod.Spec.Tasks.GetLeaves() {
					tasks <- t
				}
			}()

			go func() {
				defer handleCrash(stopEverything)
				<-ctx.Done()
				close(tasks)
			}()

			wg := sync.WaitGroup{}

			stopAndWait := make(map[string]func())

			maybeStartDownstream := func(name string) {
				for _, downstream := range pod.Spec.Tasks.GetDownstream(name) {
					ok := true
					for _, upstream := range downstream.Dependencies {
						x := statuses.GetStatus(upstream)
						ok = ok && (x.IsSuccess() || x.IsReady())
					}
					if ok {
						tasks <- downstream
					}
				}
			}

			for t := range tasks {
				name := t.Name

				if f, ok := stopAndWait[name]; ok {
					f()
				}

				logEntries[name] = &types.LogEntry{}
				if statuses.GetStatus(name) == nil {
					statuses.TaskStatuses = append(statuses.TaskStatuses, &types.TaskStatus{
						Name: name,
					})
				}

				logEntry := logEntries[name]

				prc := proc.New(t, pod.Spec)

				processCtx, stopProcess := context.WithCancel(ctx)
				defer stopProcess()

				pwg := &sync.WaitGroup{}

				stopAndWait[name] = func() {
					stopProcess()
					pwg.Wait()
				}

				go func(t types.Task, stopProcess func()) {
					defer handleCrash(stopEverything)
					last := time.Now()
					for {
						select {
						case <-processCtx.Done():
							return
						default:
							for _, file := range t.Watch {
								stat, err := os.Stat(file)
								if err != nil {
									return
								}
								if stat.ModTime().After(last) {
									tasks <- t
									break
								}
							}
							last = time.Now()
							time.Sleep(time.Second)
						}
					}
				}(t, stopProcess)
				wg.Add(1)
				pwg.Add(1)
				go func(t types.Task, status *types.TaskStatus) {
					defer handleCrash(stopEverything)
					defer wg.Done()
					defer pwg.Done()
					logFile, err := os.Create(filepath.Join("logs", name+".log"))
					if err != nil {
						panic(err)
					}
					defer logFile.Close()
					stdout := io.MultiWriter(logFile, logEntry.Stdout())
					stderr := io.MultiWriter(logFile, logEntry.Stderr())
					for {
						select {
						case <-processCtx.Done():
							return
						default:
							err := func() error {
								runCtx, stopRun := context.WithCancel(processCtx)
								defer stopRun()
								go func() {
									defer handleCrash(stopEverything)
									<-ctx.Done()
									stopProcess()
								}()
								status.State = types.TaskState{
									Running: &types.TaskStateRunning{},
								}
								if probe := t.LivenessProbe; probe != nil {
									liveFunc := func(live bool, err error) {
										if !live {
											_, _ = fmt.Fprintf(stderr, "liveness live=%v err=%v\n", live, err)
											stopRun()
										} else {
											_, _ = fmt.Fprintf(stdout, "liveness live=%v\n", live)
										}
									}
									go probeLoop(runCtx, stopEverything, *probe, liveFunc)
								}
								if probe := t.ReadinessProbe; probe != nil {
									readyFunc := func(ready bool, err error) {
										status.Ready = ready
										if ready {
											_, _ = fmt.Fprintf(stdout, "readiness ready=%v\n", status.Ready)
											maybeStartDownstream(name)
										} else {
											_, _ = fmt.Fprintf(stderr, "readiness ready=%v err=%v\n", ready, err)
										}
									}
									go probeLoop(runCtx, stopEverything, *probe, readyFunc)
								}
								if err := prc.Run(runCtx, stdout, stderr); err != nil {
									return fmt.Errorf("failed to run: %w", err)
								}
								return nil
							}()
							if err != nil {
								if errors.Is(err, context.Canceled) {
									return
								}
								status.State = types.TaskState{
									Terminated: &types.TaskStateTerminated{Reason: "error"},
								}
								_, _ = fmt.Fprintln(stderr, err.Error())
							} else {
								status.State = types.TaskState{
									Terminated: &types.TaskStateTerminated{Reason: "success"},
								}
								maybeStartDownstream(name)
								if !t.IsBackground() {
									return
								}
							}
							time.Sleep(2 * time.Second)
						}
					}
				}(t, statuses.GetStatus(t.Name))

				time.Sleep(time.Second / 4)
			}

			wg.Wait()

			return nil
		},
	}

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func handleCrash(stop func()) {
	if r := recover(); r != nil {
		fmt.Println(r)
		debug.PrintStack()
		stop()
	}
}
