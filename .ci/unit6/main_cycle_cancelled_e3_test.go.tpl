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

func TestHTTP_Cycle_cancelled(t *testing.T) {
	t.Run("reduced timeout cancelled iteration", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)
			resp, err := http.Get("http://localhost:8080/cycle/?timeout=3")
			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, resp.StatusCode, http.StatusOK, "http status is unexpected")
				wholeBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(wholeBody), "iteration 0 is processed")
				assert.Contains(t, string(wholeBody), "iteration 1 is cancelled")
				assert.NotContains(t, string(wholeBody), "iteration 2")
				assert.NotContains(t, string(wholeBody), "iteration 3")
				assert.NotContains(t, string(wholeBody), "iteration 4")
			}
			cmd.Process.Kill()
			cmd.Wait()
		})
	})
}
