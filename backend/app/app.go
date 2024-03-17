package app

import (
	"context"
	"encoding/json"
	"errors"
	"oma/contract"
	"oma/internal/config"
	"oma/internal/playgroundlogs"
	"oma/models"

	"github.com/rs/zerolog/log"
)

type App struct {
	conf                     *config.Config
	decisionLogRepository    contract.DecisionLogRepository
	playgroundLogsRepository contract.PlaygroundLogsRepository
	opa                      contract.Opa
}

func New(conf *config.Config,
	decisionLogRepository contract.DecisionLogRepository,
	playgroundLogsRepository contract.PlaygroundLogsRepository,
	opa contract.Opa) *App {
	return &App{
		conf:                     conf,
		decisionLogRepository:    decisionLogRepository,
		playgroundLogsRepository: playgroundLogsRepository,
		opa:                      opa,
	}
}

func (a *App) Eval(ctx context.Context, req *models.EvalRequest) (*models.EvalResponse, error) {
	result, err := a.opa.Eval(req.Policy, req.Input)
	if err != nil {
		return nil, err
	}

	resp := result.MakeEvalResponse(req.Policy)
	resultJson, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, errors.New("failed to marshal result to json")
	}

	coverageJson, err := json.Marshal(resp.Coverage)
	if err != nil {
		return nil, errors.New("failed to marshal coverage to json")
	}

	_, err = a.playgroundLogsRepository.CreatePlaygroundLog(ctx, playgroundlogs.CreatePlaygroundLogParams{
		ID:        resp.Id,
		Input:     req.Input,
		Policy:    req.Policy,
		Timestamp: resp.Timestamp,
		Result:    string(resultJson),
		Coverage:  string(coverageJson),
	})
	if err != nil {
		log.Debug().Err(err).Msg("saving decision log to database")
		return nil, errors.New("failed to save decision log to database")
	}

	return resp, nil
}

func (a *App) Format(ctx context.Context, req *models.FormatRequest) (*models.FormatResponse, error) {
	policy, err := a.opa.Format(req.Policy)
	if err != nil {
		log.Debug().Err(err).Msg("formatting policy")
		return nil, err
	}

	return &models.FormatResponse{Formatted: policy}, nil
}

func (a *App) PlaygroundLogs(ctx context.Context) ([]playgroundlogs.PlaygroundLog, error) {
	logs, err := a.playgroundLogsRepository.ListPlaygroundlogs(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get decision logs from database")
		return nil, errors.New("failed to get decision logs from database")
	}

	return logs, nil
}
