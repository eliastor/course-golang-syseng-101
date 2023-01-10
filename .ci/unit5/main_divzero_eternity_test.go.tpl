package main_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEternity(t *testing.T) {
	instructions := "div 10 0"
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
		assert.NoError(t, err)
		allOut, err := io.ReadAll(stdout)
		assert.NoError(t, err)
		assert.Contains(t, string(allOut), "eternity")

		allErr, err := io.ReadAll(stderr)
		assert.NoError(t, err)
		assert.Empty(t, "", string(allErr))
	})
}
