package executor

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", errors.New("empty command provided")
	}
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Special case: ls returns exit status 1 if no files match pattern
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 && strings.Contains(command, "ls") {
				return "No matching files found.", nil
			}
		}
		return stderr.String(), err
	}

	return out.String(), nil
}
