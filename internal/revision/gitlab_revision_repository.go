package revision

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oma/models"
	"time"
)

type GitlabRevisionRepositoryConfig struct {
	GitlabPackagesURL string `json:"gitlab_packages_url"`
}

type GitlabRevisionRepository struct {
	conf *GitlabRevisionRepositoryConfig
}

func NewGitlabRevisionRepository(conf *GitlabRevisionRepositoryConfig) *GitlabRevisionRepository {
	return &GitlabRevisionRepository{
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

func (r *GitlabRevisionRepository) ListRevisions() ([]models.Revision, error) {
	resp, err := http.Get(r.conf.GitlabPackagesURL + "?sort=desc")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func (r *GitlabRevisionRepository) ListRevisionFiles(packageId string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/package_files", r.conf.GitlabPackagesURL, packageId))
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

func (r *GitlabRevisionRepository) DownloadRevision(revision *models.Revision) (*models.Bundle, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s/%s/%s", r.conf.GitlabPackagesURL, revision.PackageType, revision.Name, revision.Version, revision.FileName))
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

func (r *GitlabRevisionRepository) DownloadRevisionForPackage(packageId string, filename string) (*models.Bundle, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", r.conf.GitlabPackagesURL, packageId))
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

func UnGzTar(body io.Reader) (*models.Bundle, error) {
	gr, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	files := make(models.Bundle)

	// Extract the .tar file
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}
		content, err := io.ReadAll(tr)
		if err != nil {
			panic(err)
		}

		files[header.Name] = string(content)
	}

	return &files, nil
}
