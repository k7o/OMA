package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"log"
)

type EvalResult struct {
	Result []struct {
		Expressions []struct {
			Value    interface{} `json:"value"`
			Text     string      `json:"text"`
			Location struct {
				Row int `json:"row"`
				Col int `json:"col"`
			} `json:"location"`
		} `json:"expressions"`
	} `json:"result"`
	Errors  []EvalError `json:"errors"`
	Metrics struct {
		TimerRegoExternalResolveNs int `json:"timer_rego_external_resolve_ns"`
		TimerRegoLoadFilesNs       int `json:"timer_rego_load_files_ns"`
		TimerRegoModuleCompileNs   int `json:"timer_rego_module_compile_ns"`
		TimerRegoModuleParseNs     int `json:"timer_rego_module_parse_ns"`
		TimerRegoQueryCompileNs    int `json:"timer_rego_query_compile_ns"`
		TimerRegoQueryEvalNs       int `json:"timer_rego_query_eval_ns"`
		TimerRegoQueryParseNs      int `json:"timer_rego_query_parse_ns"`
	} `json:"metrics"`
	Profile []struct {
		TotalTimeNs int `json:"total_time_ns"`
		NumEval     int `json:"num_eval"`
		NumRedo     int `json:"num_redo"`
		NumGenExpr  int `json:"num_gen_expr"`
		Location    struct {
			File string `json:"file"`
			Row  int    `json:"row"`
			Col  int    `json:"col"`
		} `json:"location"`
	} `json:"profile"`
	Coverage struct {
		Files           map[string]Coverage `json:"files"`
		CoveredLines    int                 `json:"covered_lines"`
		NotCoveredLines int                 `json:"not_covered_lines"`
		Coverage        float64             `json:"coverage"`
	} `json:"coverage"`
}

type EvalError struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	Location struct {
		File string `json:"file"`
		Row  int    `json:"row"`
		Col  int    `json:"col"`
	} `json:"location"`
}

type Coverage struct {
	Covered []struct {
		Start struct {
			Row int `json:"row"`
		} `json:"start"`
		End struct {
			Row int `json:"row"`
		} `json:"end"`
	} `json:"covered"`
	CoveredLines int `json:"covered_lines"`
	Coverage     int `json:"coverage"`
}

func Eval(policy string, input []byte) (*EvalResponse, error) {
	// Write the module to a temporary file.
	moduleFile, cleanup, err := writeBytesToFile([]byte(policy), "rego")
	defer cleanup()
	if err != nil {
		return nil, err
	}

	// Write the input to a temporary file.
	inputFile, cleanup, err := writeBytesToFile(input, "json")
	defer cleanup()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("/opt/homebrew/bin/opa", "eval", "-d", moduleFile, "-i", inputFile, "--profile", "data", "--coverage")
	output, err := cmd.Output()
	if err != nil && len(output) == 0 {
		return nil, err
	}

	var result EvalResult
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	return &EvalResponse{
		Result: makeResult(result, policy),
		Errors: result.Errors,
		Coverage: CoverageResponse{
			Covered:      makeCoverage(result.Coverage.Files),
			CoveredLines: result.Coverage.CoveredLines,
			Coverage:     int(result.Coverage.Coverage),
		},
	}, nil
}

func makeCoverage(files map[string]Coverage) []Covered {
	covered := []Covered{}
	for _, file := range files {
		for _, c := range file.Covered {
			covered = append(covered, Covered{
				Start: c.Start.Row,
				End:   c.End.Row,
			})
		}
	}
	return covered
}

func makeResult(result EvalResult, policy string) interface{} {
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

type EvalResponse struct {
	Result   interface{}      `json:"result"`
	Errors   []EvalError      `json:"errors"`
	Coverage CoverageResponse `json:"coverage"`
}

type CoverageResponse struct {
	Covered      []Covered `json:"covered"`
	CoveredLines int       `json:"covered_lines"`
	Coverage     int       `json:"coverage"`
}

type Covered struct {
	Start int `json:"start"`
	End   int `json:"end"`
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
