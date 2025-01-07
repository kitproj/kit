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

	"github.com/kitproj/kit/internal/types"
)

type host struct {
	log  *log.Logger
	spec types.Spec
	types.Task
}

func (h *host) Run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	environ, err := types.Environ(h.spec, h.Task)
	if err != nil {
		return fmt.Errorf("error getting spec environ: %w", err)
	}

	path := h.Command[0]
	cmd := exec.Command(path, append(h.Command[1:], h.Args...)...)
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

var _ Interface = &host{}
