package proc

import (
	"context"
	"io"

	"github.com/alexec/kit/internal/types"
)

type Interface interface {
	// Init performs are one-time initialization.
	Init(ctx context.Context) error
	// Build does any build steps needed.
	Build(ctx context.Context, stdout, stderr io.Writer) error
	// Run runs the process.
	Run(ctx context.Context, stdout, stderr io.Writer) error
}

func New(c types.Container, spec types.Spec) Interface {
	if _, ok := imageIsHostfile(c.Image); c.Image == "" || ok {
		return &host{Container: c, Spec: spec}
	} else {
		return &container{Container: c, Spec: spec}
	}
}
