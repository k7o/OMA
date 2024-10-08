package contract

import (
	"fmt"
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

type RevisionRepositoryType string

const (
	GitlabPackages  RevisionRepositoryType = "gitlab_packages"
	GitlabContainer RevisionRepositoryType = "gitlab_container"
)

func (t *RevisionRepositoryType) Validate() error {
	switch *t {
	case GitlabPackages, GitlabContainer:
		return nil
	case "":
		return fmt.Errorf("REVISION_CONFIG_TYPE is required")
	default:
		return fmt.Errorf("invalid REVISION_CONFIG_TYPE: %s", *t)
	}
}
