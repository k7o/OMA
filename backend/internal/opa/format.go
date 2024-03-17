package opa

import (
	"fmt"
	"os/exec"
)

func (opa *Opa) Format(policy string) (string, error) {
	policyFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
	defer cleanup()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("/opt/homebrew/bin/opa", "fmt", policyFile)
	output, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr := string(exitErr.Stderr)
		return "", fmt.Errorf("opa fmt failed: %s", stderr)
	} else if err != nil && len(output) == 0 {
		return "", err
	}

	return string(output), nil
}
