package subdomains

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	test "github.com/webx-top/echo/testing"
)

func request(method, path string, h engine.Handler, reqRewrite ...func(*http.Request)) (int, string) {
	rec := test.Request(method, path, h, reqRewrite...)
	return rec.Code, rec.Body.String()
}
func TestSortHosts(t *testing.T) {
	a := New()
	e := echo.New()
	e.Get(`/`, func(c echo.Context) error {
		return c.String(`frontend`)
	})
	a.Add(`frontend`, e)
	e2 := echo.New()
	e2.SetPrefix(`/admin`)
	e2.Get(``, func(c echo.Context) error {
		return c.String(`backend`)
	})
	e2.Get(`/index`, func(c echo.Context) error {
		return c.String(`backend-index`)
	})
	a.Add(`backend`, e2)

	a.Ready().Commit()
	assert.Equal(t, []string{`backend`, `frontend`}, a.Hosts[``])

	c, b := request(echo.GET, "/", a)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "frontend", b)

	c, b = request(echo.GET, "/admin", a)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "backend", b)

	c, b = request(echo.GET, "/admin/index", a)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "backend-index", b)

	c, b = request(echo.GET, "/adminindex", a)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)
}
