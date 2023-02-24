package main_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var outBinPath = "./unit_test" // filepath.Join(tmpDir, "u5e1")

var exerciseID = flag.String("e", "1", "exercise number")

func setEnvs(cmd *exec.Cmd, tmpDir string) {
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
		"GOMODCACHE=" + tmpDir,
	}
}

func externalProgramTest(t *testing.T, ctx context.Context, command string, args []string, f func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error)) {
	t.Helper()

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	err := cmd.Start()

	require.Eventually(t, func() bool {
		_, err := http.Get("http://localhost:8080/")
		return err == nil
	}, 2*time.Second, 150*time.Millisecond, "the server hasn't started in 2 seconds")

	f(t, cmd, outBuf, errBuf, err)
}

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "gotest-*")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer os.RemoveAll(tmpDir)

	codepath := "."
	cmd := exec.Command("go", "build", "-o", outBinPath, codepath)
	setEnvs(cmd, tmpDir)
	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err = cmd.Run()
	if err != nil {
		fmt.Println("test error: \"go build\" failed", err)
		os.Exit(1)
	}

	defer os.Remove(outBinPath)

	m.Run()
}

func Test_HTTP(t *testing.T) {
	t.Run("1 successful iteration", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)
			resp, err := http.Get("http://localhost:8080/cycle/1")
			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, resp.StatusCode, http.StatusOK, "http status is unexpected")
				wholeBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(wholeBody), "iteration 0 is processed")
			}
			cmd.Process.Kill()
			cmd.Wait()
		})
	})

	t.Run("1 canceled iteration", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, func(t *testing.T, cmd *exec.Cmd, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cycle/1", nil)
			assert.NoError(t, err)

			_, err = http.DefaultClient.Do(req)

			assert.Error(t, err)

			assert.Eventually(t, func() bool {
				wholeErr, err := io.ReadAll(stderr)
				return err == nil && strings.Contains(string(wholeErr), "incoming request /cycle/1 was canceled by client")
			}, 2*time.Second, 150*time.Millisecond)

			cmd.Process.Kill()
			cmd.Wait()
		})
	})
}
