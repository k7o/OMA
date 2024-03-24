package models

type DownloadBundleRequest struct {
	Revision Revision `json:"revision"`
}

type ApplicationSettings struct {
	OpaServerUrl    string `json:"opa_server_url"`
	BundleServerUrl string `json:"bundle_server_url"`
}

type DownloadRevisionResponse struct {
	Files *Bundle `json:"files"`
}
