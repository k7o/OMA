package http

import (
	"encoding/json"
	"net/http"
	"oma/models"

	"github.com/go-chi/chi"
)

func (s *Server) listRevisions(w http.ResponseWriter, r *http.Request) {
	result, err := s.app.ListRevisions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) revisionFiles(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "package_id")

	result, err := s.app.RevisionFiles(r.Context(), packageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) downloadPackage(w http.ResponseWriter, r *http.Request) {
	req := &models.DownloadBundleRequest{
		Revision: models.Revision{
			PackageId: chi.URLParam(r, "package_id"),
			FileName:  chi.URLParam(r, "file_name"),
		},
	}

	result, err := s.app.DownloadRevisionPackage(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	req := &models.DownloadBundleRequest{
		Revision: models.Revision{
			PackageType: chi.URLParam(r, "package_type"),
			Name:        chi.URLParam(r, "name"),
			Version:     chi.URLParam(r, "version"),
			FileName:    chi.URLParam(r, "file_name"),
		},
	}

	result, err := s.app.DownloadRevision(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
