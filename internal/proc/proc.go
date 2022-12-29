package proc

import (
	"context"
	"io"
	"time"

	"github.com/alexec/kit/internal/types"
)

type Interface interface {
	Init(ctx context.Context) error
	Build(ctx context.Context, stdout, stderr io.Writer) error
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Stop(ctx context.Context, grace time.Duration) error
}

func New(c types.Container) Interface {
	if _, ok := imageIsHostfile(c.Image); c.Image == "" || ok {
		return &host{Container: c}
	} else {
		return &container{Container: c}
	}
}
