package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/alexec/kit/internal/types"
)

type host struct {
	types.PodSpec
	types.Task
}

func (h *host) Run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.Command(h.Command[0], append(h.Command[1:], h.Args...)...)
	cmd.Dir = h.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = append(os.Environ(), h.Env.Environ()...)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := h.stop(cmd.Process.Pid); err != nil {
			_, _ = fmt.Fprintln(stderr, err.Error())
		}
	}()
	return cmd.Wait()
}

func (h *host) stop(pid int) error {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		// already stopped
		if err.Error() == "no such process" {
			return nil
		}
		return fmt.Errorf("failed get pgid: %w", err)
	}
	if pgid == pid {
		pid = -pid
	}
	target, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}
	if err := target.Signal(syscall.SIGTERM); err == nil {
		return nil
	}
	time.Sleep(h.PodSpec.GetTerminationGracePeriod())
	if err := target.Signal(os.Kill); err == nil {
		return nil
	} else {
		return fmt.Errorf("failed to kill: %w", err)
	}
}

var _ Interface = &host{}
