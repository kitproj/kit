package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/alexec/kit/internal/types"
)

type host struct {
	types.Container
	grace   time.Duration
	process *os.Process
}

func (h *host) Init(ctx context.Context) error {
	return nil
}

func (h *host) Build(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if f, ok := imageIsHostfile(h.Image); ok {
		cmd := exec.CommandContext(ctx, f)
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
		go func() {
			<-ctx.Done()
			h.stop()
		}()
		h.process = cmd.Process
		return cmd.Wait()
	}
	return nil
}

func (h *host) Run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(ctx, h.Command[0], append(h.Command[1:], h.Args...)...)
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
	go func() {
		<-ctx.Done()
		h.stop()
	}()
	return cmd.Wait()
}

func (h *host) stop() error {
	if h.process != nil {
		pgid, _ := syscall.Getpgid(h.process.Pid)
		if err := syscall.Kill(-pgid, syscall.SIGTERM); err == nil || isNotPermitted(err) {
			return nil
		}
		time.Sleep(h.grace)
		if err := syscall.Kill(-pgid, syscall.SIGKILL); err == nil || isNotPermitted(err) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func isNotPermitted(err error) bool {
	return err != nil && err.Error() == "operation not permitted"
}

var _ Interface = &host{}

const hostfile = "Hostfile"

func imageIsHostfile(image string) (string, bool) {
	f := filepath.Join(image, hostfile)
	_, err := os.Stat(f)
	return f, err == nil
}