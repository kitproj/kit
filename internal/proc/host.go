package proc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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

	path := h.Command[0]
	cmd := exec.Command(path, append(h.Command[1:], h.Args...)...)
	cmd.Dir = h.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = append(os.Environ(), h.Env.Environ()...)
	log.Printf("starting process %q", h.Command)
	if err := cmd.Start(); err != nil {
		return err
	}
	// capture pgid straight away because it's not available after the process exits,
	// the process may exit and leave children behind.
	pid := cmd.Process.Pid
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return fmt.Errorf("failed get pgid: %w", err)
	}
	go func() {
		<-ctx.Done()
		if err := h.stop(pgid); err != nil {
			_, _ = fmt.Fprintln(stderr, err.Error())
		}
	}()
	log.Printf("waiting for process %d(%q)", pid, h.Command)
	err = cmd.Wait()
	log.Printf("process %d exited : %v", pid, err)
	return err
}

func (h *host) stop(pid int) error {
	log.Printf("terminating process %d\n", pid)
	target, err := os.FindProcess(-pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}
	if err := target.Signal(syscall.SIGTERM); err == nil {
		return nil
	}
	time.Sleep(h.PodSpec.GetTerminationGracePeriod())
	log.Printf("killing process %d\n", pid)
	if err := target.Signal(os.Kill); ignoreProcessFinishedErr(err) != nil {
		return fmt.Errorf("failed to kill: %w", err)
	}
	return nil
}

func (h *host) Reset(ctx context.Context) error {
	return nil
}

func ignoreProcessFinishedErr(err error) error {
	if err != nil && !strings.Contains(err.Error(), "process already finished") {
		return err
	}
	return nil
}

var _ Interface = &host{}
