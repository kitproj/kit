package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"syscall"
	"time"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const escape = "\x1b"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	in, err := os.ReadFile("dev.yaml")
	must(err)
	pod := &corev1.Pod{}
	must(yaml.UnmarshalStrict(in, pod))

	for _, c := range pod.Spec.Containers {
		if _, ok := states[c.Name]; ok {
			must(fmt.Errorf("duplicate name %s", c.Name))
		}
		states[c.Name] = &State{}
	}

	go func() {
		defer runtime.HandleCrash()
		for {
			width, _, err := terminal.GetSize(0)
			must(err)

			log.Printf("%s[2J", escape)
			log.Printf("%s[H", escape)
			for _, c := range pod.Spec.Containers {
				name := c.Name
				state := states[name]
				r := "▓"
				if v, ok := map[phase]string{
					livePhase:    color.BlueString("▓"),
					readyPhase:   color.GreenString("▓"),
					unreadyPhase: color.YellowString("▓"),
				}[state.phase]; ok {
					r = v
				}
				line := fmt.Sprintf("%s %-10s [%-8s]  %s", r, name, state.phase, state.log.String())
				if len(line) > width && width > 0 {
					line = line[0 : width-1]
				}
				log.Println(line)
			}
			time.Sleep(time.Second / 2)
		}
	}()

	_ = os.Mkdir("logs", 0777)

	for _, c := range pod.Spec.Containers {
		name := c.Name
		state := states[name]
		state.phase = creatingPhase

		var pd processDef

		logFile, err := os.Create(filepath.Join("logs", name+".log"))
		if err != nil {
			must(err)
		}
		stdout := io.MultiWriter(logFile, states[c.Name].Stdout())
		stderr := io.MultiWriter(logFile, states[c.Name].Stderr())
		if c.Image == "" {
			pd = &hostProcess{Container: *c.DeepCopy()}
		} else {
			pd = &containerProcess{Container: *c.DeepCopy()}
		}

		err = pd.Init(ctx)
		must(err)

		go func(state *State) {
			defer runtime.HandleCrash()
			for {
				select {
				case <-ctx.Done():
					_ = pd.Kill(context.Background())
					return
				default:
					err := func() error {
						defer func() { state.phase = exitedPhase }()
						if err := pd.Kill(ctx); err != nil {
							return err
						}
						state.phase = buildingPhase
						if err := pd.Build(ctx, stdout, stderr); err != nil {
							return err
						}
						state.phase = runningPhase
						return pd.Run(ctx, stdout, stderr)
					}()
					if err != nil {
						state.phase = errorPhase
						state.log = logEntry{"error", err.Error()}
					}
					time.Sleep(3 * time.Second)
				}
			}
		}(state)

		if p := c.LivenessProbe; p != nil {
			liveFunc := func(name string, live bool, err error) {
				if live {
					states[name].phase = livePhase
				} else {
					states[name].phase = deadPhase
				}
				if err != nil {
					state.log = logEntry{"error", err.Error()}
				}
				if !live {
					if err := pd.Kill(ctx); err != nil {
						state.log = logEntry{"error", err.Error()}
					}
				}
			}
			go probeLoop(ctx, name, *p.DeepCopy(), liveFunc)
		}
		if probe := c.ReadinessProbe; probe != nil {
			readyFunc := func(name string, ready bool, err error) {
				if ready {
					states[name].phase = readyPhase
				} else {
					states[name].phase = unreadyPhase
				}
				if err != nil {
					state.log = logEntry{"error", err.Error()}
				}
			}
			go probeLoop(ctx, name, *probe.DeepCopy(), readyFunc)
		}
	}

	<-ctx.Done()

	time.Sleep(3 * time.Second)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
