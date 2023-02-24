package main_test

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HTTP_cancelled_iterations(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cycle/2?timeout=3", nil)
		assert.NoError(t, err)

		_, err = http.DefaultClient.Do(req)

		assert.Error(t, err)

		assert.Eventually(t, func() bool {
			wholeErr, err := io.ReadAll(stderr)
			return assert.NoError(t, err) &&
				assert.Contains(t, string(wholeErr), "incoming request /cycle/2?timeout=3 was canceled by client")
		}, 2*time.Second, 150*time.Millisecond)

		cmd.Process.Kill()
		cmd.Wait()
	})
}
