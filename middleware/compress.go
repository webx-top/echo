package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type (
	gzipWriter struct {
		io.Writer
		engine.Response
	}
)

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get(echo.ContentType) == `` {
		w.Header().Set(echo.ContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w gzipWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Response.(http.Hijacker).Hijack()
}

func (w *gzipWriter) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

var writerPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		scheme := `gzip`
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Response().Header().Add(echo.Vary, echo.AcceptEncoding)
			if strings.Contains(c.Request().Header().Get(echo.AcceptEncoding), scheme) {
				rw := c.Response().Writer()
				w := writerPool.Get().(*gzip.Writer)
				w.Reset(rw)
				defer func() {
					if c.Response().Size() == 0 {
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						c.Response().SetWriter(rw)
						c.Response().Header().Del(echo.ContentEncoding)
						w.Reset(ioutil.Discard)
					}
					w.Close()
					writerPool.Put(w)
				}()
				gw := gzipWriter{Writer: w, Response: c.Response()}
				c.Response().Header().Set(echo.ContentEncoding, scheme)
				c.Response().SetWriter(gw)
			}
			return h.Handle(c)
		})
	}
}
