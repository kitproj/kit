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
		lines, err := readEnvfile(filepath.Join(workingDir, e))
		if err != nil {
			return nil, err
		}
		environ = append(environ, lines...)
	}
	return environ, nil
}

func readEnvfile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var environ []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		environ = append(environ, line)
	}
	return environ, scanner.Err()
}
