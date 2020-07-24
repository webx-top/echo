package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/mock"
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
	t := template.New(`C:\a\b\c\c.html`)
	t = template.Must(t.Parse(`{{define "C:\\a\\b\\c\\d.html"}}123333{{end}}{{template "C:\\a\\b\\c\\d.html"}}`)) //注意：define和template标签后面的参数如果含“\”，则会执行转义。所以“C:\d”需要改为“C:\\d”,否则会出错
	t.Execute(os.Stdout, nil)
	//return

	tpl := render.New("standard", "./template/")
	tpl.Init()
	ctx := echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
	demo := map[string]interface{}{
		"test": "one---",
		"r":    []string{"one", "two", "three"},
		"nested": []*Nested{
			{
				Name:     `AAA`,
				Email:    `AAA@webx.top`,
				Id:       1,
				HasChild: true,
				Children: []*Nested{
					{
						Name:     `AAA1`,
						Email:    `AAA1@webx.top`,
						Id:       11,
						HasChild: true,
						Children: []*Nested{
							{
								Name:     `AAA11`,
								Email:    `AAA11@webx.top`,
								Id:       111,
								HasChild: false,
							},
						},
					},
				},
			},
			{
				Name:     `BBB`,
				Email:    `BBB@webx.top`,
				Id:       2,
				HasChild: true,
				Children: []*Nested{
					{
						Name:     `BBB1`,
						Email:    `BBB1@webx.top`,
						Id:       21,
						HasChild: true,
						Children: []*Nested{
							{
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
		str := tpl.Fetch("test", demo, ctx)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}

	_ = fmt.Printf
	defaults.Use(mw.Log(), mw.Recover(), render.Middleware(tpl))
	defaults.Get(`/`, func(ctx echo.Context) error {
		return ctx.Render(`test`, demo)
	})
	defaults.Get(`/e`, func(ctx echo.Context) error {
		return ctx.Render(`test2`, demo)
	})
	defaults.Run(standard.New(`:4444`))

}
