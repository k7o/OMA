package contract

import (
	"oma/internal/decisionlogs"
	"oma/internal/playgroundlogs"
	"oma/models"
)

type DecisionLogRepository interface {
	decisionlogs.Querier
	MigrationEmbed
}

type PlaygroundLogsRepository interface {
	playgroundlogs.Querier
	MigrationEmbed
}

type RevisionRepository interface {
	ListRevisions() ([]models.Revision, error)
	ListRevisionFiles(packageId string) ([]string, error)
	DownloadRevisionById(revisionId string) (*models.Bundle, error)
	DownloadRevision(revision *models.Revision) (*models.Bundle, error)
	DownloadRevisionForPackage(packageId string, filename string) (*models.Bundle, error)
}
