package app

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"oma/models"
)

func (a *App) DownloadBundle(ctx context.Context, req *models.DownloadBundleRequest) (*models.DownloadBundleResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/latest/bundle.tar.gz", req.ApplicationSettings.BundleServerUrl))
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	files := make(map[string]string)

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

	return &models.DownloadBundleResponse{
		Files: files,
	}, nil
}
