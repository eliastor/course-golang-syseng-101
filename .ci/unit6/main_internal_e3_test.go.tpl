package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Process_Simple2(t *testing.T) {
	t.Run("Cancelled with 1s timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		processedStr := process(ctx, 1)
		assert.Equal(t, "iteration 1 is cancelled", processedStr)
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	})
}

func Test_Cycle_Cancelled9(t *testing.T) {
	t.Run("1 from 3 cancelled with 5s timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		lines := cycle(ctx, 3)
		assert.Equal(t, []string{
			"iteration 0 is processed",
			"iteration 1 is processed",
			"iteration 2 is cancelled",
		}, lines)
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	})
}
