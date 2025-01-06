package main

import (
	"testing"

	"github.com/kitproj/kit/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_taskNode_blocked(t *testing.T) {
	service := types.Task{Ports: []types.Port{{}}}
	t.Run("service running", func(t *testing.T) {
		n := taskNode{phase: "running", task: service}
		assert.False(t, n.blocked())
	})
	t.Run("service waiting", func(t *testing.T) {
		n := taskNode{phase: "waiting", task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service starting", func(t *testing.T) {
		n := taskNode{phase: "starting", task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service succeeded", func(t *testing.T) {
		n := taskNode{phase: "succeeded", task: service}
		assert.True(t, n.blocked())
	})
	t.Run("service failed", func(t *testing.T) {
		n := taskNode{phase: "failed", task: service}
		assert.True(t, n.blocked())
	})
	task := types.Task{}
	t.Run("task running", func(t *testing.T) {
		n := taskNode{phase: "running", task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task waiting", func(t *testing.T) {
		n := taskNode{phase: "waiting", task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task starting", func(t *testing.T) {
		n := taskNode{phase: "starting", task: task}
		assert.True(t, n.blocked())
	})
	t.Run("task succeeded", func(t *testing.T) {
		n := taskNode{phase: "succeeded", task: task}
		assert.False(t, n.blocked())
	})
	t.Run("task failed", func(t *testing.T) {
		n := taskNode{phase: "failed", task: task}
		assert.True(t, n.blocked())
	})
}
