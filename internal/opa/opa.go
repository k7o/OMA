package opa

import (
	"encoding/json"
	"errors"
	"fmt"
	"oma/models"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var checkFirstLineRegex = regexp.MustCompile(`(?s)(?m)(\d+\s+(?:error|errors) occurred)(.*?)(\d+:)`)
var policyFileRegex = regexp.MustCompile(`(?s)(?m).*(\d+): (.*)`)

var ErrExitError = errors.New("command exited")

type Opa struct{}

func New() *Opa {
	return &Opa{}
}

func (opa *Opa) Eval(bundle *models.Bundle, input string, options *models.EvalOptions) (*models.EvalResult, error) {
	if bundle == nil {
		return nil, errors.New("bundle is required")
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "temp-files-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Iterate over the map and write each file to the temporary directory
	dataFiles := []string{}
	for path, content := range *bundle {
		fullPath := filepath.Join(tempDir, path)
		// Ensure the directory structure exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return nil, err
		}
		// Write the file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return nil, err
		}

		dataFiles = append(dataFiles, fullPath)
	}

	// Write the input to a temporary file.
	inputFile, cleanup, err := writeBytesToFile([]byte(input), "json")
	defer cleanup()
	if err != nil {
		return nil, err
	}

	cmdString := fmt.Sprintf("eval -i %s --profile data", inputFile)
	if options.Coverage {
		cmdString += " --coverage"
	}

	for _, dataFile := range dataFiles {
		cmdString += fmt.Sprintf(" --data %s", dataFile)
	}

	output, err := cmd(cmdString)
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

func (opa *Opa) Lint(policy string) (string, []string, error) {
	policyFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
	defer cleanup()
	if err != nil {
		return "", nil, err
	}

	output, err := cmd(fmt.Sprintf("check %s", policyFile))
	if errors.Is(err, ErrExitError) {
		output = []byte(strings.TrimPrefix(err.Error(), ErrExitError.Error()+"\ncheck command: "))
	} else if err != nil {
		return "", nil, err
	}

	outputString := string(output)
	outputString = checkFirstLineRegex.ReplaceAllString(outputString, "$1: \n$3")

	var msg string
	var errors []string
	for i, line := range strings.Split(outputString, "\n") {
		line = strings.TrimSpace(line)
		if i == 0 {
			msg = line
			continue
		}

		if strings.HasPrefix(line, "/") {
			groups := policyFileRegex.FindStringSubmatch(line)
			log.Debug().Msgf("a: %s", groups)
			line = policyFileRegex.ReplaceAllString(line, "$1: $2")
			line = strings.TrimSpace(strings.TrimSuffix(line, ":"))
		}

		if line == "" {
			continue
		}

		errors = append(errors, fmt.Sprintf("line %s", line))
	}

	return msg, errors, nil
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
		return nil, errors.Join(ErrExitError, fmt.Errorf("%s command: %s", splits[0], stderr))
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
