package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/webx-top/echo"
)

func TestTest(t *testing.T) {
	configs := []*Config{
		{
			Method: `GET`,
			Handler: func(buf *bytes.Buffer) echo.HandlerFunc {
				return func(c echo.Context) error {
					buf.WriteString(`1`)
					return c.String(`OK`)
				}
			},
			Checker: func(t *testing.T, r *httptest.ResponseRecorder, buf *bytes.Buffer) {
				Eq(t, `1`, buf.String())
				Eq(t, `OK`, r.Body.String())
				Eq(t, http.StatusOK, r.Code)
			},
		},
		{
			Method: `GET`,
			Path:   `/2`,
			Handler: func(buf *bytes.Buffer) echo.HandlerFunc {
				return func(c echo.Context) error {
					buf.WriteString(`2`)
					return c.String(`OK`)
				}
			},
			Checker: DefaultChecker(`12`),
		},
	}
	Hit(t, configs)
}
