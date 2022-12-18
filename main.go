package main

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"sigs.k8s.io/yaml"
	"strings"
	"syscall"
	"time"
)

type event = any

type signalEvent struct{}

type processExitedEvent struct {
	name string
	err  error
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	pwd, err := os.Getwd()
	ok(err)

	log.Printf("pwd=%q", pwd)

	in, err := os.ReadFile("dev.yaml")
	ok(err)
	pod := &v1.Pod{}
	ok(yaml.UnmarshalStrict(in, pod))

	hosts, err := os.Create("hosts")
	ok(err)
	for _, c := range pod.Spec.Containers {
		_, err := hosts.WriteString(fmt.Sprintf("127.0.0.1 %s\n", c.Name))
		ok(err)
	}
	ok(hosts.Close())

	events := make(chan event)
	var cmds []*exec.Cmd

	for _, c := range pod.Spec.Containers {
		cmd := exec.Command(c.Command[0], append(c.Command[1:], c.Args...)...)
		cmd.Dir = c.WorkingDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true, // attach to myself
		}
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOSTALIASES=hosts"))

		for _, e := range c.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", e.Name, e.Value))
		}

		log.Printf("name=%s path=%s args=%q", c.Name, cmd.Path, cmd.Args)

		cmds = append(cmds, cmd)

		go func(name string, cmd *exec.Cmd) {
			err := cmd.Run()
			events <- processExitedEvent{name, err}
		}(c.Name, cmd)

		log.Printf("readinessProbe=%v", c.ReadinessProbe != nil)

		if c.ReadinessProbe != nil {
			go func(name string, probe *v1.Probe) {
				initialDelay := time.Duration(probe.InitialDelaySeconds) * time.Second
				period := time.Duration(probe.PeriodSeconds) * time.Second
				if period == 0 {
					period = 10 * time.Second
				}
				log.Printf("name=%s initialDelay=%v, period=%v", name, initialDelay, period)
				time.Sleep(initialDelay)
				for {
					if httpGet := probe.HTTPGet; httpGet != nil {
						proto := strings.ToLower(string(httpGet.Scheme))
						if proto == "" {
							proto = "http"
						}
						resp, err := http.Get(fmt.Sprintf("%s://localhost:%v%s", proto, httpGet.Port.IntValue(), httpGet.Path))
						ready := resp != nil && resp.StatusCode == http.StatusOK
						log.Printf("name=%s ready=%v err=%v", name, ready, err)
					} else {
						log.Fatalf("only httpGet supported: %s", c.Name)
					}
					time.Sleep(period)
				}
			}(c.Name, c.ReadinessProbe)
		}
	}

	go func() {
		<-ctx.Done()
		events <- signalEvent{}
	}()

	log.Printf("running...")

	waitingFor := len(pod.Spec.Containers)

	for event := range events {
		switch obj := event.(type) {
		case signalEvent:
			log.Printf("exiting...")
			for _, cmd := range cmds {
				_ = cmd.Process.Kill()
			}
		case processExitedEvent:
			log.Printf("name=%s err=%q", obj.name, obj.err)
			waitingFor--
			if waitingFor == 0 {
				return
			}
		}
	}
}

func ok(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}
