package types

import (
	"fmt"
)

func Environ(spec Spec, task Task) ([]string, error) {
	specEnviron, err := spec.Environ()
	if err != nil {
		return nil, fmt.Errorf("error getting spec environ: %w", err)
	}
	taskEnviron, err := task.Environ()
	if err != nil {
		return nil, fmt.Errorf("error getting spec environ: %w", err)
	}

	return append(specEnviron, taskEnviron...), nil
}
