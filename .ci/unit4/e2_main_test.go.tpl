package main_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setEnvs(cmd *exec.Cmd, tmpDir string) {
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
		"GOMODCACHE=" + tmpDir,
	}
}

func TestMain(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gotest-*")
	if !assert.NoError(t, err) {
		assert.FailNow(t, "test error: cannot create temporary dir")
	}
	defer os.RemoveAll(tmpDir)

	t.Run("BUILD", func(t *testing.T) {
		codepath := "."
		cmd := exec.Command("go", "build", "-o", "e2", codepath)
		setEnvs(cmd, tmpDir)
		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf

		err = cmd.Run()
		if !assert.NoError(t, err, errBuf.String()) {
			assert.FailNow(t, "test error: \"go build\" failed", err)
			return
		}
	})

	t.Run("RUN", func(t *testing.T) {
		cmd := exec.Command("timeout", "20", "./e2")
		setEnvs(cmd, tmpDir)
		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf

		go func() {
			err = cmd.Run()
			if !assert.NoError(t, err, errBuf.String()) {
				assert.FailNow(t, "test error: \"go run\" failed", err)
				return
			}
		}()

		time.Sleep(1 * time.Second)

		defer cmd.Process.Kill()

		t.Run("ROOT_HANDLER", func(t *testing.T) {
			resp, err := http.Get("http://localhost:8080/")
			assert.NoError(t, err, "Error while executing http.Get to lcoalhost:8080")

			respBody, err := io.ReadAll(resp.Body)

			assert.NoError(t, err, "Error occurend while reading request body")

			assert.Equal(t, "GET URI: / handling with / handler", strings.TrimRight(string(respBody), "\n"))

		})

		t.Run("ROOT_HANDLER_WITH_REQUEST_URI", func(t *testing.T) {
			resp, err := http.Get("http://localhost:8080/a/b/c?foo=bar")
			assert.NoError(t, err, "Error while executing http.Get to lcoalhost:8080")

			respBody, err := io.ReadAll(resp.Body)

			assert.NoError(t, err, "Error occurend while reading request body")

			assert.Equal(t, "GET URI: /a/b/c?foo=bar handling with / handler", strings.TrimRight(string(respBody), "\n"))
		})

		t.Run("ECHO HANDLER", func(t *testing.T) {
			testString := "AAAAAAAAAAAAAAAAAAAA"
			resp, err := http.Post("http://localhost:8080/echo", "application/x-www-form-urlencoded", strings.NewReader(testString))
			assert.NoError(t, err, "Error while executing http.Get to lcoalhost:8080")

			respBody, err := io.ReadAll(resp.Body)

			assert.NoError(t, err, "Error occurend while reading request body")

			assert.Equal(t, testString, string(respBody))
		})

		t.Run("GZIP_HANDLER", func(t *testing.T) {
			testString := "AAAAAAAAAAAAAAAAAAAA="

			buf := bytes.NewBuffer(nil)
			gzipper := gzip.NewWriter(buf)
			gzipper.Name = "/dev/stdin"
			gzipper.Write([]byte(testString))
			gzipper.Flush()
			gzipper.Close()
			resp, err := http.Post("http://localhost:8080/ungzip", "application/x-www-form-urlencoded", buf)
			assert.NoError(t, err, "Error while executing http.Get to lcoalhost:8080")

			respBody, err := io.ReadAll(resp.Body)

			assert.NoError(t, err, "Error occurend while reading request body")
			assert.Equal(t, 200, resp.StatusCode, string(respBody))
			assert.Equal(t, testString, string(respBody))
		})

		t.Run("GZIP_HANDLER_MALFORMED_QUERY", func(t *testing.T) {
			testString := "AAAAAAAAAAAAAAAAAAAA="

			resp, err := http.Post("http://localhost:8080/ungzip", "application/x-www-form-urlencoded", strings.NewReader(testString))
			assert.NoError(t, err, "Error while executing http.Get to lcoalhost:8080")

			respBody, err := io.ReadAll(resp.Body)

			assert.NoError(t, err, "Error occurend while reading request body")
			assert.Equal(t, 400, resp.StatusCode, string(respBody))
		})

		assert.NoError(t, err, "Can't run code")
	})
}
