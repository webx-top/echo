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
	api.Use(mw.Log())
	api.Use(mw.Recover())

	s.Add("api@api.localhost:1323", api)

	api.Get("/", func(c echo.Context) error {
		return c.String("API")
	})

	//------
	// Blog
	//------

	blog := echo.New()
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
	site.Use(mw.Log())
	site.Use(mw.Recover())

	s.Add("@localhost:1323", site)

	site.Get("/", func(c echo.Context) error {
		return c.String("Welcome!")
	})

	s.Run(":1323")
}
