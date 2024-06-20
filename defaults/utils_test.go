package defaults

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/admpub/fasthttp"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	fasthttpng "github.com/webx-top/echo/engine/fasthttp"
)

type tkey struct{}

var testContextKey tkey

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

	eCtx.Internal().Set(`_`, `0`)
	eCtx.Request().Form().Set(`test`, `1`)

	req := eCtx.Request().StdRequest().WithContext(echo.AsStdContext(eCtx))
	assert.Equal(t, `1`, req.Form.Get(`test`))
	assert.Equal(t, eCtx.Internal().String(`_`), MustGetContext(req.Context()).Internal().String(`_`))

	req.Form.Set(`reqTest`, `req1`)
	req2 := req.WithContext(context.WithValue(req.Context(), testContextKey, `000`))
	assert.Equal(t, `req1`, req2.Form.Get(`reqTest`))
	assert.Equal(t, req.Form.Get(`reqTest`), req2.Form.Get(`reqTest`))
	assert.Equal(t, `000`, req2.Context().Value(testContextKey).(string))
	assert.Equal(t, `*context.valueCtx`, fmt.Sprintf("%T", req2.Context()))
	assert.Equal(t, `*echo.xContext`, fmt.Sprintf("%T", MustGetContext(req2.Context())))
	assert.Equal(t, `0`, MustGetContext(req2.Context()).Internal().String(`_`))

	req3 := req.WithContext(context.WithValue(req.Context(), testContextKey, `000`))
	assert.Equal(t, `req1`, req3.Form.Get(`reqTest`))
	assert.Equal(t, req.Form.Get(`reqTest`), req3.Form.Get(`reqTest`))
	assert.Equal(t, `000`, req3.Context().Value(testContextKey).(string))
	assert.Equal(t, `*context.valueCtx`, fmt.Sprintf("%T", req3.Context()))
	assert.Equal(t, `*echo.xContext`, fmt.Sprintf("%T", MustGetContext(req3.Context())))
	assert.Equal(t, `0`, MustGetContext(req3.Context()).Internal().String(`_`))
	assert.True(t, IsMockContext(eCtx))
}

func TestFastHTTPContext(t *testing.T) {
	rCtx := &fasthttp.RequestCtx{}
	req := fasthttpng.NewRequest(rCtx)
	eCtx := echo.NewContext(req, fasthttpng.NewResponse(req), Default)
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
