package main

import (
	"context"
	"io"
)

type ProcessDef interface {
	Init(ctx context.Context) error
	Build(ctx context.Context, stdout, stderr io.Writer) error
	Run(ctx context.Context, stdout, stderr io.Writer) error
	Stop(ctx context.Context) error
}
