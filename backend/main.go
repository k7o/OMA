package main

import (
	"context"
	"oma/app"
	"oma/internal/config"
	"oma/internal/db"
	"oma/internal/decisionlogs"
	"oma/internal/opa"
	"oma/internal/playgroundlogs"
	"oma/transport/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()
	conf := &config.Config{
		LogLevel: zerolog.DebugLevel,
		Transport: config.TransportConfig{
			HTTP: http.Config{
				Port: 8080,
			},
		},
	}

	zerolog.SetGlobalLevel(conf.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db, err := db.InitInMemoryDatabase(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("initializing database")
	}

	decisionLogRepository := decisionlogs.New(db)
	playgroundLogRepository := playgroundlogs.New(db)
	opa := opa.New()

	app := app.New(conf, decisionLogRepository, playgroundLogRepository, opa)
	server := http.New(&conf.Transport.HTTP, app)
	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("running server")
	}
}
