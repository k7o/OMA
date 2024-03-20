package http

import (
	"encoding/json"
	"net/http"
	"oma/models"
)

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	req, err := jsonReqBody[models.DownloadBundleRequest](w, r)
	if err != nil {
		return
	}

	result, err := s.app.DownloadBundle(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
