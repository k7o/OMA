package contract

import (
	"oma/internal/decisionlogs"
	"oma/internal/playgroundlogs"
)

type DecisionLogRepository interface {
	decisionlogs.Querier
	MigrationEmbed
}

type PlaygroundLogsRepository interface {
	playgroundlogs.Querier
	MigrationEmbed
}
