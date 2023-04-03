package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoff(t *testing.T) {
	assert.Equal(t, backoff{time.Second}.next(), backoff{time.Second * 2})
	assert.Equal(t, backoff{1000 * time.Second}.next(), backoff{time.Second * 8})
}
