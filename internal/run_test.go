package internal

import (
	"bytes"
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/kitproj/kit/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRunSubgraph(t *testing.T) {

	setup := func(t *testing.T) (context.Context, context.CancelFunc, *log.Logger, *bytes.Buffer) {
		ctx, cancel := context.WithCancel(context.Background())
		buffer := &bytes.Buffer{}
		logger := log.New(buffer, "", 0)
		t.Cleanup(func() {
			t.Log(buffer.String())
			buffer.Reset()
		})
		return ctx, cancel, logger, buffer

	}

	t.Run("Task not found", func(t *testing.T) {
		ctx, cancel, logger, _ := setup(t)
		defer cancel()
		err := RunSubgraph(
			ctx,
			cancel,
			logger,
			&types.Workflow{},
			[]string{"job"},
			nil,
		)
		assert.EqualError(t, err, "task \"job\" not found in workflow")
	})

	t.Run("Skipped task not found", func(t *testing.T) {
		ctx, cancel, logger, _ := setup(t)
		defer cancel()
		err := RunSubgraph(
			ctx,
			cancel,
			logger,
			&types.Workflow{},
			nil,
			[]string{"job"},
		)
		assert.EqualError(t, err, "skipped task \"job\" not found in workflow")
	})

	t.Run("Single successful job", func(t *testing.T) {
		ctx, cancel, logger, _ := setup(t)
		defer cancel()
		wf := &types.Workflow{
			Tasks: map[string]types.Task{
				"job": {Command: []string{"true"}},
			},
		}
		err := RunSubgraph(
			ctx,
			cancel,
			logger,
			wf,
			[]string{"job"},
			nil,
		)
		assert.NoError(t, err)
	})

	t.Run("Single failing job", func(t *testing.T) {
		ctx, cancel, logger, _ := setup(t)
		defer cancel()
		wf := &types.Workflow{
			Tasks: map[string]types.Task{
				"job": {Command: []string{"false"}},
			},
		}
		err := RunSubgraph(
			ctx,
			cancel,
			logger,
			wf,
			[]string{"job"},
			nil,
		)
		assert.EqualError(t, err, "failed tasks: [job: exit status 1]")
	})
	t.Run("Single running service", func(t *testing.T) {
		ctx, cancel, logger, buffer := setup(t)
		defer cancel()

		wf := &types.Workflow{
			Tasks: map[string]types.Task{
				"service": {Command: []string{"cat"}, Ports: []types.Port{{}}},
			},
		}
		// services block until they are ready, so we must run them in in a goroutine
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := RunSubgraph(
				ctx,
				cancel,
				logger,
				wf,
				[]string{"service"},
				nil,
			)
			assert.NoError(t, err)
		}()

		time.Sleep(time.Second)

		cancel()

		wg.Wait()

		assert.Contains(t, buffer.String(), "[service] (succeeded)")
	})
	t.Run("Single failing service", func(t *testing.T) {
		ctx, cancel, logger, buffer := setup(t)
		defer cancel()

		wf := &types.Workflow{
			Tasks: map[string]types.Task{
				"service": {Command: []string{"false"}, Ports: []types.Port{{}}},
			},
		}
		// services block until they are ready, so we must run  them in in a goroutine
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := RunSubgraph(
				ctx,
				cancel,
				logger,
				wf,
				[]string{"service"},
				nil,
			)
			assert.EqualError(t, err, "failed tasks: [service: exit status 1]")
		}()

		time.Sleep(time.Second)
		cancel()

		wg.Wait()

		assert.Contains(t, buffer.String(), "[service] (failed) exit status 1")
	})
}
