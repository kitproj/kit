package internal

import (
	"sync"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	Name string     `json:"name"`
	Task types.Task `json:"task"`
	// logFile is the log file path
	logFile string
	// Phase represents the current execution state of the task:
	// - "pending": Task is registered but waiting for dependencies to complete
	// - "waiting": Task's dependencies are satisfied and is ready to execute
	// - "starting": Task is initializing but not yet fully running
	// - "running": Task is actively executing
	// - "stalled": Task is running but not making progress
	// - "succeeded": Task completed successfully
	// - "failed": Task completed with errors
	// - "cancelled": Task was manually stopped
	// - "skipped": Task was intentionally not executed
	Phase string `json:"phase"`
	// the message for the task phase, e.g. "exit code 1'
	Message string `json:"message,omitempty"`
	// cancel function
	cancel func()
	// a mutex
	mu *sync.Mutex
}

func (n TaskNode) blocked() bool {
	switch n.Phase {
	case "running", "stalled":
		return n.Task.GetType() == types.TaskTypeJob
	case "succeeded", "skipped":
		return false
	default:
		return true
	}
}
