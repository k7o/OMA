package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	router.Route("/api", func(r chi.Router) {
		r.Post("/eval", EvalHandler)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

type EvalRequest struct {
	Policy string `json:"policy"`
	Input  string `json:"input"`
	Data   string `json:"data"`
}

func EvalHandler(w http.ResponseWriter, r *http.Request) {
	var req EvalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := Eval(req.Policy, []byte(req.Input))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resultJSON))
}
