package proc

import (
	"context"
	"io"

	"github.com/kitproj/kit/internal/types"
)

type noop struct{}

func (n noop) Run(ctx context.Context, stdout, stderr io.Writer) error {
	return nil
}

func (n noop) GetMetrics(ctx context.Context) (*types.Metrics, error) {
	return &types.Metrics{}, nil
}
