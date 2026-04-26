package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLifecycleHook_GetCommand(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var h *LifecycleHook
		assert.Nil(t, h.GetCommand())
	})
	t.Run("Empty", func(t *testing.T) {
		h := &LifecycleHook{}
		assert.Nil(t, h.GetCommand())
	})
	t.Run("Command", func(t *testing.T) {
		h := &LifecycleHook{Command: Strings{"echo", "hello"}}
		assert.Equal(t, Strings{"echo", "hello"}, h.GetCommand())
	})
	t.Run("Sh", func(t *testing.T) {
		h := &LifecycleHook{Sh: "echo hello"}
		assert.Equal(t, Strings{"sh", "-c", "echo hello"}, h.GetCommand())
	})
	t.Run("CommandPreferredOverSh", func(t *testing.T) {
		h := &LifecycleHook{Command: Strings{"echo", "hi"}, Sh: "echo hello"}
		assert.Equal(t, Strings{"echo", "hi"}, h.GetCommand())
	})
}

func TestLifecycle_GetOnSuccessHook(t *testing.T) {
	t.Run("NilLifecycle", func(t *testing.T) {
		var l *Lifecycle
		assert.Nil(t, l.GetOnSuccessHook())
	})
	t.Run("NoOnSuccess", func(t *testing.T) {
		l := &Lifecycle{}
		assert.Nil(t, l.GetOnSuccessHook())
	})
	t.Run("WithOnSuccess", func(t *testing.T) {
		hook := &LifecycleHook{Sh: "echo success"}
		l := &Lifecycle{OnSuccess: hook}
		assert.Equal(t, hook, l.GetOnSuccessHook())
	})
}

func TestLifecycle_GetOnFailureHook(t *testing.T) {
	t.Run("NilLifecycle", func(t *testing.T) {
		var l *Lifecycle
		assert.Nil(t, l.GetOnFailureHook())
	})
	t.Run("NoOnFailure", func(t *testing.T) {
		l := &Lifecycle{}
		assert.Nil(t, l.GetOnFailureHook())
	})
	t.Run("WithOnFailure", func(t *testing.T) {
		hook := &LifecycleHook{Sh: "echo failed"}
		l := &Lifecycle{OnFailure: hook}
		assert.Equal(t, hook, l.GetOnFailureHook())
	})
}

func TestTask_GetOnSuccessHook(t *testing.T) {
	t.Run("NoLifecycle", func(t *testing.T) {
		task := &Task{}
		assert.Nil(t, task.GetOnSuccessHook())
	})
	t.Run("WithOnSuccess", func(t *testing.T) {
		hook := &LifecycleHook{Sh: "echo success"}
		task := &Task{Lifecycle: &Lifecycle{OnSuccess: hook}}
		assert.Equal(t, hook, task.GetOnSuccessHook())
	})
}

func TestTask_GetOnFailureHook(t *testing.T) {
	t.Run("NoLifecycle", func(t *testing.T) {
		task := &Task{}
		assert.Nil(t, task.GetOnFailureHook())
	})
	t.Run("WithOnFailure", func(t *testing.T) {
		hook := &LifecycleHook{Sh: "echo failed"}
		task := &Task{Lifecycle: &Lifecycle{OnFailure: hook}}
		assert.Equal(t, hook, task.GetOnFailureHook())
	})
}
