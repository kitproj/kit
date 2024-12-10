package proc

import (
	"context"
	"io"

	"github.com/kitproj/kit/internal/types"
)

type Interface interface {
	// Run runs the process.
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Reset(ctx context.Context) error
}

func New(t types.Task, spec types.PodSpec) Interface {

	if t.Image == "" {
		if len(t.Command) == 0 {
			return &noop{}
		}
		return &host{Task: t, spec: spec}
	} else {
		return &container{Task: t, spec: spec}
	}
}
