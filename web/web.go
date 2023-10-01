//go:build !dev

package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS

func FS() fs.FS {
	fs, _ := fs.Sub(dist, "dist")
	return fs
}
