package main

import (
	"flag"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
)

func main() {
	port := flag.String(`p`, "8080", "port")
	flag.Parse()
	e := echo.New()
	e.Use(mw.Log())

	d := render.New(`standard`, `./template`)
	d.Init(true)

	e.Use(render.Middleware(d))

	e.Get("/", func(c echo.Context) error {

		// It uses template file ./template/index.html
		return c.Render(`index`, map[string]interface{}{
			"Name": "Webx",
		})
	})

	switch `` {
	case `fast`:
		// FastHTTP
		e.Run(fasthttp.New(":" + *port))

	default:
		// Standard
		e.Run(standard.New(":" + *port))
	}
}
