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
	// GetMetrics returns current resource metrics for this process
	GetMetrics(ctx context.Context) (*types.Metrics, error)
}

func New(name string, t types.Task, log *log.Logger, spec types.Spec) Interface {
	if t.Image != "" {
		return &container{
			name: name,
			log:  log,
			spec: spec,
			Task: t,
		}
	}
	if len(t.GetCommand()) > 0 {
		return &host{
			log:  log,
			spec: spec,
			Task: t,
		}
	}
	if len(t.Manifests) > 0 {
		return &k8s{
			name: name,
			log:  log,
			Task: t,
			spec: spec,
		}
	}
	return &noop{}
}
