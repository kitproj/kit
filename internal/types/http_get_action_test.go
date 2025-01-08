package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPGetAction_URL(t *testing.T) {
	a := HTTPGetAction{
		Scheme: "https",
		Port:   8080,
	}
	assert.Equal(t, "https://localhost:8080", a.URL().String())
}
