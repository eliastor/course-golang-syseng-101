package unit1

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExercise5_Syntax(t *testing.T) {
}

func runCommandAndTest(t *testing.T, testfunc func(t *testing.T, stdOut string, stdErr string), wdir string, name string, args ...string) {

	cmd := exec.Command(name, args...)
	tmpDir := t.TempDir()

	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GOPATH=" + os.Getenv("GOPATH"),
		"HOME=" + os.Getenv("HOME"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
		"GOMODCACHE=" + tmpDir,
		// "GOMODCACHE=" + os.Getenv("GOMODCACHE"),
		"GOMOD=" + wdir + "/go.mod",
		"GOSUMDB=sum.golang.org",
		"GOPROXY=https://proxy.golang.org,direct",
	}
	if wdir != "" {
		cmd.Dir = wdir
	}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if !assert.NoError(t, err, errBuf.String()) {
		assert.FailNow(t, "test error: \""+name+"\" failed", err)
		return
	}
	testfunc(t, outBuf.String(), errBuf.String())
	outBuf = nil
	errBuf = nil
}

func TestExercise5_Output(t *testing.T) {
	codepath := "../../unit1/exercises/e5/main.go"

	defer os.Remove(path.Dir(codepath) + "/go.mod")
	defer os.Remove(path.Dir(codepath) + "/go.sum")
	runCommandAndTest(t, func(t *testing.T, stdOut, stdErr string) {
		// t.Log("go mod init OUT:", stdOut)
		// t.Log("go mod init ERR:", stdErr)
	}, path.Dir(codepath), "go", "mod", "init", "exercise5")

	runCommandAndTest(t, func(t *testing.T, stdOut, stdErr string) {
		// t.Log("go mod init OUT:", stdOut)
		// t.Log("go mod init ERR:", stdErr)
	}, path.Dir(codepath), "go", "get", "golang.org/x/tour/pic")

	runCommandAndTest(t, func(t *testing.T, stdOut, stdErr string) {
		// assert.Empty(t, stdErr, "stderr of your code must be empty", stdErr)
		assert.Greater(t, len(stdOut), 6, "Output is too short. Is it output from golang.org/x/tour/pic ?")
		assert.Equal(t, "IMAGE:", stdOut[0:6], "Output is not image. Is it output from golang.org/x/tour/pic ?")
	}, path.Dir(codepath), "timeout", "15", "go", "run", "main.go")

}
