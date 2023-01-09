package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings_MarshalJSON(t *testing.T) {
	x := Strings{"a", `"b"`}
	data, err := x.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"a \"\"\"b\"\"\""`, string(data))
}
