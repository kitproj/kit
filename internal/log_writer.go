package internal

import (
	"bytes"
	"log"
)

type logWriter struct {
	prefix string
	suffix string
	buffer bytes.Buffer
	logger *log.Logger
}

func (lw *logWriter) Write(p []byte) (int, error) {
	// reset color and bold

	for _, b := range p {
		if b == '\n' {
			lw.logger.Printf("%s%s%s\n", lw.prefix, lw.buffer.String(), lw.suffix)
			lw.buffer.Reset()
		} else {
			lw.buffer.WriteByte(b)
		}
	}

	return len(p), nil
}
