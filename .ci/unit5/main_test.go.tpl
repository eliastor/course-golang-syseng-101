package main_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var outBinPath = "./unit_test" // filepath.Join(tmpDir, "u5e1")

func setEnvs(cmd *exec.Cmd, tmpDir string) {
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"CGO_ENABLED=0",
		"GOCACHE=" + tmpDir,
		"GOMODCACHE=" + tmpDir,
	}
}

func externalProgramTest(t *testing.T, ctx context.Context, command string, args []string, stdin io.Reader, f func(t *testing.T, stdout, stderr io.Reader, err error)) {
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
	cmd.Stdin = stdin
	err := cmd.Run()

	f(t, outBuf, errBuf, err)
}

func onlyErrorTest(t *testing.T, outBinPath, instructions string, expectedOut string) {
	t.Helper()

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
		assert.Error(t, err)

		allOut, err := io.ReadAll(stdout)
		assert.NoError(t, err)
		assert.Empty(t, string(allOut))

		allErr, err := io.ReadAll(stderr)
		assert.NoError(t, err)
		assert.Contains(t, string(allErr), expectedOut)
	})
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

	m.Run()

}

func TestRun(t *testing.T) {
	t.Run("NormalInstructions", func(t *testing.T) {
		instructions := "mul 2 2\ndiv 20 2\npow 3 3\npow 2 3\npow 2 4\npow 3 2"

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)
			allOut, err := io.ReadAll(stdout)
			assert.NoError(t, err)
			assert.Contains(t, string(allOut), "4\n10\n27\n8\n16\n9\n")

			allErr, err := io.ReadAll(stderr)
			assert.NoError(t, err)
			assert.Empty(t, "", string(allErr))
		})
	})

	t.Run("Overflow mul", func(t *testing.T) {
		onlyErrorTest(t, outBinPath,
			"mul 10000000000 10000000000",
			"Computation error: error in expression 10000000000 * 10000000000: integer overflow",
		)
	})

	t.Run("Overflow double mul", func(t *testing.T) {
		onlyErrorTest(t, outBinPath,
			"mul 10000000000 10000000000\nmul 10000000000 10000000000",
			"Computation error: error in expression 10000000000 * 10000000000: integer overflow",
		)
	})

	t.Run("normal pow", func(t *testing.T) {
		instructions := "pow 2 3\npow 2 4\npow 3 2\npow 3 3"
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
			assert.NoError(t, err)

			allOut, err := io.ReadAll(stdout)
			assert.NoError(t, err)
			assert.Equal(t, "8\n16\n9\n27", strings.TrimSpace(string(allOut)))

			allErr, err := io.ReadAll(stderr)
			assert.NoError(t, err)
			assert.Empty(t, string(allErr))
		})
	})

	t.Run("overflow pow", func(t *testing.T) {
		instructions := "pow 999 1000"
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		externalProgramTest(t, ctx, outBinPath, nil, strings.NewReader(instructions), func(t *testing.T, stdout, stderr io.Reader, err error) {
			assert.Error(t, err)

			allOut, err := io.ReadAll(stdout)
			assert.NoError(t, err)
			assert.Empty(t, string(allOut))

			allErr, err := io.ReadAll(stderr)
			assert.NoError(t, err)
			assert.Contains(t, string(allErr), "Computation error")
			assert.Contains(t, string(allErr), "999")
			assert.Contains(t, string(allErr), "1000")
			assert.Contains(t, string(allErr), "integer overflow")
		})
	})

}
