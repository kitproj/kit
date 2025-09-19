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

	"github.com/kitproj/kit/internal/metrics"
	"github.com/kitproj/kit/internal/types"
)

type host struct {
	log  *log.Logger
	spec types.Spec
	types.Task
	pid            int
	profFSSnapshot *metrics.ProcFSSnapshot
}

func (h *host) Run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	environ, err := types.Environ(h.spec, h.Task)
	if err != nil {
		return fmt.Errorf("error getting spec environ: %w", err)
	}

	command := h.GetCommand()
	path := command[0]
	cmd := exec.CommandContext(ctx, path, append(command[1:], h.Args...)...)
	cmd.Dir = h.WorkingDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = append(environ, os.Environ()...)
	log := h.log
	log.Println("starting process")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	// capture pgid straight away because it's not available after the process exits,
	// the process may exit and leave children behind.
	pid := cmd.Process.Pid
	h.pid = pid
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return fmt.Errorf("failed get pgid: %w", err)
	}
	go func() {
		<-ctx.Done()
		if err := h.stop(pgid); err != nil {
			log.Printf("failed to stop process: %v", err)
		}
	}()
	return cmd.Wait()
}

func (h *host) stop(pid int) error {
	target, err := os.FindProcess(-pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}
	log := h.log
	if err := target.Signal(syscall.SIGTERM); ignoreProcessFinishedErr(err) != nil {
		log.Printf("failed to terminate: %v", err)
	}
	gracePeriod := h.spec.GetTerminationGracePeriod()
	time.Sleep(gracePeriod)
	err = target.Signal(os.Kill)
	if ignoreProcessFinishedErr(err) != nil {
		return fmt.Errorf("failed to kill: %w", err)
	}
	return nil
}

func ignoreProcessFinishedErr(err error) error {
	if err != nil && !strings.Contains(err.Error(), "process already finished") {
		return err
	}
	return nil
}

func (h *host) GetMetrics(ctx context.Context) (*types.Metrics, error) {
	command := metrics.GetPSCommand(h.pid)
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get process metrics for pid %d: %w", h.pid, err)
	}

	return metrics.ParsePSOutput(string(output))
}
