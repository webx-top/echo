package mock

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/webx-top/echo/engine"
)

type Response struct {
	header *Header
}

func NewResponse() *Response {
	return &Response{
		header: &Header{},
	}
}

func (r *Response) Header() engine.Header {
	return r.header
}

func (r *Response) WriteHeader(code int) {
}

func (r *Response) KeepBody(on bool) {
}

func (r *Response) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (r *Response) Status() int {
	return 0
}

func (r *Response) Size() int64 {
	return 0
}

func (r *Response) Committed() bool {
	return false
}

func (r *Response) SetWriter(w io.Writer) {
}

func (r *Response) Writer() io.Writer {
	return ioutil.Discard
}

func (r *Response) Object() interface{} {
	return nil
}

func (r *Response) Error(errMsg string, args ...int) {
	r.Write(engine.Str2bytes(errMsg))
}

func (r *Response) Hijack(fn func(net.Conn)) {
}

func (r *Response) Body() []byte {
	return nil
}

func (r *Response) Redirect(url string, code int) {
}

func (r *Response) NotFound() {
}

func (r *Response) SetCookie(cookie *http.Cookie) {
}

func (r *Response) ServeFile(file string) {
}

func (r *Response) Stream(step func(io.Writer) bool) {
}

func (r *Response) StdResponseWriter() http.ResponseWriter {
	return nil
}
