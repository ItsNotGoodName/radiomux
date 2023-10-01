// echoext contains extensions for Echo.
package echoext

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

// MountFS adds GET handlers to all files and directories for the given filesystem.
func MountFS(e *echo.Echo, f fs.FS) error {
	httpFS := http.FS(f)
	fsHandler := func(c echo.Context) error {
		http.StripPrefix("/", http.FileServer(httpFS)).ServeHTTP(c.Response(), c.Request())
		return nil
	}

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				e.GET("/"+name+"/*", fsHandler)
			} else if name == "index.html" {
				indexHandler := mountIndexGet(httpFS)
				e.GET("/", indexHandler)
				e.GET("/index.html", indexHandler)
			} else {
				e.GET("/"+name, fsHandler)
			}
		}
	} else if err != fs.ErrNotExist {
		return err
	}

	return nil
}

// mountIndexGet returns an index.html handler for the given filesystem.
// This is required because HTTP.FileServer redirects to the root of the directory if it sees an "index.html".
func mountIndexGet(httpFS http.FileSystem) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response()

		index, err := httpFS.Open("/index.html")
		if err != nil {
			msg, code := mountToHTTPError(err)
			http.Error(w, msg, code)
			return nil
		}

		stat, err := index.Stat()
		if err != nil {
			msg, code := mountToHTTPError(err)
			http.Error(w, msg, code)
			return nil
		}

		http.ServeContent(w, c.Request(), "index.html", stat.ModTime(), index)
		return nil
	}
}

// mountToHTTPError is copied from "net/http/fs.go".
func mountToHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
