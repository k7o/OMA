package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var uiAssets embed.FS

// Assets returns the embedded filesystem for the UI.
func Assets() (fs.FS, error) {
	return fs.Sub(uiAssets, "dist")
}
