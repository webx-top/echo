package echo_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
)

func request(method, path string, e *Echo) (int, string) {
	req, _ := http.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(standard.NewRequest(req), standard.NewResponse(rec, req, nil))
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

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("3")
			return next.Handle(c)
		}
	})

	// Route
	e.Get("/", func(c Context) error {
		return c.String("OK")
	})

	c, b := request(GET, "/", e)
	assert.Equal(t, "-1123", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoMiddlewareError(t *testing.T) {
	e := New()
	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return errors.New("error")
		}
	})
	e.Get("/", NotFoundHandler)
	c, _ := request(GET, "/", e)
	assert.Equal(t, http.StatusInternalServerError, c)
}

func TestEchoHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.Get("/ok", func(c Context) error {
		return c.String("OK")
	})

	c, b := request(GET, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoMeta(t *testing.T) {
	e := New()

	middleware := e.MetaMiddleware(
		H{"authorization": true, "data": H{"by": "middleware"}},
		func(next HandlerFunc) HandlerFunc {
			return func(c Context) error {
				return next.Handle(c)
			}
		},
	)
	g := e.Group("/root")
	g.Use(middleware)

	g.Get("/", e.MetaHandler(
		H{"version": 1.0, "data": H{"by": "handler"}},
		func(c Context) error {
			//Dump(c.Meta())
			return c.String("OK")
		},
	))

	var meta H

	for _, route := range e.Routes() {
		if route.Path == "/root/" {
			meta = route.Meta
		}
	}
	assert.Equal(t, H{
		"authorization": true,
		"version":       1.0,
		"data": H{
			"by": "handler",
		},
	}, meta)

	c, b := request(GET, "/root/", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}
