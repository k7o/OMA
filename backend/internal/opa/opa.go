package opa

import (
	"encoding/json"
	"fmt"
	"oma/models"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Opa struct{}

func New() *Opa {
	return &Opa{}
}

func (opa *Opa) Eval(policy string, input string) (*models.EvalResult, error) {
	// Write the module to a temporary file.
	policyFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
	defer cleanup()
	if err != nil {
		return nil, err
	}

	// Write the input to a temporary file.
	inputFile, cleanup, err := writeBytesToFile([]byte(input), "json")
	defer cleanup()
	if err != nil {
		return nil, err
	}

	output, err := cmd(fmt.Sprintf("eval -d %s -i %s --profile data --coverage", policyFile, inputFile))
	if err != nil {
		return nil, err
	}

	var result models.EvalResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

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

func cmd(command string) ([]byte, error) {
	splits := strings.Split(command, " ")
	if command == "" || len(splits) == 0 {
		log.Debug().Msg("empty opa command")
		return nil, fmt.Errorf("empty opa command")
	}

	cmd := exec.Command("/opt/homebrew/bin/opa", splits...)
	output, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok && len(output) == 0 {
		stderr := string(exitErr.Stderr)
		return nil, fmt.Errorf("%s command: %s", splits[0], stderr)
	} else if err != nil && len(output) == 0 {
		return nil, err
	}

	return output, nil
}

func writeBytesToFile(data []byte, ext string) (string, func(), error) {
	file, err := os.CreateTemp("", fmt.Sprintf("tempfile*.%s", ext))
	if err != nil {
		return "", nil, err
	}

	// Write the data to the file.
	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
		return "", nil, err
	}

	// Define the cleanup function.
	cleanup := func() {
		os.Remove(file.Name())
	}

	return file.Name(), cleanup, nil
}
