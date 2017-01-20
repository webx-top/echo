package main

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	//"github.com/webx-top/echo/engine/fasthttp"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	_ "github.com/webx-top/echo/middleware/render/sse"
)

func main() {
	e := echo.New()
	e.Use(mw.Log(), mw.Recover())
	e.Use(render.Middleware(render.New(`sse`, ``)))

	e.Get("/room/:roomid", roomGET)
	e.Post("/room/:roomid", roomPOST)
	e.Delete("/room/:roomid", roomDELETE)
	e.Get("/stream/:roomid", stream)

	e.Run(standard.New(":8080"))
	//e.Run(fasthttp.New(":8080")) //not implemented
}

func stream(c echo.Context) error {
	roomid := c.Param("roomid")
	listener := openListener(roomid)
	defer closeListener(roomid, listener)

	c.Stream(func(w io.Writer) bool {
		c.Render("message", <-listener)
		return true
	})
	return nil
}

func roomGET(c echo.Context) error {
	roomid := c.Param("roomid")
	userid := fmt.Sprint(rand.Int31())
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(200)
	return html.Execute(c.Response(), echo.H{
		"roomid": roomid,
		"userid": userid,
	})
}

func roomPOST(c echo.Context) error {
	roomid := c.Param("roomid")
	userid := c.Form("user")
	message := c.Form("message")
	room(roomid).Submit(userid + ": " + message)

	return c.JSON(echo.H{
		"status":  "success",
		"message": message,
	})
}

func roomDELETE(c echo.Context) error {
	roomid := c.Param("roomid")
	deleteBroadcast(roomid)
	return nil
}
