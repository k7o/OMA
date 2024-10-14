package revision

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oma/models"
	"time"

	"github.com/rs/zerolog/log"
)

type GitlabPackagesRevisionRepositoryConfig struct {
	URL          string `json:"url"`
	PrivateToken string `json:"private_token"`
}

func (c *GitlabPackagesRevisionRepositoryConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

type GitlabPackagesRevisionRepository struct {
	conf *GitlabPackagesRevisionRepositoryConfig
}

func NewGitlabPackagesRevisionRepository(conf *GitlabPackagesRevisionRepositoryConfig) *GitlabPackagesRevisionRepository {
	return &GitlabPackagesRevisionRepository{
		conf: conf,
	}
}

type GitlabPackage struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	PackageType string `json:"package_type"`
	Status      string `json:"status"`
	Links       struct {
		WebPath string `json:"web_path"`
	} `json:"_links"`
	CreatedAt        time.Time   `json:"created_at"`
	LastDownloadedAt interface{} `json:"last_downloaded_at"`
}

func (r *GitlabPackagesRevisionRepository) ListRevisions() ([]models.Revision, error) {
	url := r.conf.URL + "?sort=desc"
	if r.conf.PrivateToken != "" {
		url += fmt.Sprintf("&private_token=%s", r.conf.PrivateToken)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("error reading response body")
		}
		log.Error().Int("status_code", resp.StatusCode).Bytes("body", body).Msg("received bad statuscode from gitlab package repository")
	}

	var revisions []GitlabPackage
	if err := json.NewDecoder(resp.Body).Decode(&revisions); err != nil {
		return nil, err
	}

	var result []models.Revision
	for _, revision := range revisions {
		result = append(result, models.Revision{
			PackageId:   fmt.Sprintf("%d", revision.ID),
			Name:        revision.Name,
			Version:     revision.Version,
			PackageType: revision.PackageType,
			CreatedAt:   revision.CreatedAt,
		})
	}
	return result, nil
}

type RevisionFiles []struct {
	ID        int       `json:"id"`
	PackageID int       `json:"package_id"`
	CreatedAt time.Time `json:"created_at"`
	FileName  string    `json:"file_name"`
}

func (r *GitlabPackagesRevisionRepository) ListRevisionFiles(packageId string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/package_files", r.conf.URL, packageId)
	if r.conf.PrivateToken != "" {
		url += fmt.Sprintf("?private_token=%s", r.conf.PrivateToken)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var files RevisionFiles
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		result = append(result, file.FileName)
	}

	return result, nil
}

func (r *GitlabPackagesRevisionRepository) DownloadRevisionById(revisionId string) (*models.Bundle, error) {
	url := fmt.Sprintf("%s/?package_version=%s", r.conf.URL, revisionId)
	if r.conf.PrivateToken != "" {
		url += fmt.Sprintf("&private_token=%s", r.conf.PrivateToken)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packages []GitlabPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("revision not found")
	}

	files, err := r.ListRevisionFiles(fmt.Sprintf("%d", packages[0].ID))
	if err != nil {
		return nil, err
	}

	bundle := make(models.Bundle)
	for _, filename := range files {
		b, err := r.DownloadRevision(&models.Revision{
			PackageId:   fmt.Sprintf("%d", packages[0].ID),
			FileName:    filename,
			Name:        packages[0].Name,
			PackageType: packages[0].PackageType,
			Version:     packages[0].Version,
			CreatedAt:   packages[0].CreatedAt,
		})

		if err != nil {
			return nil, err
		}

		for k, v := range *b {
			bundle[k] = v
		}
	}

	return &bundle, nil
}

func (r *GitlabPackagesRevisionRepository) DownloadRevision(revision *models.Revision) (*models.Bundle, error) {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", r.conf.URL, revision.PackageType, revision.Name, revision.Version, revision.FileName)
	if r.conf.PrivateToken != "" {
		url += fmt.Sprintf("?private_token=%s", r.conf.PrivateToken)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	files, err := UnGzTar(resp.Body)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *GitlabPackagesRevisionRepository) DownloadRevisionForPackage(packageId string, filename string) (*models.Bundle, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", r.conf.URL, packageId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packageInfo GitlabPackage
	if err := json.NewDecoder(resp.Body).Decode(&packageInfo); err != nil {
		return nil, err
	}

	return r.DownloadRevision(&models.Revision{
		PackageId:   packageId,
		FileName:    filename,
		Name:        packageInfo.Name,
		PackageType: packageInfo.PackageType,
		Version:     packageInfo.Version,
		CreatedAt:   packageInfo.CreatedAt,
	})

}
