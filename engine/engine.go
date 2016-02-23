package engine

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"time"

	"github.com/webx-top/echo/logger"
)

type (
	HandlerFunc func(Request, Response)

	Engine interface {
		SetHandler(HandlerFunc)
		SetLogger(logger.Logger)
		Start()
	}

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

	Header interface {
		Add(string, string)
		Del(string)
		Get(string) string
		Set(string, string)
		Object() interface{}
	}

	UrlValuer interface {
		Add(string, string)
		Del(string)
		Get(string) string
		Set(string, string)
		Encode() string
		All() map[string][]string
	}

	URL interface {
		SetPath(string)
		Path() string
		QueryValue(string) string
		RawQuery() string
		Object() interface{}
	}

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
)
