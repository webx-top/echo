package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/test"
)

func TestRecover(t *testing.T) {
	e := echo.New()
	e.SetDebug(true)
	req := test.NewRequest(echo.GET, "/", nil)
	rec := test.NewResponseRecorder()
	c := echo.NewContext(req, rec, e)
	h := Recover()(echo.HandlerFunc(func(c echo.Context) error {
		panic("test")
	}))
	h.Handle(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Status())
	assert.Contains(t, rec.Body.String(), "panic recover")
}
