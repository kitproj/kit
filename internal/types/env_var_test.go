package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVar_String(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		s, err := EnvVar{Name: "FOO", Value: "1"}.String()
		assert.NoError(t, err)
		assert.Equal(t, "FOO=1", s)
	})
}
