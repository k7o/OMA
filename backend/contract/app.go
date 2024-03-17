package contract

import (
	"context"
	"oma/internal/playgroundlogs"
	"oma/models"
)

type App interface {
	Eval(ctx context.Context, req *models.EvalRequest) (*models.EvalResponse, error)
	PlaygroundLogs(ctx context.Context) ([]playgroundlogs.PlaygroundLog, error)
}
