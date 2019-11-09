package mvc

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
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
	funcMap := map[string]interface{}{}
	s.FuncMapCopyTo(funcMap)
	tmplConfig := &render.Config{}
	//==================copy from Module#InitRenderer()
	blog.Renderer, blog.Resource = blog.NewRenderer(tmplConfig, blog, funcMap)
	if blog.Handler != nil {
		blog.Handler.SetRenderer(blog.Renderer)
		blog.Handler.Get(blog.Resource.Path+`/*`, func(c echo.Context) error {
			return c.String(c.P(0))
		})
	} else {
		blog.Group.Get(strings.TrimPrefix(blog.Resource.Path, `/`+blog.Name)+`/*`, func(c echo.Context) error {
			return c.String(c.P(0))
		})
	}
	blog.Use(mw.SimpleFuncMap(funcMap))
	//===================================================
	//blog.InitRenderer(tmplConfig, funcMap)

	blog.Register(`/index`, func(ctx *Context) error {
		return ctx.String(`OK:/blog/index`)
	})

	s.Commit()

	c, b := request(echo.GET, "/blog/index", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK:/blog/index", b)

	// assets
	c, b = request(echo.GET, "/blog/assets/webx.js", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "webx.js", b)

	// change domain
	c, b = request(echo.GET, "/setdomain?domain=blog.webx.top&module=blog", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	c, b = request(echo.GET, "/blog/index", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)

	// assets
	c, b = request(echo.GET, "/blog/assets/webx.js", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)

	// assets
	c, b = request(echo.GET, "http://blog.webx.top/assets/webx.js", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "webx.js", b)

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

	// assets
	c, b = request(echo.GET, "http://blog.webx.top/assets/webx.js", s)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)

	// assets
	c, b = request(echo.GET, "/blog/assets/webx.js", s)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "webx.js", b)

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
	s.Commit()
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
