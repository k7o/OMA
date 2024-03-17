package main

import (
	"context"
	"log"
	"oma/app"
	"oma/internal/config"
	"oma/internal/db"
	"oma/internal/decisionlogs"
	"oma/internal/opa"
	"oma/transport/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()
	conf := &config.Config{
		Transport: config.TransportConfig{
			HTTP: http.Config{
				Port: 8080,
			},
		},
	}

	db, err := db.InitInMemoryDatabase(ctx)
	if err != nil {
		log.Fatal("Error initializing database", err)
	}

	decisionLogRepository := decisionlogs.New(db)
	opa := opa.New()

	app := app.New(conf, decisionLogRepository, opa)
	server := http.New(&conf.Transport.HTTP, app)
	if err := server.Run(); err != nil {
		log.Fatal("Error running server", err)
	}
}
