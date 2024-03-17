package opa

import (
	"encoding/json"
	"fmt"
	"oma/models"
	"os"
	"os/exec"
	"strings"

	"log"
)

type Opa struct{}

func New() *Opa {
	return &Opa{}
}

func (opa *Opa) Eval(policy string, input string) (*models.EvalResult, error) {
	// Write the module to a temporary file.
	moduleFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
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

	cmd := exec.Command("/opt/homebrew/bin/opa", "eval", "-d", moduleFile, "-i", inputFile, "--profile", "data", "--coverage")
	output, err := cmd.Output()
	if err != nil && len(output) == 0 {
		return nil, err
	}

	var result models.EvalResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (opa *Opa) MakeEvalResponse(result *models.EvalResult, policy string) *models.EvalResponse {
	return &models.EvalResponse{
		Result: makeResult(result, policy),
		Errors: result.Errors,
		Coverage: models.CoverageResponse{
			Covered:      makeCoverage(result.Coverage.Files),
			CoveredLines: result.Coverage.CoveredLines,
			Coverage:     int(result.Coverage.Coverage),
		},
	}
}

func makeCoverage(files map[string]models.Coverage) []models.Covered {
	covered := []models.Covered{}
	for _, file := range files {
		for _, c := range file.Covered {
			covered = append(covered, models.Covered{
				Start: c.Start.Row,
				End:   c.End.Row,
			})
		}
	}
	return covered
}

func makeResult(result *models.EvalResult, policy string) interface{} {
	if len(result.Result) == 0 {
		return nil
	} else if result.Result[0].Expressions == nil {
		return nil
	} else if len(result.Result[0].Expressions) == 0 {
		return nil
	}

	lines := strings.Split(policy, "\n")
	packageNesting := []string{}
	if len(lines) > 0 {
		if strings.HasPrefix(lines[0], "package ") {
			packageNesting = strings.Split(strings.TrimPrefix(lines[0], "package "), ".")
		}
	}

	return getPackageResult(result.Result[0].Expressions[0].Value, packageNesting)
}

func getPackageResult(result interface{}, splits []string) interface{} {
	if len(splits) == 0 {
		return result
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		return getPackageResult(resultMap[splits[0]], splits[1:])
	}

	log.Fatalf("Expected map[string]interface{} but got %T", result)
	return nil
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
