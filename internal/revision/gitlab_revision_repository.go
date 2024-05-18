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
	PackagesURL  string `json:"packages_url"`
	PrivateToken string `json:"private_token"`
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
	url := r.conf.PackagesURL + "?sort=desc"
	if r.conf.PrivateToken != "" {
		url += fmt.Sprintf("&private_token=%s", r.conf.PrivateToken)
	}

	resp, err := http.Get(url)
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
	url := fmt.Sprintf("%s/%s/package_files", r.conf.PackagesURL, packageId)
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

func (r *GitlabRevisionRepository) DownloadRevisionById(revisionId string) (*models.Bundle, error) {
	url := fmt.Sprintf("%s/?package_version=%s", r.conf.PackagesURL, revisionId)
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

func (r *GitlabRevisionRepository) DownloadRevision(revision *models.Revision) (*models.Bundle, error) {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", r.conf.PackagesURL, revision.PackageType, revision.Name, revision.Version, revision.FileName)
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

func (r *GitlabRevisionRepository) DownloadRevisionForPackage(packageId string, filename string) (*models.Bundle, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", r.conf.PackagesURL, packageId))
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
