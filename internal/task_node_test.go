package internal

import (
	"testing"

	"github.com/kitproj/kit/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_taskNode_blocked(t *testing.T) {
	service := types.Task{Ports: []types.Port{{}}}
	t.Run("service running", func(t *testing.T) {
		n := TaskNode{Phase: "running", Task: service}
		assert.False(t, n.blocked())
	})
	t.Run("service waiting", func(t *testing.T) {
		n := TaskNode{Phase: "waiting", Task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service starting", func(t *testing.T) {
		n := TaskNode{Phase: "starting", Task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service succeeded", func(t *testing.T) {
		n := TaskNode{Phase: "succeeded", Task: service}
		assert.False(t, n.blocked())
	})
	t.Run("service cancelled", func(t *testing.T) {
		n := TaskNode{Phase: "cancelled", Task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service failed", func(t *testing.T) {
		n := TaskNode{Phase: "failed", Task: service}
		assert.True(t, n.blocked())
	})
	task := types.Task{}
	t.Run("task running", func(t *testing.T) {
		n := TaskNode{Phase: "running", Task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task waiting", func(t *testing.T) {
		n := TaskNode{Phase: "waiting", Task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task starting", func(t *testing.T) {
		n := TaskNode{Phase: "starting", Task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task succeeded", func(t *testing.T) {
		n := TaskNode{Phase: "succeeded", Task: task}
		assert.False(t, n.blocked())
	})
	t.Run("task failed", func(t *testing.T) {
		n := TaskNode{Phase: "failed", Task: task}
		assert.True(t, n.blocked())
	})
}
