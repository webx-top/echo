package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
)

type (
	Static struct {
		Root      string
		Browse    bool
		Index     string
		isAbsRoot bool
	}

	FaviconOptions struct {
	}
)

func (s Static) Handle(c echo.Context) error {
	if !s.isAbsRoot {
		s.isAbsRoot = true
		s.Root, _ = filepath.Abs(s.Root)
	}
	file := c.P(0)
	absFile := filepath.Join(s.Root, file)
	if !strings.HasPrefix(absFile, s.Root) {
		return echo.ErrNotFound
	}
	fi, err := os.Stat(absFile)
	if err != nil {
		return err
	}
	w := c.Response()
	if fi.IsDir() {
		/* NOTE:
		Not checking the Last-Modified header as it caches the response `304` when
		changing differnt directories for the same path.
		*/

		// Index file
		indexFile := filepath.Join(file, s.Index)
		fi, err = os.Stat(indexFile)
		if err != nil || fi.IsDir() {
			if s.Browse {
				fs := http.Dir(s.Root)
				d, err := fs.Open(file)
				if err != nil {
					return echo.ErrNotFound
				}
				defer d.Close()
				dirs, err := d.Readdir(-1)
				if err != nil {
					return err
				}

				// Create a directory index
				w.Header().Set(echo.ContentType, echo.TextHTMLCharsetUTF8)
				if _, err = fmt.Fprintf(w, "<pre>\n"); err != nil {
					return err
				}
				for _, d := range dirs {
					name := d.Name()
					color := "#212121"
					if d.IsDir() {
						color = "#e91e63"
						name += "/"
					}
					if _, err = fmt.Fprintf(w, "<a href=\"%s\" style=\"color: %s;\">%s</a>\n", name, color, name); err != nil {
						return err
					}
				}
				_, err = fmt.Fprintf(w, "</pre>\n")
				return err
			}
			return echo.ErrNotFound
		} else {
			absFile = indexFile
		}
	}
	w.ServeFile(absFile)
	return nil
	// TODO:
	// http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
}

// Favicon serves the default favicon - GET /favicon.ico.
func Favicon(root string, options ...FaviconOptions) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			return nil
		})
	}
}
