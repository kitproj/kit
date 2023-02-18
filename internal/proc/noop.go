package proc

import (
	"context"
	"io"
)

type noop struct{}

func (n noop) Run(ctx context.Context, stdout, stderr io.Writer) error {
	return nil
}
func (n noop) Reset(ctx context.Context) error {
	return nil
}
