package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/admpub/jet/v6"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/pprof"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	. "github.com/webx-top/echo/middleware/render/jet"
)

func main() {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	heapAllocStart := memStat.HeapAlloc
	jet.SetDefaultExtensions(`.html`)
	t := New(`./template/`)
	t.Init()
	//t.SetDebug(true)
	ctx := echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
	for i := 0; i < 5000; i++ {
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
	runtime.ReadMemStats(memStat)
	heapAllocEnd := memStat.HeapAlloc
	heapAllocIncr := heapAllocEnd - heapAllocStart
	fmt.Println(`~~~~~~~~~~~~~~~~~~>heapAllocIncr: `, com.FormatBytes(heapAllocIncr))

	defaults.Use(mw.Log(), mw.Recover(), render.Middleware(t))
	pprof.Wrapper(defaults.Default)
	defaults.Run(standard.New(`:4444`))
}
