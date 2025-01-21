package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("Spaces", func(t *testing.T) {
			s := Strings{}
			err := s.UnmarshalJSON([]byte(`"a b"`))
			assert.NoError(t, err)
			assert.Len(t, s, 2)
		})
		t.Run("SingleQuotes", func(t *testing.T) {
			s := Strings{}
			err := s.UnmarshalJSON([]byte(`"a 'b'"`))
			assert.NoError(t, err)
			assert.Len(t, s, 2)
		})
		t.Run("DoubleQuotes", func(t *testing.T) {
			s := Strings{}
			err := s.UnmarshalJSON([]byte(`"a \"b\""`))
			assert.NoError(t, err)
			assert.Len(t, s, 2)
		})
	})
}
