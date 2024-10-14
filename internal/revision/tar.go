package revision

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"oma/models"
)

// UnGzTar reads a tar.gz file from an io.Reader and returns a Bundle
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
