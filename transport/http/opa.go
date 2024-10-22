package http

import (
	"encoding/json"
	"net/http"
	"oma/models"
)

func (s *Server) eval(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.EvalRequest](w, r)
	if err != nil {
		return
	}

	result, err := s.app.Eval(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) format(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.FormatRequest](w, r)
	if err != nil {
		return
	}

	result, err := s.app.Format(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) lint(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.LintRequest](w, r)
	if err != nil {
		return
	}

	result, err := s.app.Lint(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) testAll(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.EvalRequest](w, r)
	if err != nil {
		return
	}

	result, err := s.app.TestAll(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
