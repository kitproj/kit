package main

import (
	"context"
	"fmt"
	"github.com/alexec/kit/internal/proc"
	"github.com/alexec/kit/internal/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"k8s.io/utils/strings/slices"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"sync"
	"syscall"
	"time"
)

func up() *cobra.Command {
	var include []string
	var exclude []string

	cmd := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
			defer stop()

			_ = os.Mkdir("logs", 0777)

			pod := &types.Kit{}

			pod.Status = &types.Status{}

			go func() {
				defer handleCrash(stop)
				for {
					width, _, _ := terminal.GetSize(0)
					if width == 0 {
						width = 80
					}

					log.Printf("%s[2J", escape)
					log.Printf("%s[H", escape)
					for _, state := range pod.GetContainerStatuses() {
						r := "▓"
						if v, ok := map[types.Phase]string{
							types.LivePhase:    color.BlueString("▓"),
							types.ReadyPhase:   color.GreenString("▓"),
							types.UnreadyPhase: color.YellowString("▓"),
						}[state.Phase]; ok {
							r = v
						}
						line := fmt.Sprintf("%s %-10s [%-8s]  %s", r, state.Name, state.Phase, state.Log.String())
						if len(line) > width && width > 0 {
							line = line[0 : width-1]
						}
						log.Println(line)
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

			terminationGracePeriod := pod.Spec.GetTerminationGracePeriod()

			for _, containers := range [][]types.Container{pod.Spec.InitContainers, pod.Spec.Containers} {
				wg := &sync.WaitGroup{}

				if pod.Status.InitContainerStatuses == nil {
					for _, c := range containers {
						pod.Status.InitContainerStatuses = append(pod.Status.InitContainerStatuses, &types.ContainerStatus{Name: c.Name})
					}
				} else {
					for _, c := range containers {
						pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, &types.ContainerStatus{Name: c.Name})
					}
				}

				for _, c := range containers {

					name := c.Name
					state := pod.GetContainerStatuses().Get(c.Name)

					state.Phase = types.DeadPhase

					if slices.Contains(exclude, name) {
						state.Phase = types.ExcludedPhase
						continue
					}

					if include != nil && !slices.Contains(include, name) {
						state.Phase = types.ExcludedPhase
						continue
					}

					state.Phase = types.CreatingPhase

					logFile, err := os.Create(filepath.Join("logs", name+".log"))
					if err != nil {
						return err
					}
					stdout := io.MultiWriter(logFile, state.Stdout())
					stderr := io.MultiWriter(logFile, state.Stderr())
					var pd proc.Proc
					if c.Image == "" {
						pd = &proc.HostProc{Container: c}
					} else {
						pd = &proc.ContainerProc{Container: c}
					}

					if err = pd.Init(ctx); err != nil {
						return err
					}

					wg.Add(1)
					go func() {
						defer handleCrash(stop)
						defer wg.Done()
						for {
							select {
							case <-ctx.Done():
								return
							default:
								err := func() error {
									defer func() { state.Phase = types.ExitedPhase }()
									if err := pd.Stop(ctx, terminationGracePeriod); err != nil {
										return fmt.Errorf("failed to stop: %v", err)
									}
									state.Phase = types.BuildingPhase
									if err := pd.Build(ctx, stdout, stderr); err != nil {
										return fmt.Errorf("failed to build: %v", err)
									}
									state.Phase = types.RunningPhase
									if err := pd.Run(ctx, stdout, stderr); err != nil {
										return fmt.Errorf("failed to run: %v", err)
									}
									return nil
								}()
								if err != nil && err != context.Canceled {
									state.Phase = types.ErrorPhase
									state.Log = types.LogEntry{Level: "error", Msg: err.Error()}
								} else {
									return
								}
							}
						}
					}()

					go func() {
						<-ctx.Done()

						if err := pd.Stop(context.Background(), terminationGracePeriod); err != nil {
							state.Phase = types.ErrorPhase
							state.Log = types.LogEntry{Level: "error", Msg: fmt.Sprintf("failed to stop: %v", err)}
						}
					}()

					if probe := c.LivenessProbe; probe != nil {
						liveFunc := func(live bool, err error) {
							if live {
								state.Phase = types.LivePhase
							} else {
								state.Phase = types.DeadPhase
							}
							if err != nil {
								state.Log = types.LogEntry{Level: "error", Msg: err.Error()}
							}
							if !live {
								if err := pd.Stop(ctx, terminationGracePeriod); err != nil {
									state.Log = types.LogEntry{Level: "error", Msg: err.Error()}
								}
							}
						}
						go probeLoop(ctx, stop, *probe, liveFunc)
					}
					if probe := c.ReadinessProbe; probe != nil {
						readyFunc := func(ready bool, err error) {
							if ready {
								state.Phase = types.ReadyPhase
							} else {
								state.Phase = types.UnreadyPhase
							}
							if err != nil {
								state.Log = types.LogEntry{Level: "error", Msg: err.Error()}
							}
						}
						go probeLoop(ctx, stop, *probe, readyFunc)
					}
				}

				wg.Wait()
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&exclude, "exclude", "e", nil, "exclude")
	cmd.Flags().StringArrayVarP(&include, "include", "i", nil, "include")
	return cmd
}
func handleCrash(stop func()) {
	if r := recover(); r != nil {
		log.Println(r)
		stop()
	}
}
