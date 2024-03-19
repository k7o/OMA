package contract

import (
	"context"
	"oma/internal/decisionlogs"
	"oma/internal/playgroundlogs"
	"oma/models"
)

type App interface {
	Eval(ctx context.Context, req *models.EvalRequest) (*models.EvalResponse, error)
	Format(ctx context.Context, req *models.FormatRequest) (*models.FormatResponse, error)
	Lint(ctx context.Context, req *models.LintRequest) (*models.LintResponse, error)
	TestAll(ctx context.Context, req *models.EvalRequest) (*models.TestAllResponse, error)

	PlaygroundLogs(ctx context.Context) ([]playgroundlogs.PlaygroundLog, error)

	PushDecisionLogs(ctx context.Context, req *models.DecisionLogRequest) error
	ListDecisionLogs(ctx context.Context) ([]decisionlogs.DecisionLog, error)
}
