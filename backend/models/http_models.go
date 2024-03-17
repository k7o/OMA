package models

import (
	"log"
	"strings"
	"time"

	"github.com/dgryski/trifles/uuid"
)

type EvalRequest struct {
	Policy string `json:"policy"`
	Input  string `json:"input"`
	Data   string `json:"data"`
}

type EvalResponse struct {
	Id        string           `json:"id"`
	Result    interface{}      `json:"result"`
	Errors    []EvalError      `json:"errors"`
	Coverage  CoverageResponse `json:"coverage"`
	Timestamp time.Time        `json:"timestamp"`
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

type FormatRequest struct {
	Policy string `json:"policy"`
}

type FormatResponse struct {
	Formatted string `json:"formatted"`
}

func (result *EvalResult) MakeEvalResponse(policy string) *EvalResponse {
	return &EvalResponse{
		Id:     uuid.UUIDv4(),
		Result: makeResult(result, policy),
		Errors: result.Errors,
		Coverage: CoverageResponse{
			Covered:      makeCoverage(result.Coverage.Files),
			CoveredLines: result.Coverage.CoveredLines,
			Coverage:     int(result.Coverage.Coverage),
		},
		Timestamp: time.Now(),
	}
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

func makeResult(result *EvalResult, policy string) interface{} {
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
