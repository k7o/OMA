package http

import (
	"encoding/json"
	"net/http"
	"oma/models"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func (s *Server) listRevisions(w http.ResponseWriter, r *http.Request) {
	result, err := s.app.ListRevisions(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("listing revisions")
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
		log.Error().Err(err).Msg("listing revision files")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) downloadRevisionById(w http.ResponseWriter, r *http.Request) {
	result, err := s.app.DownloadRevisionById(r.Context(), chi.URLParam(r, "revision_id"))
	if err != nil {
		log.Error().Err(err).Msg("downloading revision by id")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) downloadPackage(w http.ResponseWriter, r *http.Request) {
	req := &models.DownloadBundleRequest{
		Revision: models.Revision{
			PackageId:   chi.URLParam(r, "package_id"),
			FileName:    chi.URLParam(r, "file_name"),
			PackageType: r.URL.Query().Get("package_type"),
		},
	}

	result, err := s.app.DownloadRevisionPackage(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("downloading package")
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
		log.Error().Err(err).Msg("downloading revision")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
