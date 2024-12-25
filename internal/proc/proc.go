package proc

import (
	"context"
	"io"
	"log"

	"github.com/kitproj/kit/internal/types"
)

type Interface interface {
	// Run runs the process.
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Reset(ctx context.Context) error
}

func New(t types.Task, log *log.Logger, spec types.PodSpec) Interface {
	if t.Image != "" {
		return &container{log: log, Task: t, spec: spec}
	}
	if len(t.Command) > 0 {
		return &host{log: log, Task: t, spec: spec}
	}
	if t.Sh != "" {
		return &shell{Task: t, spec: spec}
	}
	return &noop{}
}
