package contract

import (
	"context"
	"oma/models"
)

type App interface {
	Eval(ctx context.Context, req *models.EvalRequest) (*models.EvalResponse, error)
}
