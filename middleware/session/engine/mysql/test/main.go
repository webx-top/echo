package main

import (
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/encoding/dbconfig"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/middleware/session/engine/mysql"
)

func main() {
	e := echo.New()
	sessionOptions := &echo.SessionOptions{
		Engine: `mysql`,
		Name:   `SESSIONID`,
		CookieOptions: &echo.CookieOptions{
			Path:     `/`,
			Domain:   ``,
			MaxAge:   0,
			Secure:   false,
			HttpOnly: true,
		},
	}
	mysql.RegWithOptions(&mysql.Options{
		Config: dbconfig.Config{
			Host:    `127.0.0.1`,
			Engine:  `mysql`,
			User:    `root`,
			Pass:    `root`,
			Name:    `test`,
			Charset: `utf8`,
		},
		Table: `session`,
		KeyPairs: [][]byte{
			[]byte(`123456789012345678901234567890ab`),
		},
	})
	e.Use(session.Sessions(sessionOptions))

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
