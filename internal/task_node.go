package internal

import (
	"sync"
	"time"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	Name string `json:"name"`
	task types.Task
	// logFile is the log file path
	logFile string
	// the phase of the task, e.g. "pending", "waiting", "running", "succeeded", "failed", "cancelled", "skipped"
	Phase string `json:"phase"`
	// the message for the task phase, e.g. "exit code 1'
	Message string `json:"message,omitempty"`
	// Started at is the start time of the task
	StartedAt time.Time `json:"startedAt"`
	// Completed at is the end time of the task
	UpdatedAt time.Time `json:"updatedAt"`
	// cancel function
	cancel func()
	// a mutex
	mu *sync.Mutex
}

func (n TaskNode) blocked() bool {
	switch n.Phase {
	case "running":
		return n.task.GetType() == types.TaskTypeJob
	case "succeeded", "skipped":
		return false
	default:
		return true
	}
}
