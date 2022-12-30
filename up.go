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
	corev1 "k8s.io/api/core/v1"
	k8sstrings "k8s.io/utils/strings"
	"sigs.k8s.io/yaml"
)

func up() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start-up processes",
		RunE: func(cmd *cobra.Command, args []string) error {

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
						line := fmt.Sprintf("%s %-10s [%-7s]  %s", icon, k8sstrings.ShortenString(state.Name, 10), reason, logEntries[c.Name].String())
						log.Println(k8sstrings.ShortenString(line, width))
					}
					time.Sleep(time.Second / 2)
				}
			}()

			in, err := os.ReadFile(kitFile)
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
			if err = os.WriteFile(kitFile, data, 0o0644); err != nil {
				return err
			}

			for _, containers := range [][]types.Container{pod.Spec.InitContainers, pod.Spec.Containers} {
				wg := &sync.WaitGroup{}

				if init {
					for _, c := range pod.Spec.InitContainers {
						logEntries[c.Name] = &types.LogEntry{}
						status.InitContainerStatuses = append(status.InitContainerStatuses, corev1.ContainerStatus{
							Name: c.Name,
						})
					}
				} else {
					for _, c := range pod.Spec.Containers {
						logEntries[c.Name] = &types.LogEntry{}
						status.ContainerStatuses = append(status.ContainerStatuses, corev1.ContainerStatus{
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
					go func(name, image string, livenessProbe, readinessProbe *types.Probe) {
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
									state.State = corev1.ContainerState{
										Waiting: &corev1.ContainerStateWaiting{Reason: "building"},
									}
									if err := prc.Build(runCtx, stdout, stderr); err != nil {
										return fmt.Errorf("failed to build: %v", err)
									}
									state.State = corev1.ContainerState{
										Running: &corev1.ContainerStateRunning{},
									}
									logEntry.Msg = ""
									go func() {
										defer handleCrash(stopEverything)
										var last time.Time
										for {
											select {
											case <-runCtx.Done():
												return
											default:
												stat, err := os.Stat(image)
												if err != nil {
													return
												}
												if !last.IsZero() && stat.ModTime().After(last) {
													stopRun()
												}
												last = stat.ModTime()
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
									state.State = corev1.ContainerState{
										Terminated: &corev1.ContainerStateTerminated{Reason: "error"},
									}
									logEntry.Level = "error"
									logEntry.Msg = err.Error()
								} else {
									state.State = corev1.ContainerState{
										Terminated: &corev1.ContainerStateTerminated{Reason: "exited"},
									}
									if init {
										return
									}
								}
							}
						}
					}(c.Name, c.Image, c.LivenessProbe.DeepCopy(), c.ReadinessProbe)

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
