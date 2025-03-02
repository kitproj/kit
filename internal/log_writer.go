package internal

import (
	"bytes"
	"log"
)

type logWriter struct {
	//
	fixes  func() (string, string)
	buffer bytes.Buffer
	logger *log.Logger
}

func (lw *logWriter) Write(p []byte) (int, error) {
	prefix, suffix := lw.fixes()

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
