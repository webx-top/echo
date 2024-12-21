package testing

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

type typeCtxKey string

var ctxKeyMock typeCtxKey = `$mock$`

func IsMock(c context.Context) bool {
	v, _ := c.Value(ctxKeyMock).(bool)
	return v
}

func ContextWithMockTag(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKeyMock, true)
}

// Request testing
func Request(method, path string, handler engine.Handler, reqRewrite ...func(*http.Request)) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err == nil {
		req = req.WithContext(ContextWithMockTag(req.Context()))
	}
	for _, rfn := range reqRewrite {
		rfn(req)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(WrapRequest(req), WrapResponse(req, rec))
	//rec.Code, rec.Body.String(),rec.Header
	return rec
}

func NewStdRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	return req
}

func NewStdResponse() http.ResponseWriter {
	return httptest.NewRecorder()
}

func NewRequestAndResponse(method, path string) (engine.Request, engine.Response) {
	req := NewStdRequest(method, path)
	return WrapRequest(req), WrapResponse(req, NewStdResponse())
}

func WrapRequest(req *http.Request) engine.Request {
	return standard.NewRequest(req)
}

func WrapResponse(req *http.Request, rw http.ResponseWriter) engine.Response {
	return standard.NewResponse(rw, req, log.New().Sync())
}
