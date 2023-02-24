package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Process_Simple(t *testing.T) {
	t.Run("Processed with 10s timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		processedStr := process(ctx, 1)
		assert.Equal(t, "iteration 1 is processed", processedStr)
		assert.NoError(t, ctx.Err())
	})
}
