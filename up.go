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

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kit [config_file]",
		Short: "Start-up processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := defaultConfigFile

			ctx, stopEverything := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
			defer stopEverything()

			_ = os.Mkdir("logs", 0777)

			pod := &types.Pod{}
			status := &types.Status{}
			logEntries := make(map[string]*types.LogEntry)
			stoppers := make(map[string]func())
			processes := make(map[string]proc.Interface)

			go func() {
				defer handleCrash(stopEverything)
				for {
					width, _, _ := terminal.GetSize(0)
					if width == 0 {
						width = 80
					}
					log.Printf("%s[2J", escape)   // clear screen
					log.Printf("%s[0;0H", escape) // move to 0,0
					for _, t := range pod.Spec.Tasks {
						state := status.GetContainerStatus(t.Name)
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
						log.Println(prefix + " " + msg)
					}
					log.Println()
					log.Printf("kit %s", tag)
					time.Sleep(time.Second / 2)
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

			for t := range tasks {
				name := t.Name
				logEntries[name] = &types.LogEntry{}
				if status.GetContainerStatus(name) == nil {
					status.TaskStatuses = append(status.TaskStatuses, &types.TaskStatus{
						Name: name,
					})
				}

				logEntry := logEntries[name]

				prc := proc.New(t, pod.Spec)

				processes[name] = prc

				processCtx, stopProcess := context.WithCancel(ctx)

				stoppers[name] = stopProcess

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
								}
							}
							last = time.Now()
							time.Sleep(time.Second)
						}
					}
				}(t, stopProcess)
				go func(t types.Task, status *types.TaskStatus) {
					defer handleCrash(stopEverything)
					wg.Add(1)
					defer wg.Done()
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
											for _, downstream := range pod.Spec.Tasks.GetDownstream(t.Name) {
												tasks <- downstream
											}
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
								for _, downstream := range pod.Spec.Tasks.GetDownstream(t.Name) {
									if stop, ok := stoppers[downstream.Name]; ok {
										stop()
										delete(stoppers, downstream.Name)
									}
									tasks <- downstream
								}
								if !t.IsBackground() {
									return
								}
							}
							time.Sleep(2 * time.Second)
						}
					}
				}(t, status.GetContainerStatus(t.Name))

				time.Sleep(time.Second / 4)
			}

			wg.Wait()

			return nil
		},
	}
}

func handleCrash(stop func()) {
	if r := recover(); r != nil {
		log.Println(r)
		debug.PrintStack()
		stop()
	}
}
