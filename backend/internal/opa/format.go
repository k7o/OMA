package opa

import (
	"fmt"
)

func (opa *Opa) Format(policy string) (string, error) {
	policyFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
	defer cleanup()
	if err != nil {
		return "", err
	}

	output, err := cmd(fmt.Sprintf("fmt %s", policyFile))
	if err != nil {
		return "", err
	}

	return string(output), nil
}
