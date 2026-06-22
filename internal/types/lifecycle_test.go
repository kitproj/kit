package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLifecycle_GetOnSuccess(t *testing.T) {
	t.Run("NilLifecycle", func(t *testing.T) {
		var l *Lifecycle
		assert.Equal(t, "", l.GetOnSuccess())
	})
	t.Run("NoOnSuccess", func(t *testing.T) {
		l := &Lifecycle{}
		assert.Equal(t, "", l.GetOnSuccess())
	})
	t.Run("WithOnSuccess", func(t *testing.T) {
		l := &Lifecycle{OnSuccess: "notify"}
		assert.Equal(t, "notify", l.GetOnSuccess())
	})
}

func TestLifecycle_GetOnFailure(t *testing.T) {
	t.Run("NilLifecycle", func(t *testing.T) {
		var l *Lifecycle
		assert.Equal(t, "", l.GetOnFailure())
	})
	t.Run("NoOnFailure", func(t *testing.T) {
		l := &Lifecycle{}
		assert.Equal(t, "", l.GetOnFailure())
	})
	t.Run("WithOnFailure", func(t *testing.T) {
		l := &Lifecycle{OnFailure: "alert"}
		assert.Equal(t, "alert", l.GetOnFailure())
	})
}

func TestTask_GetOnSuccess(t *testing.T) {
	t.Run("NoLifecycle", func(t *testing.T) {
		task := &Task{}
		assert.Equal(t, "", task.GetOnSuccess())
	})
	t.Run("WithOnSuccess", func(t *testing.T) {
		task := &Task{Lifecycle: &Lifecycle{OnSuccess: "notify"}}
		assert.Equal(t, "notify", task.GetOnSuccess())
	})
}

func TestTask_GetOnFailure(t *testing.T) {
	t.Run("NoLifecycle", func(t *testing.T) {
		task := &Task{}
		assert.Equal(t, "", task.GetOnFailure())
	})
	t.Run("WithOnFailure", func(t *testing.T) {
		task := &Task{Lifecycle: &Lifecycle{OnFailure: "alert"}}
		assert.Equal(t, "alert", task.GetOnFailure())
	})
}
