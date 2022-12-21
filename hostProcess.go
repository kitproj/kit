package main

import (
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"os"
	"os/exec"
	"syscall"
)

type HostProcess struct {
	corev1.Container
	process *os.Process
}

func (h *HostProcess) Init(ctx context.Context) error {
	return nil
}

func (h *HostProcess) Build(ctx context.Context, stdout, stderr io.Writer) error {
	return nil
}

func (h *HostProcess) Run(ctx context.Context, stdout, stderr io.Writer) error {
	cmd := exec.Command(h.Command[0], append(h.Command[1:], h.Args...)...)
	cmd.Dir = h.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	for _, env := range h.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	h.process = cmd.Process
	return cmd.Wait()
}

func (h *HostProcess) Kill(ctx context.Context) error {
	defer func() { h.process = nil }()
	if h.process != nil {
		pgid, _ := syscall.Getpgid(h.process.Pid)
		return syscall.Kill(-pgid, syscall.SIGTERM)
	}
	return nil
}

var _ ProcessDef = &HostProcess{}
