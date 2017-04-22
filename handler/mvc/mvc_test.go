package mvc

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	test "github.com/webx-top/echo/testing"
)

func request(method, path string, e *Application) (int, string) {
	rec := test.Request(method, path, e)
	return rec.Code, rec.Body.String()
}

func TestChangeDomain(t *testing.T) {
	s := New(`test`)
	s.Debug(true)
	base := s.NewModule("base")

	base.Register(`/setdomain`, func(ctx *Context) error {
		domain := ctx.Query(`domain`)
		module := ctx.Query(`module`)
		if len(module) == 0 {
			return errors.New(`Module name is required`)
		}
		if s.HasModule(module) == false {
			return errors.New(`not found module: ` + module)
		}
		s.SetDomain(module, domain)
		return ctx.String(`OK`)
	})

	blog := s.NewModule("blog")

	blog.Register(`/index`, func(ctx *Context) error {
		return ctx.String(`OK:/blog/index`)
	})

	c, b := request(echo.GET, "/blog/index", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK:/blog/index", b)

	// change domain
	c, b = request(echo.GET, "/setdomain?domain=blog.webx.top&module=blog", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	c, b = request(echo.GET, "/blog/index", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)

	c, b = request(echo.GET, "http://blog.webx.top/index", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK:/blog/index", b)

	//reverse
	c, b = request(echo.GET, "/setdomain?domain=&module=blog", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	c, b = request(echo.GET, "http://blog.webx.top/index", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)

	c, b = request(echo.GET, "/blog/index", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK:/blog/index", b)
}

func TestMiddleware(t *testing.T) {
	s := New(`testMiddleware`)
	buf := new(bytes.Buffer)

	s.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-1")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("0")
			return next.Handle(c)
		}
	})

	s.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-3")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-2")
			return next.Handle(c)
		}
	})

	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("1")
			return next.Handle(c)
		}
	})

	base := s.NewModule("base", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("2")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("3")
			return next.Handle(c)
		}
	})

	// Route
	base.Router().Get("/", func(c echo.Context) error {
		return c.String("OK")
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("4")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("5")
			return next.Handle(c)
		}
	})

	blog := s.NewModule("blog@blog.webx.top", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("6")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("7")
			return next.Handle(c)
		}
	})

	// Route
	blog.Router().Get("/", func(c echo.Context) error {
		return c.String("OK")
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("8")
			return next.Handle(c)
		}
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("9")
			return next.Handle(c)
		}
	})
	c, b := request(echo.GET, "/", s)
	assert.Equal(t, "-3-2-1012345", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	buf = new(bytes.Buffer)
	c, b = request(echo.GET, "http://blog.webx.top/", s)
	assert.Equal(t, "-3-2-1016789", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	c, b = request(echo.GET, "/blog", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)
}
