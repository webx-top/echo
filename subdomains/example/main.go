package main

import (
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/subdomains"
)

func main() {
	s := subdomains.New()

	//-----
	// API
	//-----

	api := echo.New()
	api.SetDebug(true)
	api.Use(mw.Log())
	api.Use(mw.Recover())

	s.Add("api@api.localhost:1323", api)

	api.Get("/", func(c echo.Context) error {
		return c.String("API")
	})
	api.Get("/geturl/:name", func(c echo.Context) error {
		return c.String("GET-URL:" + s.URLByName(`#api#geturl`, c.Param("name")))
	}).SetName(`geturl`)

	//------
	// Blog
	//------

	blog := echo.New()
	blog.SetDebug(true)
	blog.Use(mw.Log())
	blog.Use(mw.Recover())

	s.Add("blog@blog.localhost:1323", blog)

	blog.Get("/", func(c echo.Context) error {
		return c.String("Blog")
	})

	//---------
	// Website
	//---------

	site := echo.New()
	site.SetDebug(true)
	site.Use(mw.Log())
	site.Use(mw.Recover())

	s.Add("@localhost:1323", site)

	site.Get("/", func(c echo.Context) error {
		return c.String("Welcome!")
	})

	s.Run(":1323")
}
