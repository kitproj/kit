package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/alexec/kit/internal/proc"
	"github.com/alexec/kit/internal/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"k8s.io/utils/strings/slices"
	"sigs.k8s.io/yaml"
)

func up() *cobra.Command {
	var include []string
	var exclude []string

	cmd := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
			defer stop()

			var states = map[string]*types.State{}

			in, err := os.ReadFile(kitFile)
			if err != nil {
				return err
			}
			pod := &types.Kit{}
			err = yaml.UnmarshalStrict(in, pod)
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(pod)
			if err != nil {
				return err
			}
			err = os.WriteFile(kitFile, data, 0o0644)
			if err != nil {
				return err
			}

			var names []string
			go func() {
				defer handleCrash(stop)
				for {
					width, _, _ := terminal.GetSize(0)
					if width == 0 {
						width = 80
					}

					log.Printf("%s[2J", escape)
					log.Printf("%s[H", escape)
					for _, name := range names {
						state := states[name]
						r := "▓"
						if v, ok := map[types.Phase]string{
							types.LivePhase:    color.BlueString("▓"),
							types.ReadyPhase:   color.GreenString("▓"),
							types.UnreadyPhase: color.YellowString("▓"),
						}[state.Phase]; ok {
							r = v
						}
						line := fmt.Sprintf("%s %-10s [%-8s]  %s", r, name, state.Phase, state.Log.String())
						if len(line) > width && width > 0 {
							line = line[0 : width-1]
						}
						log.Println(line)
					}
					time.Sleep(time.Second / 2)
				}
			}()

			_ = os.Mkdir("logs", 0777)

			terminationGracePeriod := pod.Spec.GetTerminationGracePeriod()

			for _, containers := range [][]types.Container{pod.Spec.InitContainers, pod.Spec.Containers} {
				wg := &sync.WaitGroup{}

				states = map[string]*types.State{}
				names = nil

				for _, c := range containers {
					if _, ok := states[c.Name]; ok {
						return fmt.Errorf("duplicate name %s", c.Name)
					}
					states[c.Name] = &types.State{}
					names = append(names, c.Name)
				}

				for _, c := range containers {
					name := c.Name
					state := states[name]

					if slices.Contains(exclude, name) {
						state.Phase = types.ExcludedPhase
						continue
					}

					if include != nil && !slices.Contains(include, name) {
						state.Phase = types.ExcludedPhase
						continue
					}

					state.Phase = types.CreatingPhase

					var pd proc.Proc

					logFile, err := os.Create(filepath.Join("logs", name+".log"))
					if err != nil {
						if err != nil {
							return err
						}
					}
					stdout := io.MultiWriter(logFile, states[c.Name].Stdout())
					stderr := io.MultiWriter(logFile, states[c.Name].Stderr())
					if c.Image == "" {
						pd = &proc.HostProc{Container: c}
					} else {
						pd = &proc.ContainerProc{Container: c}
					}

					err = pd.Init(ctx)
					if err != nil {
						return err
					}

					go func(state *types.State) {
						defer handleCrash(stop)
						wg.Add(1)
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
								time.Sleep(3 * time.Second)
							}
						}
					}(state)

					go func() {
						select {
						case <-ctx.Done():
							err := pd.Stop(context.Background(), terminationGracePeriod)
							if err != nil {
								state.Phase = types.ErrorPhase
								state.Log = types.LogEntry{Level: "error", Msg: fmt.Sprintf("failed to stop: %v", err)}
							}
						}
					}()

					if probe := c.LivenessProbe; probe != nil {
						liveFunc := func(name string, live bool, err error) {
							if live {
								states[name].Phase = types.LivePhase
							} else {
								states[name].Phase = types.DeadPhase
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
						go probeLoop(ctx, stop, name, *probe, liveFunc)
					}
					if probe := c.ReadinessProbe; probe != nil {
						readyFunc := func(name string, ready bool, err error) {
							if ready {
								states[name].Phase = types.ReadyPhase
							} else {
								states[name].Phase = types.UnreadyPhase
							}
							if err != nil {
								state.Log = types.LogEntry{Level: "error", Msg: err.Error()}
							}
						}
						go probeLoop(ctx, stop, name, *probe, readyFunc)
					}
					time.Sleep(time.Second)
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
