package testing

import (
	"net/http"
	"net/http/httptest"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

func Request(method, path string, handler engine.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(standard.NewRequest(req), standard.NewResponse(rec, req, nil))
	//rec.Code, rec.Body.String(),rec.Header
	return rec
}
