package engine

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/webx-top/echo/logger"
)

type (
	// Engine defines an interface for HTTP server.
	Engine interface {
		SetHandler(Handler)
		SetLogger(logger.Logger)
		Start()
	}

	// Request defines an interface for HTTP request.
	Request interface {
		Scheme() string
		Host() string
		URI() string
		URL() URL
		Header() Header
		Proto() string
		// ProtoMajor() int
		// ProtoMinor() int
		RemoteAddress() string
		Method() string
		SetMethod(string)
		Body() io.ReadCloser
		FormValue(string) string
		Object() interface{}

		Form() UrlValuer
		PostForm() UrlValuer
		MultipartForm() *multipart.Form
		IsTLS() bool
		Cookie(string) string
		Referer() string
		UserAgent() string
		FormFile(string) (multipart.File, *multipart.FileHeader, error)
	}

	// Response defines an interface for HTTP response.
	Response interface {
		Header() Header
		WriteHeader(int)
		Write(b []byte) (int, error)
		Status() int
		Size() int64
		Committed() bool
		SetWriter(io.Writer)
		Writer() io.Writer
		Object() interface{}

		Hijack(func(net.Conn))
		Body() []byte
		Redirect(string, int)
		NotFound()
		SetCookie(*http.Cookie)
		ServeFile(string)
	}

	// Header defines an interface for HTTP header.
	Header interface {
		Add(string, string)
		Del(string)
		Get(string) string
		Set(string, string)
		Object() interface{}
	}

	// Wrap url.Values
	UrlValuer interface {
		Add(string, string)
		Del(string)
		Get(string) string
		Set(string, string)
		Encode() string
		All() map[string][]string
		Reset(url.Values)
	}

	// URL defines an interface for HTTP request url.
	URL interface {
		SetPath(string)
		Path() string
		QueryValue(string) string
		Query() url.Values
		RawQuery() string
		Object() interface{}
	}

	// Config defines engine configuration.
	Config struct {
		Address            string
		TLSCertfile        string
		TLSKeyfile         string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxConnsPerIP      int
		MaxRequestsPerConn int
		MaxRequestBodySize int
	}

	// Handler defines an interface to server HTTP requests via `ServeHTTP(Request, Response)`
	// function.
	Handler interface {
		ServeHTTP(Request, Response)
	}

	// HandlerFunc is an adapter to allow the use of `func(Request, Response)` as HTTP handlers.
	HandlerFunc func(Request, Response)
)

// ServeHTTP serves HTTP request.
func (h HandlerFunc) ServeHTTP(req Request, res Response) {
	h(req, res)
}
