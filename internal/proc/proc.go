package proc

import (
	"context"
	"io"
	"time"
)

type Proc interface {
	Init(ctx context.Context) error
	Build(ctx context.Context, stdout, stderr io.Writer) error
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Stop(ctx context.Context, grace time.Duration) error
}
