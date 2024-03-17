package contract

import (
	"context"
	"oma/internal/decisionlogs"
)

type DecisionLogRepository interface {
	GetDecisionLog(ctx context.Context, decisionID string) (decisionlogs.DecisionLog, error)
	CreateDecisionLog(ctx context.Context, arg decisionlogs.CreateDecisionLogParams) (decisionlogs.DecisionLog, error)
	ListDecisionLogs(ctx context.Context) ([]decisionlogs.DecisionLog, error)
}
