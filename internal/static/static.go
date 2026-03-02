package static

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distEmbed embed.FS

// GetStaticFS returns the static file system for the frontend
func GetStaticFS() fs.FS {
	f, err := fs.Sub(distEmbed, "dist")
	if err != nil {
		panic(err)
	}
	return f
}
