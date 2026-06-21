package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

var terminalWriter io.Writer = os.Stdout

var isTerminalWriter = func(w io.Writer) bool {
	file, ok := w.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

func setTerminalTitle(title string) {
	if terminalWriter == nil || !isTerminalWriter(terminalWriter) {
		return
	}
	title = strings.NewReplacer("\a", "", "\x1b", "", "\r", " ", "\n", " ").Replace(title)
	_, _ = fmt.Fprintf(terminalWriter, "\033]0;%s\033\\", title)
}

func ringTerminalBell() {
	if terminalWriter == nil || !isTerminalWriter(terminalWriter) {
		return
	}
	_, _ = io.WriteString(terminalWriter, "\a")
}

func workflowTitle(name string, nodes map[string]*TaskNode) string {
	if len(nodes) == 0 {
		return fmt.Sprintf("kit %s", name)
	}

	complete := 0
	running := 0
	failures := []string{}
	for _, node := range nodes {
		switch node.Phase {
		case "failed":
			failures = append(failures, node.Name)
		case "running", "stalled":
			running++
			complete++
		case "succeeded", "skipped":
			complete++
		}
	}

	switch {
	case len(failures) == 1:
		return fmt.Sprintf("kit %s: failed (%s)", name, failures[0])
	case len(failures) > 1:
		return fmt.Sprintf("kit %s: failed (%d)", name, len(failures))
	case complete == len(nodes) && running > 0:
		return fmt.Sprintf("kit %s: ready", name)
	case complete == len(nodes):
		return fmt.Sprintf("kit %s: done", name)
	default:
		return fmt.Sprintf("kit %s: starting (%d/%d)", name, complete, len(nodes))
	}
}
