package opa

import (
	"io"
	"net/http"
	"os"
)

func Download(url string) (string, error) {
	if _, err := os.Stat("./opa"); err == nil {
		return "./opa", nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create("./opa")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	err = os.Chmod("./opa", 0755)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
