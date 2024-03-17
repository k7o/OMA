package opa

import (
	"log"
	"os"
	"os/exec"
	"time"
)

type OpaDecisionLogsPush []struct {
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

func StartOPAServer() error {
	cmd := exec.Command("opa", "run", "--server", "--config-file=./config.yaml", "--addr=localhost:8181", "--diagnostic-addr=localhost:8282")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
