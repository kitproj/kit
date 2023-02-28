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
	_, err := w.writer.Write(output)
	// we must return the original length of p, not the length of the output, because we'll get broken pipes if we don't (SIGPIPE)
	return len(p), err
}
