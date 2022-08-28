package sh

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Execute(prog string, args ...string) (string, error) {
	return execute("", prog, args...)
}

func execute(workingDir string, prog string, args ...string) (string, error) {
	cmd := exec.Command(prog, args...)
	cmd.Dir = workingDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute %q: %w", prog, err)
	}
	return out.String(), nil
}
