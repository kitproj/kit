package proc

import (
	"context"
	"fmt"
	"github.com/alexec/joy/internal/types"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type HostProc struct {
	types.Container
	process *os.Process
}

func (h *HostProc) Init(ctx context.Context) error {
	return nil
}

func (h *HostProc) Build(ctx context.Context, stdout, stderr io.Writer) error {
	return nil
}

func (h *HostProc) Run(ctx context.Context, stdout, stderr io.Writer) error {
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

func (h *HostProc) Stop(ctx context.Context, grace time.Duration) error {
	if h.process != nil {
		pgid, _ := syscall.Getpgid(h.process.Pid)
		if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil && !isNotPermitted(err) {
			return err
		}
		time.Sleep(grace)
		if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil && !isNotPermitted(err) {
			return err
		}
		return nil
	}
	return nil
}

func isNotPermitted(err error) bool {
	return err.Error() == "operation not permitted"
}

var _ Proc = &HostProc{}
