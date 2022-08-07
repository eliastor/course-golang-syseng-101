package main_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnvs(cmd *exec.Cmd, tmpDir string) {
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
		"GOMODCACHE=" + tmpDir,
	}
}

type Logrecord struct {
	IP         string `json:"ip"`
	Username   string `json:"user"`
	Timestamp  uint64 `json:"time"`
	HTTPMethod string `json:"method"`
	URIPath    string `json:"path"`
	Size       uint   `json:"size"`
	HTTPCode   uint   `json:"code"`
}

func TestMain(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gotest-*")
	if !assert.NoError(t, err) {
		assert.FailNow(t, "test error: cannot create temporary dir")
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Build", func(t *testing.T) {
		codepath := "."
		cmd := exec.Command("go", "build", "-o", "e1", codepath)
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

	t.Run("fake", func(t *testing.T) {
		cmd := exec.Command("timeout", "100", "./e1")
		setEnvs(cmd, tmpDir)
		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf

		err = cmd.Run()
		if !assert.NoError(t, err, errBuf.String()) {
			assert.FailNow(t, "test error: \"go run\" failed", err)
			return
		}
		assert.NoError(t, err, "Can't run code")
		assert.Empty(t, errBuf.String(), "stderr of your code must be empty", errBuf.String())
		assert.NotEmpty(t, outBuf.String(), "stdout of your code must not be empty")

		outputLines := strings.Split(outBuf.String(), "\n")
		outputLines = outputLines[:len(outputLines)-1]
		assert.Equal(t, 10000, len(outputLines))
		for _, line := range outputLines {
			l := Logrecord{}
			err := json.Unmarshal([]byte(line), &l)
			require.NoError(t, err, "Can't unmarshal json. Output must be lines with valid json "+line)
			require.NotZero(t, l.HTTPCode, line)
			require.NotZero(t, l.Timestamp, line)
			require.NotEmpty(t, l.HTTPMethod, line)
			require.NotEmpty(t, l.IP, line)
			require.NotEmpty(t, l.URIPath, line)
		}
	})

	t.Run("stdin", func(t *testing.T) {
		stdin := `86.132.122.254 leet_coder - [18/07/2022:06:20:40 +0000] "GET /articles HTTP/1.1" 200 14425
119.151.42.59 sarah_cooper - [18/07/2022:06:20:40 +0000] "GET /trending HTTP/1.1" 404 11211`

		expected := `{"ip":"86.132.122.254","user":"leet_coder","time":1658125240,"method":"GET","path":"/articles","size":14425,"code":200}
{"ip":"119.151.42.59","user":"sarah_cooper","time":1658125240,"method":"GET","path":"/trending","size":11211,"code":404}
`
		cmd := exec.Command("timeout", "15", "./e1", "-")
		setEnvs(cmd, tmpDir)
		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf
		cmd.Stdin = strings.NewReader(stdin)

		err = cmd.Run()
		if !assert.NoError(t, err, errBuf.String()) {
			assert.FailNow(t, "test error: \"go run\" failed", err)
			return
		}
		assert.NoError(t, err, "Can't run code")
		assert.Empty(t, errBuf.String(), "stderr of your code must be empty", errBuf.String())
		assert.NotEmpty(t, outBuf.String(), "stdout of your code must not be empty")

		outputLines := strings.Split(outBuf.String(), "\n")
		expectedLines := strings.Split(expected, "\n")
		assert.Subset(t, outputLines, expectedLines, "not all lines are printed")
		assert.Subset(t, expectedLines, outputLines, "odd lines are printed")

	})

	t.Run("file", func(t *testing.T) {
		cmd := exec.Command("timeout", "30", "./e1", "../e0/fake.log")
		setEnvs(cmd, tmpDir)
		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf

		err = cmd.Run()
		if !assert.NoError(t, err, errBuf.String()) {
			assert.FailNow(t, "test error: \"go run\" failed", err)
			return
		}
		assert.NoError(t, err, "Can't run code")
		assert.Empty(t, errBuf.String(), "stderr of your code must be empty", errBuf.String())
		assert.NotEmpty(t, outBuf.String(), "stdout of your code must not be empty")

		outputLines := strings.Split(outBuf.String(), "\n")
		outputLines = outputLines[:len(outputLines)-1]
		assert.Equal(t, 10000, len(outputLines))
		for _, line := range outputLines {
			l := Logrecord{}
			err := json.Unmarshal([]byte(line), &l)
			require.NoError(t, err, "Can't unmarshal json. Output must be lines with valid json "+line)
			require.NotZero(t, l.HTTPCode, line)
			require.NotZero(t, l.Timestamp, line)
			require.NotEmpty(t, l.HTTPMethod, line)
			require.NotEmpty(t, l.IP, line)
			require.NotEmpty(t, l.URIPath, line)
		}
	})
}
