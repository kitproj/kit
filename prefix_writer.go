package main

import "io"

type prefixWriter struct {
	prefix string
	writer io.Writer
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	prefixBytes := []byte(w.prefix)
	output := make([]byte, len(prefixBytes)+len(p))
	copy(output, prefixBytes)
	copy(output[len(prefixBytes):], p)
	return w.writer.Write(output)
}
