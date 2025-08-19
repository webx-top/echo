package render_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/render"
	test "github.com/webx-top/echo/testing"
)

func requestJSON(method, path string, e *echo.Echo) (int, string) {
	rec := test.Request(method, path, e, func(r *http.Request) {
		r.Header.Set(`Accept`, `application/json`)
	})
	return rec.Code, rec.Body.String()
}

func requestHTML(method, path string, e *echo.Echo) (int, string) {
	rec := test.Request(method, path, e, func(r *http.Request) {
		r.Header.Set(`Accept`, `text/html`)
	})
	return rec.Code, rec.Body.String()
}

func requestJSON2(method, path string, e *echo.Echo) (int, string) {
	rec := test.Request(method, path, e, func(r *http.Request) {
		r.Header.Set(`Accept`, `application/json, text/javascript; q=0.01`)
	})
	return rec.Code, rec.Body.String()
}

func TestEchoMiddleware(t *testing.T) {
	e := echo.New()
	e.SetHTTPErrorHandler(render.HTTPErrorHandler(nil))
	buf := new(bytes.Buffer)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("0")
			buf.WriteString(c.Format())
			return next.Handle(c)
		}
	}, render.Auto(), func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("1")
			return next.Handle(c)
		}
	})

	// Route
	e.Get("/", func(c echo.Context) error {
		return c.Render(`no`, "OK")
	})
	e.Get("/noperm", func(c echo.Context) error {
		return c.NewError(code.Unauthenticated, ``)
	})

	e.RebuildRouter()

	c, b := requestJSON2(echo.GET, "/", e)
	assert.Equal(t, "0json1", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `{"Code":1,"State":"Success","Info":null,"Data":"OK"}`, b)

	buf.Reset()

	c, b = requestJSON2(echo.GET, "/noperm", e)
	assert.Equal(t, "0json1", buf.String())
	assert.Equal(t, http.StatusUnauthorized, c)
	assert.Equal(t, `{"Code":-1,"State":"Unauthenticated","Info":"Unauthenticated","Zone":"","Data":{}}`, b)

	buf.Reset()

	c, b = requestHTML(echo.GET, "/noperm", e)
	assert.Equal(t, "0html1", buf.String())
	assert.Equal(t, http.StatusUnauthorized, c)
	assert.Equal(t, `Unauthenticated`, b)
}
