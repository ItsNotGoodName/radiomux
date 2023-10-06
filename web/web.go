//go:build !dev

package web

import (
	"embed"
	"net/http"
)

//go:embed dist
var dist embed.FS

func DistFS() http.FileSystem {
	return http.FS(dist)
}
