package main

import (
	"fmt"
	"time"

	. "github.com/webx-top/echo/middleware/render/pongo2"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
)

func main() {
	t := New(`./template/`)
	t.Init()
	demo := map[string]interface{}{
		`name`: `webx`,
		"test": "times---",
		"r":    []string{"one", "two", "three"},
	}
	//t.SetDebug(true)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := t.Fetch("test", demo, nil)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}

	_ = fmt.Printf
	defaults.Use(mw.Log(), mw.Recover(), render.Middleware(t))
	defaults.Get(`/`, func(ctx echo.Context) error {
		return ctx.Render(`test`, demo)
	})
	defaults.Run(standard.New(`:4444`))
}
