package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isPortOpen(t *testing.T) {
	t.Run("RandomPort", func(t *testing.T) {
		err := isPortFree(0)
		assert.NoError(t, err)
	})
	t.Run("SudoPort", func(t *testing.T) {
		err := isPortFree(1)
		assert.Error(t, err)
	})
}
