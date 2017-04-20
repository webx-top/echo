package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware/session"
)

func request(method, path string, e *echo.Echo) (int, string, http.Header) {
	req, _ := http.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(standard.NewRequest(req), standard.NewResponse(rec, req, nil))
	return rec.Code, rec.Body.String(), rec.HeaderMap
}

func TestSession(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(nil))
	e.Get(`/`, func(ctx echo.Context) error {
		ctx.Session().Set(`count`, 1)
		ctx.SetCookie(`user`, `test`)
		return ctx.String(`ok`)
	})
	code, resp, header := request(`GET`, `/`, e)
	assert.Equal(t, 200, code)
	assert.Equal(t, `ok`, resp)
	assert.Equal(t, `user=test; Path=/`, header["Set-Cookie"][0])
	assert.Equal(t, `SID=`, header["Set-Cookie"][1][0:4])
}
