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
	if h.Task.HasMutex() {
		if _, err := stdout.Write([]byte(fmt.Sprintf("waiting for mutex %q to unlock...\n", h.Mutex))); err != nil {
			return err
		}
		mutex := KeyLock(h.Mutex)
		mutex.Lock()
		defer mutex.Unlock()
		if _, err := stdout.Write([]byte(fmt.Sprintf("locked mutex %q\n", h.Mutex))); err != nil {
			return err
		}
	}

	cmd := exec.CommandContext(ctx, h.Command[0], append(h.Command[1:], h.Args...)...)
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
			_, _ = stderr.Write([]byte(err.Error()))
		}
	}()
	return cmd.Wait()
}

func (h *host) stop(pid int) error {
	pgid, _ := syscall.Getpgid(pid)
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err == nil || isNotPermitted(err) {
		return nil
	}
	time.Sleep(h.PodSpec.GetTerminationGracePeriod())
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err == nil || isNotPermitted(err) {
		return nil
	} else {
		return err
	}
}

func isNotPermitted(err error) bool {
	return err != nil && err.Error() == "operation not permitted"
}

var _ Interface = &host{}
