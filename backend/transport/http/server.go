package http

import (
	"encoding/json"
	"net/http"
	"oma/contract"
	"oma/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

func (s *Server) playgroundLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := s.app.PlaygroundLogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func jsonReqBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	t := new(T)

	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return t, nil
}
