package main

import "time"

type backoff struct {
	time.Duration
}

func (t backoff) next() backoff {
	const max = time.Second * 16
	if t.Duration > max {
		return backoff{max}
	}
	return backoff{t.Duration * 2}
}

var defaultBackoff = backoff{time.Second}
