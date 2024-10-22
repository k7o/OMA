package http

import (
	"encoding/json"
	"net/http"
	"oma/models"

	"github.com/rs/zerolog/log"
)

func (s *Server) pushDecisionLog(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.DecisionLogRequest](w, r)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decode request body")
		return
	}

	if err := s.app.PushDecisionLogs(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) listDecisionLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := s.app.ListDecisionLogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
