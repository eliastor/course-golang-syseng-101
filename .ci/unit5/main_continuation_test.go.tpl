package main_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContinuation(t *testing.T) {
	t.Run("continuation of execution", func(t *testing.T) {
		instructions := "mul 2 3\ndiv 3 0\ndiv 81 9"
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)

			allOut, err := io.ReadAll(stdout)
			assert.NoError(t, err)
			assert.Contains(t, string(allOut), "6\neternity\n9")

			allErr, err := io.ReadAll(stderr)
			assert.NoError(t, err)
			assert.Empty(t, string(allErr))
		})
	})
}
