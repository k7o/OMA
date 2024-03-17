package models

type EvalRequest struct {
	Policy string `json:"policy"`
	Input  string `json:"input"`
	Data   string `json:"data"`
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
