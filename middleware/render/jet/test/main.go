package main

import (
	"fmt"
	"time"

	"github.com/admpub/jet"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/mock"
	. "github.com/webx-top/echo/middleware/render/jet"
)

func main() {
	jet.SetDefaultExtensions(`.html`)
	t := New(`./template/`)
	t.Init()
	//t.SetDebug(true)
	ctx := echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := t.Fetch("test", map[string]interface{}{
			`name`: `webx`,
			"test": "times---" + fmt.Sprintf("%v", i),
			"r":    []string{"one", "two", "three"},
		}, ctx)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}
}
