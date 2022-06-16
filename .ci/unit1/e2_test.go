package unit1

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExercise2_Syntax(t *testing.T) {
}

func TestExercise2_Output(t *testing.T) {
	expected := `googleDNS: 8.8.8.8
loopback: 127.0.0.1
`

	codepath := "../../unit1/exercises/e2/main.go"

	cmd := exec.Command("timeout", "10", "go", "run", codepath)
	tmpDir, err := os.MkdirTemp("", "gotest-*")
	if !assert.NoError(t, err) {
		assert.FailNow(t, "test error: cannot create temporary dir")
	}
	defer os.RemoveAll(tmpDir)

	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
	}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err = cmd.Run()
	if !assert.NoError(t, err, errBuf.String()) {
		assert.FailNow(t, "test error: \"go run\" failed", err)
		return
	}

	assert.NoError(t, err, "Can't compile and run code")
	assert.Empty(t, errBuf.String(), "stderr of your code must be empty", errBuf.String())
	assert.NotEmpty(t, outBuf.String(), "stdout of your code must not be empty")

	outputLines := strings.Split(outBuf.String(), "\n")
	expectedLines := strings.Split(expected, "\n")
	assert.Subset(t, outputLines, expectedLines, "not all lines are printed")
	assert.Subset(t, expectedLines, outputLines, "odd lines are printed")
}
