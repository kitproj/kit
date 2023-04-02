package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevel_Less(t *testing.T) {
	assert.True(t, LogLevelDebug.Less(LogLevelInfo))
	assert.True(t, LogLevelInfo.Less(LogLevelWarn))
	assert.True(t, LogLevelWarn.Less(LogLevelError))
	assert.True(t, LogLevelError.Less(LogLevelOff))
}

func TestLogLevelOf(t *testing.T) {
	assert.Equal(t, LogLevelDebug, LogLevelOf("DEBUG"))
	assert.Equal(t, LogLevelInfo, LogLevelOf("INFO"))
	assert.Equal(t, LogLevelWarn, LogLevelOf("WARN"))
	assert.Equal(t, LogLevelError, LogLevelOf("ERROR"))
	assert.Equal(t, LogLevelInfo, LogLevelOf(""))
}
