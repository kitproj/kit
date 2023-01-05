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

			wg := &sync.WaitGroup{}

			for _, t := range pod.Spec.Tasks {
				name := t.Name
				logEntries[name] = &types.LogEntry{}
				status.TaskStatuses = append(status.TaskStatuses, types.TaskStatus{
					Name: name,
				})

				state := status.GetContainerStatus(name)

				logEntry := logEntries[name]

				logFile, err := os.Create(filepath.Join("logs", name+".log"))
				if err != nil {
					return err
				}
				stdout := io.MultiWriter(logFile, logEntry.Stdout())
				stderr := io.MultiWriter(logFile, logEntry.Stderr())
				prc := proc.New(t, pod.Spec)

				processes[name] = prc

				processCtx, stopProcess := context.WithCancel(ctx)

				wg.Add(1)
				go func(t types.Task) {
					defer handleCrash(stopEverything)
					defer wg.Done()
					for {
						select {
						case <-ctx.Done():
							return
						default:
							err := func() error {
								runCtx, stopRun := context.WithCancel(processCtx)
								defer stopRun()
								state.State = types.TaskState{
									Running: &types.TaskStateRunning{},
								}
								logEntry.Msg = ""
								go func() {
									defer handleCrash(stopEverything)
									last := time.Now()
									for {
										select {
										case <-runCtx.Done():
											return
										default:
											next := time.Now()

											for _, file := range t.Watch {
												stat, err := os.Stat(file)
												if err != nil {
													return
												}
												if stat.ModTime().After(last) {
													stopRun()
												}
											}

											last = next
											time.Sleep(time.Second)
										}
									}
								}()
								if probe := t.LivenessProbe; probe != nil {
									liveFunc := func(name string, live bool, err error) {
										if !live {
											stopRun()
										}
									}
									go probeLoop(runCtx, stopEverything, name, *probe, liveFunc)
								}
								if probe := t.ReadinessProbe; probe != nil {
									readyFunc := func(name string, ready bool, err error) {
										state.Ready = ready
										if err != nil {
											logEntry.Level = "error"
											logEntry.Msg = err.Error()

										}
									}
									go probeLoop(ctx, stopEverything, name, *probe, readyFunc)
								}

								go func() {
									defer handleCrash(stopEverything)
									<-ctx.Done()
									stopProcess()
								}()
								if err := prc.Run(runCtx, stdout, stderr); err != nil {
									return fmt.Errorf("failed to run: %v", err)
								}
								return nil
							}()
							if err != nil {
								if errors.Is(err, context.Canceled) {
									return
								}
								state.State = types.TaskState{
									Terminated: &types.TaskStateTerminated{Reason: "error"},
								}
								logEntry.Level = "error"
								logEntry.Msg = err.Error()
							} else {
								state.State = types.TaskState{
									Terminated: &types.TaskStateTerminated{Reason: "exited"},
								}
								if !t.IsBackground() {
									return
								}
							}
							time.Sleep(2 * time.Second)
						}
					}
				}(*t.DeepCopy())

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
