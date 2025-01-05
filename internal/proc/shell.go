package proc

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/kitproj/kit/internal/types"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type shell struct {
	spec types.PodSpec
	types.Task
}

func (s shell) Run(ctx context.Context, stdout, stderr io.Writer) error {

	environ, err := types.Environ(s.spec, s.Task)
	if err != nil {
		return err
	}

	r, err := interp.New(
		interp.Env(expand.ListEnviron(append(environ, os.Environ()...)...)),
		interp.Dir(s.WorkingDir),
		interp.StdIO(os.Stdin, stdout, stderr),
	)
	if err != nil {
		return err
	}
	p, err := syntax.NewParser().
		Parse(strings.NewReader(s.Sh), "")
	if err != nil {
		return err
	}

	err = r.Run(ctx, p)
	if err != nil {
		return err
	}

	return nil
}
