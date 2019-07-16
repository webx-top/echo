package test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	testings "github.com/webx-top/echo/testing"
)

type HandlerTest func(*bytes.Buffer) echo.HandlerFunc
type MiddlewareTest func(*bytes.Buffer) echo.MiddlewareFuncd

func Hit(t *testing.T, configs []*Config, middelwares ...MiddlewareTest) {
	e := echo.New()
	buf := new(bytes.Buffer)
	for _, h := range middelwares {
		e.Use(h(buf))
	}
	for _, cfg := range configs {
		ms := make([]interface{}, len(cfg.Middlewares))
		for k, m := range cfg.Middlewares {
			ms[k] = m(buf)
		}
		e.Match([]string{cfg.Method}, cfg.Path, cfg.Handler(buf), ms...)
		r := testings.Request(cfg.Method, cfg.Path, e, cfg.ReqRewrite...)
		//assert.Equal(t, "-1123", buf.String())
		//assert.Equal(t, http.StatusOK, r.Code)
		//assert.Equal(t, "OK", r.Body.String())
		cfg.Checker(t, r, buf)
	}
}

var (
	Eq          = assert.Equal
	NotEq       = assert.NotEqual
	True        = assert.True
	False       = assert.False
	NotNil      = assert.NotNil
	Empty       = assert.Empty
	NotEmpty    = assert.NotEmpty
	Len         = assert.Len
	Contains    = assert.Contains
	NotContains = assert.NotContains
	Subset      = assert.Subset
	NotSubset   = assert.NotSubset
)
