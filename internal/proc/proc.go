package proc

import (
	"context"
	"io"
	"time"

	"github.com/alexec/kit/internal/types"
)

type Proc interface {
	Init(ctx context.Context) error
	Build(ctx context.Context, stdout, stderr io.Writer) error
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Stop(ctx context.Context, grace time.Duration) error
}

func New(c types.Container) Proc {
	if c.Image == "" {
		return &host{Container: c}
	} else {
		return &container{Container: c}
	}
}
