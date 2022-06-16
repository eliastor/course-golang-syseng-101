package unit1

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExercise3_Syntax(t *testing.T) {
}

func TestExercise3_Output(t *testing.T) {
	codepath := "../../unit1/exercises/e3/main.go"

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
	if !assert.Equal(t, 3, len(outputLines), "there must be 2 lines with output and one blank line in the stdout") {
		t.FailNow()
	}
	assert.Regexp(t, "^[0-9]+(\\.[0-9]+)? <nil>$", outputLines[0], "first call must return square root and nil error")
	assert.Equal(t, "0 cannot Sqrt negative number: -2", outputLines[1], "second call must output error")
	assert.Equal(t, "", outputLines[2], "Println must be used, missing last newline")
}
