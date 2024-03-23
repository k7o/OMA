package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"log"
)

var (
	ErrorBundleNotInitialized = errors.New("bundle not initialized")
)

type Bundle map[string]string

func (b *Bundle) GetFile(name string) string {
	return (*b)[name]
}

func (b *Bundle) TarGz() (*bytes.Buffer, error) {
	if b == nil {
		return nil, ErrorBundleNotInitialized
	}

	// Create a buffer to hold the tar.gz data
	var buf bytes.Buffer

	// Create a new gzip writer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	// Create a new tar writer
	tw := tar.NewWriter(gz)
	defer tw.Close()

	// Iterate over the files and add them to the tar archive
	for name, content := range *b {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			log.Fatal(err)
		}
	}

	// Close the tar writer
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	// Close the gzip writer
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}

	return &buf, nil
}
