package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPSocketAction_URL(t *testing.T) {
	x := &TCPSocketAction{Port: 8080}
	assert.Equal(t, "tcp://localhost:8080", x.URL().String())
}
