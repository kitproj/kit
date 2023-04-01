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
	types.PodSpec
	types.Task
}

func (h *host) Run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	path := h.Command[0]
	cmd := exec.Command(path, append(h.Command[1:], h.Args...)...)
	cmd.Dir = h.WorkingDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = append(os.Environ(), h.Env.Environ()...)
	log.Printf("%s: starting process %q\n", h.Name, h.Command)
	err := cmd.Start()
	log.Printf("%s: started process %q: %v\n", h.Name, h.Command, err)
	if err != nil {
		return err
	}
	// capture pgid straight away because it's not available after the process exits,
	// the process may exit and leave children behind.
	pid := cmd.Process.Pid
	log.Printf("%s: getting pgid for %d\n", h.Name, pid)
	pgid, err := syscall.Getpgid(pid)
	log.Printf("%s: pgid for %d is %d %v\n", h.Name, pid, pgid, err)
	if err != nil {
		return fmt.Errorf("failed get pgid: %w", err)
	}
	go func() {
		<-ctx.Done()
		log.Printf("%s: context cancelled, stopping process", h.Name)
		if err := h.stop(pgid); err != nil {
			_, _ = fmt.Fprintln(stderr, err.Error())
		}
	}()
	log.Printf("%s: waiting for process %d pgid %d (%q)", h.Name, pid, pgid, h.Command)
	err = cmd.Wait()
	log.Printf("%s: process exited %d: %v", h.Name, pid, err)
	return err
}

func (h *host) stop(pid int) error {
	log.Printf("%s: stopping process %d", h.Name, pid)
	log.Printf("%s: finding process %d\n", h.Name, pid)
	target, err := os.FindProcess(-pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}
	log.Printf("%s: terminating process %d\n", h.Name, pid)
	if err := target.Signal(syscall.SIGTERM); ignoreProcessFinishedErr(err) != nil {
		log.Printf("%s: failed to terminate: %v", h.Name, err)
	}
	gracePeriod := h.PodSpec.GetTerminationGracePeriod()
	log.Printf("%s: waiting %v before killing %d\n", h.Name, gracePeriod, pid)
	time.Sleep(gracePeriod)
	log.Printf("%s: killing process %d\n", h.Name, pid)
	err = target.Signal(os.Kill)
	log.Printf("%s: killed process %d: %v\n", h.Name, pid, err)
	if ignoreProcessFinishedErr(err) != nil {
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
