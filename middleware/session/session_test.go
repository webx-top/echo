package session_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session"
	test "github.com/webx-top/echo/testing"
)

func request(method, path string, e *echo.Echo, reqRewrite ...func(*http.Request)) (int, string, http.Header) {
	rec := test.Request(method, path, e, reqRewrite...)
	return rec.Code, rec.Body.String(), rec.HeaderMap
}

func TestSession(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(nil))
	e.Get(`/`, func(ctx echo.Context) error {
		//echo.Dump(ctx.Request().Header().Std())
		i, _ := ctx.Session().Get(`count`).(int)
		i++
		ctx.Session().Set(`count`, i)
		ctx.SetCookie(`user`, `test-`+strconv.Itoa(i))
		return ctx.String(`ok`)
	})
	e.Get(`/result`, func(ctx echo.Context) error {
		return ctx.String(fmt.Sprintf(`%v:%v`, ctx.Session().Get(`count`), ctx.GetCookie(`user`)))
	})

	headers := http.Header{}
	rew := func(headers http.Header) func(req *http.Request) {
		return func(req *http.Request) {
			for _, h := range headers["Set-Cookie"] {
				req.Header.Add(`Cookie`, h)
			}
		}
	}
	for i := 1; i < 5; i++ {
		code, resp, header := request(`GET`, `/`, e, rew(headers))
		assert.Equal(t, 200, code)
		assert.Equal(t, `ok`, resp)
		assert.Equal(t, `user=test-`+strconv.Itoa(i)+`; Path=/`, header["Set-Cookie"][0])
		assert.Equal(t, `SID=`, header["Set-Cookie"][1][0:4])
		headers = header
		code, resp, _ = request(`GET`, `/result`, e, rew(headers))
		assert.Equal(t, 200, code)
		assert.Equal(t, strconv.Itoa(i)+`:test-`+strconv.Itoa(i), resp)
	}
}
