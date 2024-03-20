package models

type DownloadBundleRequest struct {
	ApplicationSettings ApplicationSettings `json:"application_settings"`
}

type ApplicationSettings struct {
	OpaServerUrl    string `json:"opa_server_url"`
	BundleServerUrl string `json:"bundle_server_url"`
}

type DownloadBundleResponse struct {
	Files map[string]string `json:"files"`
}

type Bundle struct {
	Files map[string]struct {
		FileValue string `json:"file_value"`
		Files     Bundle `json:"files"`
	}
}
