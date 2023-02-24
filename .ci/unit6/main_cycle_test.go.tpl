package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Cycle(t *testing.T) {
	t.Run("All 3 processed with 10s timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		lines := cycle(ctx, 3)
		assert.Equal(t, []string{
			"iteration 0 is processed",
			"iteration 1 is processed",
			"iteration 2 is processed",
		}, lines)
		assert.NoError(t, ctx.Err())
	})
}
