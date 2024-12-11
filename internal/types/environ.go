package types

import (
	"fmt"
	"os"
	"strings"
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

	// create a map of environ
	environMap := make(map[string]string)
	for _, e := range append(podEnviron, taskEnviron...) {
		parts := strings.SplitN(e, "=", 2)
		environMap[parts[0]] = parts[1]
	}

	// add values for os.Environ() but only if already set
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if _, ok := environMap[parts[0]]; ok {
			environMap[parts[0]] = parts[1]
		}
	}

	// convert map back to slice
	environ := make([]string, 0, len(environMap))
	for k, v := range environMap {
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}

	return environ, nil
}
