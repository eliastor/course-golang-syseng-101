package main_test

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTP_Timeout(t *testing.T) {
	t.Run("default timeout 6 cycles fail", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)

			_, err = http.Get("http://localhost:8080/cycle/6")

			assert.NoError(t, err)

			assert.Eventually(t, func() bool {
				wholeErr, err := io.ReadAll(stderr)
				return err == nil && strings.Contains(string(wholeErr), "incoming request /cycle/6 reached timeout")
			}, 2*time.Second, 150*time.Millisecond, "stderr of the server doesn't contain \"incoming request /cycle/6 reached timeout\"")

			cmd.Process.Kill()
			cmd.Wait()
		})
	})

	t.Run("extended timeout 2 cycles successful", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)

			resp, err := http.Get("http://localhost:8080/cycle/2?timeout=5")

			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, resp.StatusCode, http.StatusOK, "http status is unexpected")
				wholeBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(wholeBody), "iteration 0 is processed")
				assert.Contains(t, string(wholeBody), "iteration 1 is processed")
			}
			cmd.Process.Kill()
			cmd.Wait()
		})
	})

}
