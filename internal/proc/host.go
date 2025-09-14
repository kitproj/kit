package proc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kitproj/kit/internal/types"
)

type host struct {
	log  *log.Logger
	spec types.Spec
	types.Task
	pid int
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

	// The `ps` command is used to get process information.
	// -o %cpu,%mem specifies the desired output format: CPU and memory percentage.
	// -p filters the output for the given PID.
	cmd := exec.CommandContext(ctx, "ps", "-o", "%cpu,rss", "-p", strconv.Itoa(h.pid))

	// Execute the command and capture its output.
	output, err := cmd.Output()
	if err != nil {
		// This error typically occurs if the PID does not exist.
		return nil, fmt.Errorf("failed to get process metrics for pid %d: %w", h.pid, err)
	}

	// The output from `ps` includes a header line, so we need to parse the second line.
	// Example output:
	// %CPU %MEM
	//  0.1  0.2
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected ps output for pid %d: %s", h.pid, output)
	}

	// The second line contains the data. We split it by whitespace.
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return nil, fmt.Errorf("unexpected ps output format for pid %d: %s", h.pid, lines[1])
	}

	// Parse the CPU usage from the first field.
	cpu, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPU usage '%s': %w", fields[0], err)
	}

	rssMemoryKB, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS memory '%s': %w", fields[1], err)
	}

	// Convert RSS memory from KB to bytes.
	memory := uint64(rssMemoryKB * 1024)

	return &types.Metrics{CPU: cpu, Mem: memory}, nil

}
