package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbe_String(t *testing.T) {
	p := Probe{
		TCPSocket: &TCPSocketAction{
			Port: 8080,
		},
		InitialDelaySeconds: 1,
	}

	assert.Equal(t, "tcp://localhost:8080?initialDelay=1s", p.String())
}
