package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/dgryski/trifles/uuid"
)

type EvalRequest struct {
	Options EvalOptions `json:"options"`
	Bundle  Bundle      `json:"bundle"`
	Input   string      `json:"input"`
	Data    string      `json:"data"`
}

type EvalOptions struct {
	Coverage bool   `json:"coverage"`
	Path     string `json:"path"`
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

type LintRequest struct {
	Policy string `json:"policy"`
}

type LintResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type TestAllResponse struct {
	Results []EvalResponse `json:"results"`
}

type DecisionLogRequest []DecisionLogRequestItem
type DecisionLogRequestItem struct {
	Labels struct {
		ID      string `json:"id"`
		Version string `json:"version"`
	} `json:"labels"`
	DecisionID string `json:"decision_id"`
	Bundles    map[string]struct {
		Revision string `json:"revision"`
	} `json:"bundles"`
	Path        string      `json:"path"`
	Result      interface{} `json:"result"`
	Input       interface{} `json:"input"`
	RequestedBy string      `json:"requested_by"`
	Timestamp   time.Time   `json:"timestamp"`
	Metrics     struct {
		CounterServerQueryCacheHit int `json:"counter_server_query_cache_hit"`
		TimerRegoExternalResolveNs int `json:"timer_rego_external_resolve_ns"`
		TimerRegoInputParseNs      int `json:"timer_rego_input_parse_ns"`
		TimerRegoQueryEvalNs       int `json:"timer_rego_query_eval_ns"`
		TimerServerHandlerNs       int `json:"timer_server_handler_ns"`
	} `json:"metrics"`
	ReqID int `json:"req_id"`
}

func (result *EvalResult) MakeEvalResponse(bundle *Bundle) *EvalResponse {
	return &EvalResponse{
		Id:     uuid.UUIDv4(),
		Result: parseResult(result),
		Errors: parseErrors(result.Errors),
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

func parseResult(result *EvalResult) interface{} {
	if len(result.Result) == 0 {
		return nil
	} else if result.Result[0].Expressions == nil {
		return nil
	} else if len(result.Result[0].Expressions) == 0 {
		return nil
	}

	return result.Result[0].Expressions[0].Value
}

var tempFileRegex = regexp.MustCompile("temp-files-[^/]*")

func parseErrors(errors []EvalError) []EvalError {
	if len(errors) == 0 {
		return errors
	}

	tempdir := ""
	indices := tempFileRegex.FindStringIndex(errors[0].Location.File)
	// If the pattern is found, trim the string up to the start index of the pattern
	if indices != nil {
		tempdir = errors[0].Location.File[:indices[1]]
	}

	for i := range errors {
		errors[i].Location.File = strings.TrimPrefix(errors[i].Location.File, tempdir)
	}

	return errors
}
