package app

import (
	"context"
	"oma/models"
)

func (a *App) ListRevisions(ctx context.Context) ([]models.Revision, error) {
	revisions, err := a.revisionRepository.ListRevisions()
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func (a *App) RevisionFiles(ctx context.Context, packageId string) ([]string, error) {
	files, err := a.revisionRepository.ListRevisionFiles(packageId)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (a *App) DownloadRevisionById(ctx context.Context, revisionId string) (*models.DownloadRevisionResponse, error) {
	revision, err := a.revisionRepository.DownloadRevisionById(revisionId)
	if err != nil {
		return nil, err
	}

	return &models.DownloadRevisionResponse{
		Files: revision,
	}, nil
}

func (a *App) DownloadRevisionPackage(ctx context.Context, req *models.DownloadBundleRequest) (*models.DownloadRevisionResponse, error) {
	revision, err := a.revisionRepository.DownloadRevisionForPackage(req.Revision.PackageId, req.Revision.FileName)
	if err != nil {
		return nil, err
	}

	return &models.DownloadRevisionResponse{
		Files: revision,
	}, nil
}

func (a *App) DownloadRevision(ctx context.Context, req *models.DownloadBundleRequest) (*models.DownloadRevisionResponse, error) {
	revision, err := a.revisionRepository.DownloadRevision(&req.Revision)
	if err != nil {
		return nil, err
	}

	return &models.DownloadRevisionResponse{
		Files: revision,
	}, nil
}
