package main

import (
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware/session"
	_ "github.com/webx-top/echo/middleware/session/engine/cookie"
)

func main() {
	e := echo.New()
	e.Use(session.Middleware(nil))
	e.Get(`/`, func(ctx echo.Context) error {
		n, y := ctx.Session().Get(`count`).(int)
		if y {
			n++
		}
		ctx.Session().Set(`count`, n)
		ctx.SetCookie(`user`, `test`)
		return ctx.String(fmt.Sprintf(`session: %v`, ctx.Session().Get(`count`)))
	})
	e.Get(`/result`, func(ctx echo.Context) error {
		n := ctx.Session().Get(`count`)
		u := ctx.GetCookie(`user`)
		return ctx.String(fmt.Sprintf(`session: %v cookie: %s`, n, u))
	})
	e.Run(standard.New(`:4444`))
}
