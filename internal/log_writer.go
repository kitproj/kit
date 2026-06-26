package internal

import (
	"bytes"
	"log"
	"sync"
)

type logWriter struct {
	// prefixSuffixProvider returns the prefix and suffix to use when logging.
	prefixSuffixProvider func() (string, string)
	buffer               bytes.Buffer
	logger               *log.Logger
	// mu guards buffer: the same writer is used for both stdout and stderr,
	// which os/exec copies on separate goroutines
	mu sync.Mutex
}

func (lw *logWriter) Write(p []byte) (int, error) {
	prefix, suffix := lw.prefixSuffixProvider()

	lw.mu.Lock()
	defer lw.mu.Unlock()

	for _, b := range p {
		if b == '\n' {
			lw.logger.Printf("%s%s%s\n", prefix, lw.buffer.String(), suffix)
			lw.buffer.Reset()
		} else {
			lw.buffer.WriteByte(b)
		}
	}

	return len(p), nil
}
