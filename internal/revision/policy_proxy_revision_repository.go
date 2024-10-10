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

type PolicyProxyRevisionRepositoryConfig struct {
	URL string `json:"url"`
}

func (c *PolicyProxyRevisionRepositoryConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

type PolicyProxyRevisionRepository struct {
	conf *PolicyProxyRevisionRepositoryConfig
}

func NewPolicyProxyRevisionRepository(conf *PolicyProxyRevisionRepositoryConfig) *PolicyProxyRevisionRepository {
	return &PolicyProxyRevisionRepository{
		conf: conf,
	}
}

type PolicyProxyPackage struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type PolicyProxyRevisionFiles []struct {
	FileName string `json:"file_name"`
}

func (r *PolicyProxyRevisionRepository) ListRevisions() ([]models.Revision, error) {

	log.Info().Msg("ListRevisions")

	url := r.conf.URL

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
		log.Error().Int("status_code", resp.StatusCode).Bytes("body", body).Msg("received bad statuscode from policy proxy")
	}

	var revisions []PolicyProxyPackage
	if err := json.NewDecoder(resp.Body).Decode(&revisions); err != nil {
		log.Info().Msg(err.Error())
		return nil, err
	}

	var result []models.Revision
	for _, revision := range revisions {
		result = append(result, models.Revision{
			PackageId:   revision.ID,
			Name:        revision.Name,
			Version:     revision.Version,
			PackageType: "policy-proxy",
			CreatedAt:   revision.CreatedAt,
		})
	}
	return result, nil
}

func (r *PolicyProxyRevisionRepository) ListRevisionFiles(packageId string) ([]string, error) {
	log.Info().Msg("ListRevisionFiles")

	url := fmt.Sprintf("%s/%s/package_files", r.conf.URL, packageId)

	log.Info().Msg(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// bodyBytes, err := io.ReadAll(resp.Body)

	// // Log the body
	// bodyString := string(bodyBytes)
	// log.Info().Msg(bodyString)

	var files PolicyProxyRevisionFiles
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		log.Info().Msg(err.Error())
		return nil, err
	}

	var result []string
	for _, file := range files {
		result = append(result, file.FileName)
	}

	return result, nil
}

func (r *PolicyProxyRevisionRepository) DownloadRevisionById(revisionId string) (*models.Bundle, error) {
	log.Info().Msg("DownloadRevisionById")

	url := fmt.Sprintf("%s/?package_version=%s", r.conf.URL, revisionId)

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

func (r *PolicyProxyRevisionRepository) DownloadRevision(revision *models.Revision) (*models.Bundle, error) {
	log.Info().Msg("DownloadRevision")

	url := fmt.Sprintf("%s/%s/%s/%s/%s", r.conf.URL, revision.PackageType, revision.Name, revision.Version, revision.FileName)

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

func (r *PolicyProxyRevisionRepository) DownloadRevisionForPackage(packageId string, filename string) (*models.Bundle, error) {
	log.Info().Msg("DownloadRevisionForPackage")

	resp, err := http.Get(fmt.Sprintf("%s/%s", r.conf.URL, packageId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packageInfo PolicyProxyPackage
	if err := json.NewDecoder(resp.Body).Decode(&packageInfo); err != nil {
		return nil, err
	}

	return r.DownloadRevision(&models.Revision{
		PackageId:   packageId,
		FileName:    filename,
		Name:        packageInfo.Name,
		PackageType: "policy-proxy",
		Version:     packageInfo.Version,
		CreatedAt:   packageInfo.CreatedAt,
	})
}
