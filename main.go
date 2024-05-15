package main

import (
	"context"
	"oma/app"
	"oma/internal/config"
	internalDb "oma/internal/db"
	"oma/internal/decisionlogs"
	"oma/internal/opa"
	"oma/internal/playgroundlogs"
	"oma/internal/revision"
	"oma/transport/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	ctx := context.Background()
	conf := &config.Config{}

	if err := envconfig.Init(&conf); err != nil {
		log.Fatal().Err(err).Msg("initializing environment variables")
	}

	zerolog.SetGlobalLevel(conf.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db, err := internalDb.InitInMemoryDatabase(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("initializing database")
	}

	decisionLogRepository := decisionlogs.New(db)
	playgroundLogRepository := playgroundlogs.New(db)
	revisionRepository := revision.NewGitlabRevisionRepository(&conf.RevisionConfig.Gitlab)
	opa := opa.New()

	err = internalDb.Migrate(ctx, db, decisionLogRepository, playgroundLogRepository)
	if err != nil {
		log.Fatal().Err(err).Msg("migrating database")
	}

	log.Info().Msgf("Started listening on port %d", conf.Transport.HTTP.Port)

	app := app.New(conf, decisionLogRepository, playgroundLogRepository, opa, revisionRepository)
	server := http.New(&conf.Transport.HTTP, app)
	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("running server")
	}
}
