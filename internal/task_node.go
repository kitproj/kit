package internal

import (
	"sync"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	Name string `json:"name"`
	task types.Task
	// logFile is the log file path
	logFile string
	// the phase of the task, e.g. "pending", "waiting", "running", "stalled", "succeeded", "failed", "cancelled", "skipped"
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
		return n.task.GetType() == types.TaskTypeJob
	case "succeeded", "skipped":
		return false
	default:
		return true
	}
}
