package mock

import (
	"bytes"
	"io"
	"net"
	"net/http"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

type Response struct {
	header    engine.Header
	status    int
	size      int64
	committed bool
	writer    io.Writer
	body      []byte
	keepBody  bool
}

func NewResponse(writers ...io.Writer) *Response {
	var writer io.Writer
	if len(writers) > 0 {
		writer = writers[0]
	}
	if writer == nil {
		writer = bytes.NewBuffer(nil)
	}
	return &Response{
		header: &standard.Header{Header: http.Header{}},
		writer: writer,
	}
}

func (r *Response) Header() engine.Header {
	return r.header
}

func (r *Response) WriteHeader(code int) {
	if r.committed {
		log.Warn("response already committed")
		return
	}
	r.status = code
	r.committed = true
}

func (r *Response) KeepBody(on bool) {
	r.keepBody = on
}

func (r *Response) Write(b []byte) (n int, err error) {
	if !r.committed {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.WriteHeader(r.status)
	}
	if r.keepBody {
		r.body = append(r.body, b...)
	}
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

func (r *Response) Object() interface{} {
	return nil
}

func (r *Response) Error(errMsg string, args ...int) {
	if len(args) > 0 {
		r.status = args[0]
	} else {
		r.status = http.StatusInternalServerError
	}
	r.Write(engine.Str2bytes(errMsg))
	r.WriteHeader(r.status)
}

func (r *Response) Hijack(fn func(net.Conn)) {
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Redirect(url string, code int) {
}

func (r *Response) NotFound() {
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.header.Add(engine.HeaderSetCookie, cookie.String())
}

func (r *Response) ServeFile(file string) {
}

func (r *Response) Stream(step func(io.Writer) bool) {
}

func (r *Response) StdResponseWriter() http.ResponseWriter {
	return nil
}
