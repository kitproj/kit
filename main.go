package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alexec/kit/internal/proc"
	"github.com/alexec/kit/internal/types"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"sigs.k8s.io/yaml"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const escape = "\x1b"

var states = map[string]*types.State{}

func main() {

	include := ""
	exclude := ""
	flag.StringVar(&include, "i", "", "include")
	flag.StringVar(&exclude, "e", "", "exclude")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	in, err := os.ReadFile("kit.yaml")
	must(err)
	pod := &types.Kit{}
	must(yaml.UnmarshalStrict(in, pod))

	var names []string
	go func() {
		defer handleCrash(stop)
		for {
			width, _, err := terminal.GetSize(0)
			must(err)

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
				must(fmt.Errorf("duplicate name %s", c.Name))
			}
			states[c.Name] = &types.State{}
			names = append(names, c.Name)
		}

		for _, c := range containers {
			name := c.Name
			state := states[name]

			if strings.Contains(","+exclude+",", name) {
				state.Phase = types.ExcludedPhase
				continue
			}

			if include != "" && !strings.Contains(","+include+",", name) {
				state.Phase = types.ExcludedPhase
				continue
			}

			state.Phase = types.CreatingPhase

			var pd proc.Proc

			logFile, err := os.Create(filepath.Join("logs", name+".log"))
			if err != nil {
				must(err)
			}
			stdout := io.MultiWriter(logFile, states[c.Name].Stdout())
			stderr := io.MultiWriter(logFile, states[c.Name].Stderr())
			if c.Image == "" {
				pd = &proc.HostProc{Container: c}
			} else {
				pd = &proc.ContainerProc{Container: c}
			}

			err = pd.Init(ctx)
			must(err)

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
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleCrash(stop func()) {
	if r := recover(); r != nil {
		log.Println(r)
		stop()
	}
}
