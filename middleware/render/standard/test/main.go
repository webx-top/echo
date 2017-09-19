package main

import (
	"fmt"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
)

type Nested struct {
	Name     string
	Email    string
	Id       int
	HasChild bool
	Children []*Nested
}

func main() {
	tpl := render.New("standard", "./template/")
	tpl.Init()
	demo := map[string]interface{}{
		"test": "one---",
		"r":    []string{"one", "two", "three"},
		"nested": []*Nested{
			&Nested{
				Name:     `AAA`,
				Email:    `AAA@webx.top`,
				Id:       1,
				HasChild: true,
				Children: []*Nested{
					&Nested{
						Name:     `AAA1`,
						Email:    `AAA1@webx.top`,
						Id:       11,
						HasChild: true,
						Children: []*Nested{
							&Nested{
								Name:     `AAA11`,
								Email:    `AAA11@webx.top`,
								Id:       111,
								HasChild: false,
							},
						},
					},
				},
			},
			&Nested{
				Name:     `BBB`,
				Email:    `BBB@webx.top`,
				Id:       2,
				HasChild: true,
				Children: []*Nested{
					&Nested{
						Name:     `BBB1`,
						Email:    `BBB1@webx.top`,
						Id:       21,
						HasChild: true,
						Children: []*Nested{
							&Nested{
								Name:     `BBB11`,
								Email:    `BBB11@webx.top`,
								Id:       211,
								HasChild: false,
							},
						},
					},
				},
			},
		},
	}

	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := tpl.Fetch("test", demo, nil)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}

	_ = fmt.Printf
	defaults.Use(mw.Log(), mw.Recover(), render.Middleware(tpl))
	defaults.Get(`/`, func(ctx echo.Context) error {
		return ctx.Render(`test`, demo)
	})
	defaults.Run(standard.New(`:4444`))

}
