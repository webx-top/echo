package render_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render"
	test "github.com/webx-top/echo/testing"
)

func request(method, path string, e *echo.Echo) (int, string) {
	rec := test.Request(method, path, e, func(r *http.Request) {
		r.Header.Set(`Accept`, `application/json`)
	})
	return rec.Code, rec.Body.String()
}

func TestEchoMiddleware(t *testing.T) {
	e := echo.New()
	buf := new(bytes.Buffer)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("0")
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

	c, b := request(echo.GET, "/", e)
	assert.Equal(t, "01", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `{"Code":1,"State":"Success","Info":null,"Data":"OK"}`, b)
}
