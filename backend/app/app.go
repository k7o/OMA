package app

import (
	"context"
	"oma/contract"
	"oma/internal/config"
	"oma/models"
)

type App struct {
	conf                  *config.Config
	decisionLogRepository contract.DecisionLogRepository
	opa                   contract.Opa
}

func New(conf *config.Config, decisionLogRepository contract.DecisionLogRepository,
	opa contract.Opa) *App {
	return &App{
		conf:                  conf,
		decisionLogRepository: decisionLogRepository,
		opa:                   opa,
	}
}

func (a *App) Eval(ctx context.Context, req *models.EvalRequest) (*models.EvalResponse, error) {
	result, err := a.opa.Eval(req.Policy, req.Input)
	if err != nil {
		return nil, err
	}

	return a.opa.MakeEvalResponse(result, req.Policy), nil
}
