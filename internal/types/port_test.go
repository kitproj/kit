package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPort_Unstring(t *testing.T) {

	t.Run("Invalid", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("foo")
		assert.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("8080:80")
		assert.NoError(t, err)
		assert.Equal(t, uint16(8080), p.ContainerPort)
		assert.Equal(t, uint16(80), p.HostPort)
	})

	t.Run("NoHostPort", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("8080")
		assert.NoError(t, err)
		assert.Equal(t, uint16(8080), p.ContainerPort)
		assert.Equal(t, uint16(8080), p.HostPort)
	})
}

func TestPort_String(t *testing.T) {
	t.Run("NoHostPort", func(t *testing.T) {
		p := &Port{ContainerPort: 8080}
		assert.Equal(t, "8080", p.String())
	})

	t.Run("WithHostPort", func(t *testing.T) {
		p := &Port{ContainerPort: 8080, HostPort: 80}
		assert.Equal(t, "8080:80", p.String())
	})
}
