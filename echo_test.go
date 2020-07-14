package echo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/admpub/log"
	"github.com/stretchr/testify/assert"

	"github.com/webx-top/echo"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	mw "github.com/webx-top/echo/middleware"
	test "github.com/webx-top/echo/testing"
	"github.com/webx-top/validation"
)

func init() {
	mw.DefaultLogWriter = log.Writer(log.LevelInfo)
}

func request(method, path string, e *Echo, reqRewrite ...func(*http.Request)) (int, string) {
	rec := test.Request(method, path, e, reqRewrite...)
	return rec.Code, rec.Body.String()
}

func TestEchoMiddleware(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)

	e.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-1")
			return next.Handle(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("1")
			return next.Handle(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("2")
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
	e.Use(mw.Log(), func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return errors.New("error")
		}
	})
	e.Get("/", NotFoundHandler)

	e.RebuildRouter()

	c, _ := request(GET, "/", e)
	assert.Equal(t, http.StatusInternalServerError, c)
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
	g := e.Host(".admpub.com")
	g.Get("/host", func(c Context) error {
		if c.Queryx(`route`).Bool() {
			return c.JSON(c.Route())
		}
		return c.String(c.Host())
	}).SetName(`host`)
	g = e.Host("<uid:[0-9]+>.<name>.com").SetAlias(`user`)
	g.Get("/host2", func(c Context) error {
		if c.Queryx(`route`).Bool() {
			return c.JSON(c.Route())
		}
		return c.String(c.HostParam(`uid`) + `.` + c.HostParam(`name`) + `.com`)
	}).SetName(`host2`)

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
	c, b = request(GET, "/host", e, func(req *http.Request) {
		req.Host = "test.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "test.admpub.com", b)
	c, b = request(GET, "/host", e, func(req *http.Request) {
		req.Host = "test-b.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "test-b.admpub.com", b)
	c, b = request(GET, "/host?route=1", e, func(req *http.Request) {
		req.Host = "test-b.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `{"Host":".admpub.com","Method":"GET","Path":"/host","Name":"host","Format":"/host","Params":[],"Prefix":"","Meta":{}}`, b)

	c, b = request(GET, "/host2", e, func(req *http.Request) {
		req.Host = "123.coscms.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "123.coscms.com", b)
	assert.Equal(t, "10000.admpub.com/host2", e.TypeHost(`user`, 10000, `admpub`).URI(`host2`))
	assert.Equal(t, "10001.admpub.com/host2", e.TypeHost(`user`, echo.H{`uid`: 10001, `name`: `admpub`}).URI(`host2`))
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

type MetaRequest struct {
	Name string `valid:"required"`
}

func NewMetaRequest() echo.MetaValidator {
	return &MetaRequest{}
}

func (m *MetaRequest) Validate(c echo.Context) error {
	v := validation.New()
	pass, err := v.Valid(m)
	if err != nil {
		return err
	}
	if !pass {
		verr := v.Error()
		err = c.NewError(code.InvalidParameter, verr.WithField().Error()).SetZone(verr.Field)
	}
	return err
}

func (m *MetaRequest) Filters(c echo.Context) []FormDataFilter {
	return nil
}

func TestEchoMeta(t *testing.T) {
	e := New()
	e.SetDebug(true)
	g := e.Group("/root")

	g.Get("/", e.MetaHandler(
		H{"version": 1.0, "data": H{"by": "handler"}},
		func(c Context) error {
			return c.JSON(c.Route().Meta)
		},
	))

	g.Post("/post", e.MetaHandler(
		nil,
		func(c Context) error {
			data := c.Internal().Get(`validated`).(*MetaRequest)
			return c.String(data.Name)
		},
		NewMetaRequest,
	))

	e.RebuildRouter()

	var meta H

	for _, route := range e.Routes() {
		if route.Path == "/root/" {
			meta = route.Meta
		}
	}
	expected := H{
		"version": 1.0,
		"data": H{
			"by": "handler",
		},
	}
	assert.Equal(t, expected, meta)

	c, b := request(GET, "/root/", e)
	assert.Equal(t, http.StatusOK, c)
	expected2, _ := json.MarshalIndent(expected, "", "  ")
	assert.Equal(t, string(expected2), b)

	c, b = request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)
}

func TestEchoData(t *testing.T) {
	data := NewData(nil)
	data.SetCode(0)
	assert.Equal(t, 0, data.Code.Int())
	assert.Equal(t, `0`, fmt.Sprintf(`%d`, data.Code))
	assert.Equal(t, `Failure`, fmt.Sprintf(`%v`, data.Code))
	assert.Equal(t, `Failure`, fmt.Sprintf(`%s`, data.Code))
	assert.Equal(t, `Failure`, data.State)
}
