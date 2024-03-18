package models

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
