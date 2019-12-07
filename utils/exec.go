package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

func execCapture(app string, args ...string) (string, string, error) {
	cmd := exec.Command(app, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}
	if errStdout != nil || errStderr != nil {
		return "", "", fmt.Errorf("failed to capture stdout or stderr")
	}

	return stdoutBuf.String(), stderrBuf.String(), err
}

func execOutput(app string, args ...string) (string, string, error) {
	cmd := exec.Command(app, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

// Exec ...
func Exec(app string, printOutput bool, args ...string) (string, string, error) {
	if printOutput {
		return execCapture(app, args...)
	}
	return execOutput(app, args...)
}
