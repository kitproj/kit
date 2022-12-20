package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	network "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
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

	cli, err := client.NewClientWithOpts(client.FromEnv)
	must(err)
	defer cli.Close()

	for _, c := range pod.Spec.Containers {
		if _, ok := states[c.Name]; ok {
			must(fmt.Errorf("duplicate name %s", c.Name))
		}
		states[c.Name] = &state{
			kill: func() error {
				return nil
			},
		}
	}

	width, _, err := terminal.GetSize(0)
	must(err)

	go func() {
		defer runtime.HandleCrash()
		for {
			log.Printf("%s[2J", escape)
			log.Printf("%s[H", escape)
			for _, c := range pod.Spec.Containers {
				name := c.Name
				state := states[name]
				r := "▓"
				if v, ok := map[phase]string{
					readyPhase:   color.GreenString("▓"),
					unreadyPhase: color.YellowString("▓"),
				}[state.phase]; ok {
					r = v
				}
				line := fmt.Sprintf("%s %-10s [%-8s] %s | %s", r, name, state.phase, color.BlueString(state.msg), state.log.String())
				if len(line) > width {
					line = line[0 : width-1]
				}
				log.Println(line)
			}
			time.Sleep(time.Second / 2)
		}
	}()

	defer func() {
		for name, state := range states {
			states[name].phase = killingPhase
			defer func() { states[name].phase = killedPhase }()
			if err := state.kill(); err != nil {
				states[name].msg = err.Error()
			}
		}
		time.Sleep(time.Second)
	}()

	_ = os.Mkdir("logs", 0777)

	for _, c := range pod.Spec.Containers {
		states[c.Name].phase = creatingPhase

		go func(c corev1.Container) {
			err := func() error {
				logFile, err := os.Create(filepath.Join("logs", c.Name+".log"))
				if err != nil {
					return err
				}

				var environ []string
				for _, e := range c.Env {
					environ = append(environ, fmt.Sprintf("%s=%s", e.Name, e.Value))
				}
				stdout := io.MultiWriter(logFile, states[c.Name].Stdout())
				stderr := io.MultiWriter(logFile, states[c.Name].Stderr())
				image := c.Image
				if image == "" {
					cmd := exec.Command(c.Command[0], append(c.Command[1:], c.Args...)...)
					cmd.Dir = c.WorkingDir
					cmd.Stdin = os.Stdin
					cmd.Stdout = stdout
					cmd.Stderr = stderr
					cmd.SysProcAttr = &syscall.SysProcAttr{
						Setpgid: true,
					}
					cmd.Env = append(os.Environ(), environ...)

					states[c.Name].kill = func() error {
						if cmd.Process != nil {
							pgid, _ := syscall.Getpgid(cmd.Process.Pid)
							err := syscall.Kill(-pgid, syscall.SIGTERM)
							if err != nil {
								return err
							}
						}
						return logFile.Close()
					}

					go func(name string, cmd *exec.Cmd) {
						defer runtime.HandleCrash()
						defer func() {
							states[c.Name].phase = exitedPhase
						}()

						states[c.Name].phase = runningPhase
						if err := cmd.Run(); err != nil {
							states[name].phase = errorPhase
							states[name].msg = err.Error()
						}
					}(c.Name, cmd)
				} else {
					_, err := os.Stat(image)
					if err == nil {
						states[c.Name].phase = buildingPhase
						r, err := archive.TarWithOptions(filepath.Dir(image), &archive.TarOptions{})
						if err != nil {
							return err
						}
						defer r.Close()
						resp, err := cli.ImageBuild(ctx, r, types.ImageBuildOptions{Dockerfile: filepath.Base(image), Tags: []string{c.Name}})
						if err != nil {
							return err
						}
						defer resp.Body.Close()
						_, err = io.Copy(stdout, resp.Body)
						if err != nil {
							return err
						}
						image = c.Name
					} else if c.ImagePullPolicy != corev1.PullNever {
						states[c.Name].phase = pullingPhase
						r, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
						if err != nil {
							return err
						}
						_, err = io.Copy(states[c.Name].Stdout(), r)
						if err != nil {
							return err
						}
						err = r.Close()
						if err != nil {
							return err
						}
					}
					states[c.Name].phase = creatingPhase
					portSet := nat.PortSet{}
					portBindings := map[nat.Port][]nat.PortBinding{}
					for _, p := range c.Ports {
						port, err := nat.NewPort("tcp", fmt.Sprint(p.ContainerPort))
						if err != nil {
							return err
						}
						portSet[port] = struct{}{}
						hostPort := p.HostPort
						if hostPort == 0 {
							hostPort = p.ContainerPort
						}
						portBindings[port] = []nat.PortBinding{{
							HostPort: fmt.Sprint(hostPort),
						}}
					}
					list, err := cli.ContainerList(ctx, types.ContainerListOptions{
						All: true,
					})
					if err != nil {
						return err
					}
					for _, existing := range list {
						if existing.Labels["name"] == c.Name {
							err := cli.ContainerRemove(ctx, existing.ID, types.ContainerRemoveOptions{Force: true})
							if err != nil {
								return err
							}
						}
					}

					created, err := cli.ContainerCreate(ctx, &container.Config{
						Hostname: c.Name,
						Env:      environ,
						// TODO support entrypoint
						Entrypoint:   c.Command,
						Cmd:          c.Args,
						Image:        image,
						WorkingDir:   c.WorkingDir,
						Tty:          c.TTY,
						ExposedPorts: portSet,
						Labels:       map[string]string{"name": c.Name},
					}, &container.HostConfig{
						PortBindings: portBindings,
					}, &network.NetworkingConfig{}, &v1.Platform{}, c.Name)
					if err != nil {
						return err
					}
					states[c.Name].phase = startingPhase
					err = cli.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
					if err != nil {
						return err
					}
					states[c.Name].kill = func() error {
						timeout := 3 * time.Second
						err := cli.ContainerStop(context.Background(), created.ID, &timeout)
						if err != nil {
							return err
						}
						return logFile.Close()
					}
					logs, err := cli.ContainerLogs(ctx, c.Name, types.ContainerLogsOptions{
						ShowStdout: true,
						ShowStderr: true,
						Follow:     true,
					})
					if err != nil {
						return err
					}
					go func(name string, closer io.ReadCloser) {
						defer runtime.HandleCrash()
						defer func() {
							states[c.Name].phase = exitedPhase
						}()
						defer logs.Close()
						states[c.Name].phase = runningPhase
						_, err := io.Copy(stdout, logs)
						if err != nil {
							states[c.Name].phase = errorPhase
							states[name].msg = err.Error()
						}
					}(c.Name, logs)
				}

				if c.ReadinessProbe != nil {
					go func(name string, probe *corev1.Probe) {
						defer runtime.HandleCrash()
						initialDelay := time.Duration(probe.InitialDelaySeconds) * time.Second
						period := time.Duration(probe.PeriodSeconds) * time.Second
						if period == 0 {
							period = 10 * time.Second
						}
						time.Sleep(initialDelay)
						for {
							if states[name].phase.isTerminal() {
								return
							}
							if tcp := probe.TCPSocket; tcp != nil {
								_, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", tcp.Port.IntVal))
								if err != nil {
									states[name].phase = unreadyPhase
									states[name].msg = err.Error()
								} else {
									states[name].phase = readyPhase
									states[name].msg = ""
								}
							} else if httpGet := probe.HTTPGet; httpGet != nil {
								proto := strings.ToLower(string(httpGet.Scheme))
								if proto == "" {
									proto = "http"
								}
								resp, err := http.Get(fmt.Sprintf("%s://localhost:%v%s", proto, httpGet.Port.IntValue(), httpGet.Path))
								if err != nil {
									states[name].phase = unreadyPhase
									states[name].msg = err.Error()
								} else if resp.StatusCode < 300 {
									states[name].phase = readyPhase
									states[name].msg = ""
								} else {
									states[name].phase = unreadyPhase
									states[name].msg = resp.Status
								}
							} else {
								states[c.Name].log = logEntry{"error", "probe not supported"}
							}
							time.Sleep(period)
						}
					}(c.Name, c.ReadinessProbe)
				}
				return nil
			}()
			if err != nil {
				states[c.Name].phase = errorPhase
				states[c.Name].msg = err.Error()
			}
		}(c)
	}

	<-ctx.Done()

	time.Sleep(3 * time.Second)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
