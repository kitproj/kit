package types

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Envfile Strings

// Environ reads the returns the environ
func (f Envfile) Environ(workingDir string) ([]string, error) {
	var environ []string
	for _, e := range f {
		file, err := os.Open(filepath.Join(workingDir, e))
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") {
				environ = append(environ, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	return environ, nil
}
