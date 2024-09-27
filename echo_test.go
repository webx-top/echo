package echo_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/admpub/log"
	"github.com/stretchr/testify/assert"

	. "github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	test "github.com/webx-top/echo/testing"
)

func init() {
	mw.DefaultLogWriter = log.Writer(log.LevelInfo)
	log.Sync()
}

func request(method, path string, e *Echo, reqRewrite ...func(*http.Request)) (int, string) {
	rec := test.Request(method, path, e, reqRewrite...)
	return rec.Code, rec.Body.String()
}

func TestEchoMiddleware(t *testing.T) {
	e := New()
	e.SetMaxRequestBodySize(2 << 20) // 2M
	buf := new(bytes.Buffer)

	e.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			assert.Equal(t, 2<<20, c.Request().MaxSize())
			buf.WriteString("-1")
			return next.Handle(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("1")
			c.Request().SetMaxSize(3 << 20)
			return next.Handle(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("2")
			assert.Equal(t, 3<<20, c.Request().MaxSize())
			return next.Handle(c)
		}
	})

	e.Use(mw.Log(), func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("3")
			return next.Handle(c)
		}
	})

	// Route
	e.Get("/", func(c Context) error {
		assert.Equal(t, 3<<20, c.Request().MaxSize())
		return c.String("OK")
	})

	e.RebuildRouter()

	c, b := request(GET, "/", e)
	assert.Equal(t, "-1123", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoMiddlewareError(t *testing.T) {
	e := New()
	e.SetDebug(true)
	e.Use(mw.Log(), func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Equal(t, `/?id=5&name=测试`, c.RequestURI())
			assert.Equal(t, `https://github.com/?id=5&name=测试`, c.FullRequestURI())
			assert.Equal(t, ``, c.Request().URI())
			return errors.New("error")
		}
	})
	e.Get("/", NotFoundHandler)

	e.RebuildRouter()

	c, r := request(GET, "https://github.com/?id=5&name=测试", e)
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `error`, r)
}

func TestEchoHandlerError(t *testing.T) {
	e := New()
	e.SetDebug(true)
	e.Use(mw.Log())
	e.Get("/", func(c Context) error {
		return errors.New("error")
	})

	e.RebuildRouter()

	c, r := request(GET, "/", e)
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `error`, r)
}

func TestEchoRoutePath(t *testing.T) {
	e := New()
	e.Use(mw.Log())
	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Equal(t, `/`, c.Path())
			return next.Handle(c)
		}
	})
	e.Get("/", func(ctx Context) error {
		return ctx.String(ctx.Path())
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Equal(t, `/`, c.Path())
			return next.Handle(c)
		}
	})

	e.RebuildRouter()

	c, r := request(GET, "/", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `/`, r)
}

func TestGroupMiddleware(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)

	e.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-1")
			return next.Handle(c)
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("0")
			return next.Handle(c)
		}
	})

	e.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-3")
			return next.Handle(c)
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-2")
			return next.Handle(c)
		}
	})

	e.Use(mw.Log(), func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("1")
			return next.Handle(c)
		}
	})

	g := e.Group("/", func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("2")
			return next.Handle(c)
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("3")
			return next.Handle(c)
		}
	})

	// Route
	g.Get("", func(c Context) error {
		return c.String("OK")
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("4")
			return next.Handle(c)
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("5")
			return next.Handle(c)
		}
	})

	e.RebuildRouter()

	c, b := request(GET, "/", e)
	assert.Equal(t, "-3-2-1012345", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	buf = new(bytes.Buffer)
	e.RebuildRouter()
	c, b = request(GET, "/", e)
	assert.Equal(t, "-3-2-1012345", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.Get("/ok", func(c Context) error {
		return c.String("OK")
	})
	e.Get("/view/:id", func(c Context) error {
		return c.String(c.Param(`id`))
	}).SetName(`view`)
	e.Get("/file/*", func(c Context) error {
		return c.String(c.P(0))
	})
	e.Route(`GET,POST`, "/input", func(c Context) error {
		return c.String(c.Form(`test`))
	})
	e.RebuildRouter()

	assert.Equal(t, `/view/8`, e.URI(`view`, 8))

	c, b := request(GET, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
	c, b = request(GET, "/view/123", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "123", b)
	c, b = request(POST, "/view/0", e)
	assert.Equal(t, http.StatusMethodNotAllowed, c)
	assert.Equal(t, "Method Not Allowed", b)
	c, b = request(GET, "/file/path/to/file.js", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "path/to/file.js", b)
	c, b = request(GET, "/input?test=1", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "1", b)
	c, b = request(POST, "/input", e, func(r *http.Request) {
		r.PostForm = url.Values{
			`test`: []string{`2`},
		}
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "2", b)
}

func TestEchoRouter(t *testing.T) {
	e := New()

	e.Get("/router/:n/list", func(c Context) error {
		//Dump(c.Route())
		return c.String(c.Param(`n`))
	})
	e.RebuildRouter()

	c, b := request(GET, "/router/123/list", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "123", b)
}

func TestEchoRealIP(t *testing.T) {
	e := New()

	e.Get("/", func(c Context) error {
		Dump(c.Request().Header().Std())
		return c.String(c.RealIP())
	})
	e.RebuildRouter()
	ipArr := []string{`137`, `0`, `10`, `1`}
	ipArr2 := []string{`137`, `0`, `10`, `8`}
	c, b := request(GET, "/", e, func(r *http.Request) {
		r.Header.Set(`X-Forwarded-For`, strings.Join(ipArr, `.`))
		r.Header.Set(`X-Real-Ip`, strings.Join(ipArr2, `.`))
		r.RemoteAddr = `127.0.0.1:57092`
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, strings.Join(ipArr, `.`), b)
}

func TestEchoData(t *testing.T) {
	data := NewData(nil)
	data.SetCode(0)
	assert.Equal(t, 0, data.Code.Int())
	assert.Equal(t, `0`, fmt.Sprintf(`%d`, data.Code))
	assert.Equal(t, `Failure`, fmt.Sprintf(`%v`, data.Code))
	assert.Equal(t, `Failure`, data.State)
}

func TestHandlerFuncWithArg(t *testing.T) {
	type testRequest struct{ Name string }
	type testResponse struct{ Author string }
	_ = HandlerFuncWithArg[testRequest, testResponse](func(c Context, r *testRequest) (testResponse, error) {
		return testResponse{Author: r.Name}, nil
	})
}
