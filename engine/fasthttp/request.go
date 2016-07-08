// +build !appengine

package fasthttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/admpub/fasthttp"
	"github.com/webx-top/echo/engine"
)

type (
	Request struct {
		context *fasthttp.RequestCtx
		url     engine.URL
		header  engine.Header
		value   *Value
	}
)

func NewRequest(c *fasthttp.RequestCtx) *Request {
	req := &Request{
		context: c,
		url:     &URL{url: c.URI()},
		header:  &RequestHeader{&c.Request.Header},
	}
	req.value = NewValue(req)
	return req
}

func (r *Request) Host() string {
	return string(r.context.Host())
}

func (r *Request) URI() string {
	return string(r.context.RequestURI())
}

func (r *Request) URL() engine.URL {
	return r.url
}

func (r *Request) Header() engine.Header {
	return r.header
}

func (r *Request) Proto() string {
	return "HTTP/1.1"
}

func (r *Request) RemoteAddress() string {
	return r.context.RemoteAddr().String()
}

func (r *Request) Method() string {
	return string(r.context.Method())
}

func (r *Request) SetMethod(method string) {
	r.context.Request.Header.SetMethod(method)
}

func (r *Request) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(r.context.PostBody()))
}

// SetBody implements `engine.Request#SetBody` function.
func (r *Request) SetBody(reader io.Reader) {
	r.context.Request.SetBodyStream(reader, 0)
}

func (r *Request) FormValue(name string) string {
	return string(r.context.FormValue(name))
}

func (r *Request) Form() engine.URLValuer {
	return r.value
}

func (r *Request) PostForm() engine.URLValuer {
	return r.value.postArgs
}

func (r *Request) MultipartForm() *multipart.Form {
	if string(r.context.Request.Header.ContentType()) != "multipart/form-data" {
		return nil
	}
	re, err := r.context.MultipartForm()
	if err != nil {
		r.context.Logger().Printf(err.Error())
	}
	return re
}

func (r *Request) IsTLS() bool {
	return r.context.IsTLS()
}

func (r *Request) Cookie(key string) string {
	return string(r.context.Request.Header.Cookie(key))
}

func (r *Request) Referer() string {
	return string(r.context.Referer())
}

func (r *Request) UserAgent() string {
	return string(r.context.UserAgent())
}

func (r *Request) Object() interface{} {
	return r.context
}

func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	fileHeader, err := r.context.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	var file multipart.File
	file, err = fileHeader.Open()
	return file, fileHeader, err
}

func (r *Request) Scheme() string {
	return string(r.context.URI().Scheme())
}

// Size implements `engine.Request#ContentLength` function.
func (r *Request) Size() int64 {
	return int64(r.context.Request.Header.ContentLength())
}

func (r *Request) reset(c *fasthttp.RequestCtx, h engine.Header, u engine.URL) {
	r.context = c
	r.header = h
	r.url = u
	r.value = NewValue(r)
}
