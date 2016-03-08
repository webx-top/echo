// +build !appengine

package fasthttp

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/admpub/fasthttp"
	"github.com/labstack/gommon/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Response struct {
		context   *fasthttp.RequestCtx
		header    engine.Header
		status    int
		size      int64
		committed bool
		writer    io.Writer
		logger    logger.Logger
	}
)

func NewResponse(c *fasthttp.RequestCtx) *Response {
	return &Response{
		context: c,
		header:  &ResponseHeader{&c.Response.Header},
		writer:  c,
		logger:  log.New("echo"),
	}
}

func (r *Response) Object() interface{} {
	return r.context
}

func (r *Response) Header() engine.Header {
	return r.header
}

func (r *Response) WriteHeader(code int) {
	if r.committed {
		r.logger.Warn("response already committed")
		return
	}
	r.status = code
	r.context.SetStatusCode(code)
	r.committed = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.writer.Write(b)
	r.size += int64(n)
	return
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) Size() int64 {
	return r.size
}

func (r *Response) Committed() bool {
	return r.committed
}

func (r *Response) SetWriter(w io.Writer) {
	r.writer = w
}

func (r *Response) Writer() io.Writer {
	return r.writer
}

func (r *Response) Hijack(fn func(net.Conn)) {
	r.context.Hijack(fasthttp.HijackHandler(fn))
}

func (r *Response) Body() []byte {
	switch strings.ToLower(r.header.Get(`Content-Encoding`)) {
	case `gzip`:
		body, err := r.context.Response.BodyGunzip()
		if err != nil {
			r.logger.Error(err)
		}
		return body
	case `deflate`:
		body, err := r.context.Response.BodyInflate()
		if err != nil {
			r.logger.Error(err)
		}
		return body
	default:
		return r.context.Response.Body()
	}
}

func (r *Response) Redirect(url string, code int) {
	r.context.Redirect(url, code)
}

func (r *Response) NotFound() {
	r.context.NotFound()
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.header.Set("Set-Cookie", cookie.String())
}

func (r *Response) ServeFile(file string) {
	fasthttp.ServeFile(r.context, file)
}

func (r *Response) reset(c *fasthttp.RequestCtx, h engine.Header) {
	r.context = c
	r.header = h
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
	r.writer = c
}