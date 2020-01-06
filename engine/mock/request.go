package mock

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/webx-top/echo/engine"
)

var ErrNotImplemented = errors.New("not implemented")

func NewRequest() engine.Request {
	return &Request{
		header: &Header{},
		form:   &Value{},
		url:    &URL{},
	}
}

type Request struct {
	header *Header
	form   *Value
	url    *URL
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (r *Request) Scheme() string {
	return `https`
}

// Host returns HTTP request host. Per RFC 2616, this is either the value of
// the `Host` header or the host name given in the URL itself.
func (r *Request) Host() string {
	return `127.0.0.1`
}

// SetHost sets the host of the request.
func (r *Request) SetHost(string) {

}

// URI returns the unmodified `Request-URI` sent by the client.
func (r *Request) URI() string {
	return ``
}

// SetURI sets the URI of the request.
func (r *Request) SetURI(string) {

}

// URL returns `engine.URL`.
func (r *Request) URL() engine.URL {
	return r.url
}

// Header returns `engine.Header`.
func (r *Request) Header() engine.Header {
	return r.header
}

// Proto returns the HTTP proto. (HTTP/1.1 etc.)
func (r *Request) Proto() string {
	return ``
}

// ProtoMajor() int
// ProtoMinor() int

// RemoteAddress returns the client's network address.
func (r *Request) RemoteAddress() string {
	return `127.0.0.1`
}

// RealIP returns the client's network address based on `X-Forwarded-For`
// or `X-Real-IP` request header.
func (r *Request) RealIP() string {
	return `127.0.0.1`
}

// Method returns the request's HTTP function.
func (r *Request) Method() string {
	return `MOCK`
}

// SetMethod sets the HTTP method of the request.
func (r *Request) SetMethod(string) {
}

// Body returns request's body.
func (r *Request) Body() io.ReadCloser {
	return nil
}

func (r *Request) SetBody(io.Reader) {
}

// FormValue returns the form field value for the provided name.
func (r *Request) FormValue(string) string {
	return ``
}
func (r *Request) Object() interface{} {
	return nil
}

func (r *Request) Form() engine.URLValuer {
	return r.form
}

func (r *Request) PostForm() engine.URLValuer {
	return r.form
}

// MultipartForm returns the multipart form.
func (r *Request) MultipartForm() *multipart.Form {
	return nil
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (r *Request) IsTLS() bool {
	return false
}
func (r *Request) Cookie(string) string {
	return ``
}
func (r *Request) Referer() string {
	return ``
}

// UserAgent returns the client's `User-Agent`.
func (r *Request) UserAgent() string {
	return ``
}

// FormFile returns the multipart form file for the provided name.
func (r *Request) FormFile(string) (multipart.File, *multipart.FileHeader, error) {
	return nil, nil, ErrNotImplemented
}

// Size returns the size of request's body.
func (r *Request) Size() int64 {
	return 0
}

func (r *Request) BasicAuth() (string, string, bool) {
	return ``, ``, false
}

func (r *Request) StdRequest() *http.Request {
	return nil
}
