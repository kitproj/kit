package types

import (
	"fmt"
	"os"
)

func Environ(spec PodSpec, task Task) ([]string, error) {
	podEnviron, err := spec.Environ()
	if err != nil {
		return nil, fmt.Errorf("error getting spec environ: %w", err)
	}
	taskEnviron, err := task.Environ()
	if err != nil {
		return nil, fmt.Errorf("error getting spec environ: %w", err)
	}

	osEnviron := os.Environ()

	return append(append(podEnviron, taskEnviron...), osEnviron...), nil
}
