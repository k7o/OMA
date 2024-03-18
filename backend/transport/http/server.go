package http

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"oma/contract"
	"oma/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	conf *Config
	app  contract.App
}

func New(conf *Config, app contract.App) *Server {
	return &Server{
		conf: conf,
		app:  app,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	router.Route("/api", func(r chi.Router) {
		r.Post("/eval", s.eval)
		r.Post("/format", s.format)
		r.Post("/lint", s.lint)
		r.Post("/decision-log/logs", s.pushDecisionLog)
		r.Get("/decision-log/list", s.listDecisionLogs)
		r.Get("/playground-logs", s.playgroundLogs)
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}

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

func (s *Server) playgroundLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := s.app.PlaygroundLogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

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

func jsonReqBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	t := new(T)

	if r.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil, err
		}
		defer gr.Close()

		// Decode the JSON from the decompressed body
		if err := json.NewDecoder(gr).Decode(t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil, err
		}

	} else {
		if err := json.NewDecoder(r.Body).Decode(t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil, err
		}

	}

	return t, nil
}
