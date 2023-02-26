package util

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemas(t *testing.T) {
	semas := NewSemaphores(map[string]int{})
	sema := semas.Get("")
	_ = sema.Acquire(context.Background(), 1)
	ok := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = sema.Acquire(context.Background(), 1)
		ok = true
		wg.Done()
	}()
	assert.False(t, ok)
	sema.Release(1)
	wg.Wait()
}
