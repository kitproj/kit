package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashCompletion(t *testing.T) {
	output := bashCompletion("tasks.yaml")

	// Should contain the completion function
	assert.Contains(t, output, "_kit_completions()")
	assert.Contains(t, output, "complete -F _kit_completions kit")

	// Should use correct regex pattern for task names
	assert.Contains(t, output, `grep -E '^  [a-zA-Z0-9_-]+:\s*$'`)
}

func TestZshCompletion(t *testing.T) {
	output := zshCompletion("tasks.yaml")

	// Should contain the compdef header
	assert.Contains(t, output, "#compdef kit")

	// Should define the _kit function
	assert.Contains(t, output, "_kit()")

	// Should use correct regex pattern for task names
	assert.Contains(t, output, `grep -E '^  [a-zA-Z0-9_-]+:\s*$'`)

	// Should have guard against running during source
	assert.Contains(t, output, `if [ "$funcstack[1]" = "_kit" ]`)

	// Should conditionally register compdef
	assert.Contains(t, output, "compdef _kit kit")
}

func TestFishCompletion(t *testing.T) {
	output := fishCompletion("tasks.yaml")

	// Should define the tasks function
	assert.Contains(t, output, "__fish_kit_tasks")

	// Should use correct regex pattern for task names
	assert.Contains(t, output, `grep -E '^  [a-zA-Z0-9_-]+:\s*$'`)

	// Should have completions for flags
	assert.Contains(t, output, "complete -c kit")
}

func TestPrintCompletionInvalidShell(t *testing.T) {
	err := printCompletion("invalid", "tasks.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell: invalid")
}

func TestPrintCompletionValidShells(t *testing.T) {
	// These should not error (output goes to stdout)
	// We can't easily capture stdout, so just verify no error
	for _, shell := range []string{"bash", "zsh", "fish"} {
		// Redirect stdout temporarily
		old := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		err := printCompletion(shell, "tasks.yaml")

		w.Close()
		os.Stdout = old

		assert.NoError(t, err, "shell: %s", shell)
	}
}

func TestTaskNameRegexPattern(t *testing.T) {
	// Test the regex pattern we use in completion scripts
	// This simulates what the grep command does

	testYaml := `env:
  AWS_ENDPOINT_URL: http://localhost:4566

tasks:
  clean:
    sh: |
      echo "cleaning"
  build-app:
    command: go build .
    watch: src
  my_task_123:
    sh: echo "test"
`

	// Extract task names using the same pattern as our completion
	lines := strings.Split(testYaml, "\n")
	var tasks []string
	for _, line := range lines {
		// Match: starts with exactly 2 spaces, alphanumeric/hyphen/underscore, colon, end of line
		if len(line) >= 3 && line[0] == ' ' && line[1] == ' ' && line[2] != ' ' {
			// Check if line ends with : (possibly with trailing whitespace)
			trimmed := strings.TrimRight(line, " \t")
			if strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, ": ") {
				name := strings.TrimSpace(strings.TrimSuffix(trimmed, ":"))
				tasks = append(tasks, name)
			}
		}
	}

	assert.Equal(t, []string{"clean", "build-app", "my_task_123"}, tasks)
	// Should NOT include AWS_ENDPOINT_URL (has value after colon)
	// Should NOT include sh, command, watch (4-space indent)
}
