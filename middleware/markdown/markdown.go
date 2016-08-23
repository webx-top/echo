package markdown

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	md2html "github.com/russross/blackfriday"
	"github.com/webx-top/echo"
)

type (
	Options struct {
		Path         string `json:"path"` //UrlPath
		Ext          string `json:"ext"`
		Index        string `json:"index"`
		Root         string `json:"root"`
		Browse       bool   `json:"browse"`
		Preprocessor func(echo.Context, []byte) []byte
	}
)

func Markdown(options ...*Options) echo.MiddlewareFunc {
	// Default options
	opts := new(Options)
	if len(options) > 0 {
		opts = options[0]
	}
	if opts.Index == "" {
		opts.Index = "SUMMARY.md"
	}
	opts.Root, _ = filepath.Abs(opts.Root)

	if opts.Preprocessor == nil {
		opts.Preprocessor = func(c echo.Context, b []byte) []byte {
			return b
		}
	}

	length := len(opts.Path)

	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			file := c.Request().URL().Path()
			if len(file) < length || file[0:length] != opts.Path {
				return next.Handle(c)
			}
			if len(opts.Ext) > 0 && !strings.HasSuffix(file, opts.Ext) {
				return next.Handle(c)
			}
			file = filepath.Clean(file[length:])
			absFile := filepath.Join(opts.Root, file)
			if !strings.HasPrefix(absFile, opts.Root) {
				return next.Handle(c)
			}
			fp, err := os.Open(absFile)
			if err != nil {
				return err
			}
			fi, err := fp.Stat()
			if err != nil {
				fp.Close()
				return err
			}
			w := c.Response()
			if fi.IsDir() {
				// Index file
				indexFile := filepath.Join(absFile, opts.Index)
				fp.Close()
				fp, err = os.Open(indexFile)
				if err == nil {
					fi, err = fp.Stat()
					defer fp.Close()
				}
				if err != nil || fi.IsDir() {
					if opts.Browse {
						fs := http.Dir(opts.Root)
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
						w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
						if _, err = fmt.Fprintf(w, "<pre>\n"); err != nil {
							return err
						}
						for _, d := range dirs {
							name := d.Name()
							color := "#212121"
							if d.IsDir() {
								color = "#e91e63"
								name += "/"
							} else {
								if len(opts.Ext) > 0 && !strings.HasSuffix(name, opts.Ext) {
									continue
								}
							}
							if _, err = fmt.Fprintf(w, "<a href=\"%s\" style=\"color: %s;\">%s</a>\n", name, color, name); err != nil {
								return err
							}
						}
						_, err = fmt.Fprintf(w, "</pre>\n")
						return err
					}
					return echo.ErrNotFound
				}
			} else {
				defer fp.Close()
			}
			modtime := fi.ModTime()
			if t, err := time.Parse(http.TimeFormat, c.Request().Header().Get(echo.HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
				w.Header().Del(echo.HeaderContentType)
				w.Header().Del(echo.HeaderContentLength)
				return c.NoContent(http.StatusNotModified)
			}

			b, err := ioutil.ReadAll(fp)
			if err != nil {
				return err
			}
			b = md2html.MarkdownCommon(b)
			b = opts.Preprocessor(c, b)
			w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
			w.Header().Set(echo.HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(b)
			return err
		})
	}
}
