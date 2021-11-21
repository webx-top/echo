package defaults

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/admpub/fasthttp"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	fasthttpng "github.com/webx-top/echo/engine/fasthttp"
)

func TestMustGetContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), `testKey`, `testVal`)
	eCtx := MustGetContext(ctx)
	assert.Equal(t, `testVal`, eCtx.Value(`testKey`))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			eCtx.SetValue(`co_`+strconv.Itoa(i), i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 50; i++ {
		assert.Equal(t, i, eCtx.Value(`co_`+strconv.Itoa(i)))
	}
}

func TestFastHTTPContext(t *testing.T) {
	rCtx := &fasthttp.RequestCtx{}
	eCtx := echo.NewContext(fasthttpng.NewRequest(rCtx), fasthttpng.NewResponse(rCtx), Default)
	eCtx.SetValue(`testKey`, `testVal`)
	assert.Equal(t, `testVal`, eCtx.Value(`testKey`))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			eCtx.SetValue(`co_`+strconv.Itoa(i), i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 50; i++ {
		assert.Equal(t, i, eCtx.Value(`co_`+strconv.Itoa(i)))
	}
}
