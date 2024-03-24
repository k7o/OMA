package models

import "time"

type Revision struct {
	PackageId   string    `json:"package_id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	PackageType string    `json:"package_type"`
	FileName    string    `json:"file_name"`
	CreatedAt   time.Time `json:"created_at"`
}
