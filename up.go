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

func up() *cobra.Command {
	return &cobra.Command{
		Use:   "up [config_file]",
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
			init := true

			go func() {
				defer handleCrash(stopEverything)
				defer handleCrash(stopEverything)
				for {
					width, _, _ := terminal.GetSize(0)
					if width == 0 {
						width = 80
					}

					log.Printf("%s[2J", escape)
					log.Printf("%s[H", escape)
					containers := pod.Spec.Containers
					if init {
						containers = pod.Spec.InitContainers
					}
					for _, c := range containers {
						state := status.GetContainerStatus(c.Name)
						if state == nil {
							continue
						}
						icon, reason := "▓", "unknown"
						if state.State.Waiting != nil {
							reason = state.State.Waiting.Reason
						} else if state.State.Running != nil {
							if state.Ready {
								icon, reason = color.GreenString("▓"), "ready"
							} else {
								icon, reason = color.BlueString("▓"), "running"

							}
						} else if state.State.Terminated != nil {
							icon, reason = "▓", state.State.Terminated.Reason
							if reason == "error" {
								icon = color.RedString("▓")
							}
						}

						line := fmt.Sprintf("%s %-10s [%-8s] %v %s", icon, k8sstrings.ShortenString(state.Name, 10), reason, c.GetHostPorts(), logEntries[c.Name].String())
						log.Println(k8sstrings.ShortenString(line, width))
					}
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

			for _, containers := range [][]types.Container{pod.Spec.InitContainers, pod.Spec.Containers} {
				wg := &sync.WaitGroup{}

				if init {
					for _, c := range pod.Spec.InitContainers {
						logEntries[c.Name] = &types.LogEntry{}
						status.InitContainerStatuses = append(status.InitContainerStatuses, types.ContainerStatus{
							Name: c.Name,
						})
					}
				} else {
					for _, c := range pod.Spec.Containers {
						logEntries[c.Name] = &types.LogEntry{}
						status.ContainerStatuses = append(status.ContainerStatuses, types.ContainerStatus{
							Name: c.Name,
						})
					}
				}

				for _, c := range containers {

					name := c.Name
					state := status.GetContainerStatus(name)

					logEntry := logEntries[name]

					logFile, err := os.Create(filepath.Join("logs", name+".log"))
					if err != nil {
						return err
					}
					stdout := io.MultiWriter(logFile, logEntry.Stdout())
					stderr := io.MultiWriter(logFile, logEntry.Stderr())
					prc := proc.New(c, pod.Spec)

					processes[name] = prc

					if err = prc.Init(ctx); err != nil {
						return err
					}

					processCtx, stopProcess := context.WithCancel(ctx)

					wg.Add(1)
					go func(name, image string, livenessProbe, readinessProbe *types.Probe, build *types.Build) {
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
									state.State = types.ContainerState{
										Waiting: &types.ContainerStateWaiting{Reason: "building"},
									}
									if err := prc.Build(runCtx, stdout, stderr); err != nil {
										return fmt.Errorf("failed to build: %v", err)
									}
									state.State = types.ContainerState{
										Running: &types.ContainerStateRunning{},
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
												if build != nil {
													for _, file := range build.Watch {
														stat, err := os.Stat(file)
														if err != nil {
															return
														}
														if stat.ModTime().After(last) {
															stopRun()
														}
													}
												}
												last = next
												time.Sleep(time.Second)
											}
										}
									}()
									if probe := livenessProbe; probe != nil {
										liveFunc := func(name string, live bool, err error) {
											if !live {
												stopRun()
											}
										}
										go probeLoop(runCtx, stopEverything, name, *probe, liveFunc)
									}
									if probe := readinessProbe; probe != nil {
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
									state.State = types.ContainerState{
										Terminated: &types.ContainerStateTerminated{Reason: "error"},
									}
									logEntry.Level = "error"
									logEntry.Msg = err.Error()
								} else {
									state.State = types.ContainerState{
										Terminated: &types.ContainerStateTerminated{Reason: "exited"},
									}
									if init {
										return
									}
								}
								time.Sleep(2 * time.Second)
							}
						}
					}(c.Name, c.Image, c.LivenessProbe.DeepCopy(), c.ReadinessProbe.DeepCopy(), c.Build.DeepCopy())

					time.Sleep(time.Second / 4)
				}

				wg.Wait()

				init = false
			}
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
