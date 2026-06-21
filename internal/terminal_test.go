package internal

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetTerminalTitle(t *testing.T) {
	var buffer bytes.Buffer

	previousWriter := terminalWriter
	previousIsTerminal := isTerminalWriter
	terminalWriter = &buffer
	isTerminalWriter = func(io.Writer) bool { return true }
	t.Cleanup(func() {
		terminalWriter = previousWriter
		isTerminalWriter = previousIsTerminal
	})

	setTerminalTitle("hello\nworld\a")
	assert.Equal(t, "\033]0;hello world\033\\", buffer.String())
}

func TestWorkflowTitle(t *testing.T) {
	t.Run("ready", func(t *testing.T) {
		title := workflowTitle("kit", map[string]*TaskNode{
			"api": {Name: "api", Phase: "running"},
		})
		assert.Equal(t, "kit kit: ready", title)
	})

	t.Run("done", func(t *testing.T) {
		title := workflowTitle("kit", map[string]*TaskNode{
			"job": {Name: "job", Phase: "succeeded"},
		})
		assert.Equal(t, "kit kit: done", title)
	})

	t.Run("failed", func(t *testing.T) {
		title := workflowTitle("kit", map[string]*TaskNode{
			"job": {Name: "job", Phase: "failed"},
		})
		assert.Equal(t, "kit kit: failed (job)", title)
	})
}
