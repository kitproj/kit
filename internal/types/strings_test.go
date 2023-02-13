package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		x := Strings{"a", `"b"`}
		data, err := x.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `"a \"\"\"b\"\"\""`, string(data))
	})
	t.Run("UnmarshalJSON", func(t *testing.T) {
		s := &Strings{}
		err := s.UnmarshalJSON([]byte("a 'b'"))
		assert.NoError(t, err)
	})
}
