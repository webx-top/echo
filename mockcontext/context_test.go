package mockcontext

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := Acquire()
	defer Release(ctx)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			ctx.Request().Form().Set(`key`, strconv.Itoa(i))
			ctx.Response().Write([]byte(`co_` + strconv.Itoa(i) + "\n"))
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(`content`, string(ctx.Response().Body()))
	fmt.Println(`reset`, string(ctx.Response().Body()))
}
