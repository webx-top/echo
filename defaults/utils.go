package defaults

import (
	"context"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/testing"
)

func MustGetContext(ctx context.Context, args ...*echo.Echo) echo.Context {
	eCtx, ok := ctx.(echo.Context)
	if !ok {
		eCtx, ok = echo.FromStdContext(ctx)
	}
	if !ok {
		eCtx = NewMockContext(args...)
		if ctx != nil {
			req := eCtx.Request().StdRequest()
			*req = *eCtx.WithContext(ctx)
		}
	}
	return eCtx
}

func NewMockContext(args ...*echo.Echo) echo.Context {
	var e *echo.Echo
	if len(args) > 0 {
		e = args[0]
	} else {
		e = Default
	}
	return echo.NewContext(mock.NewRequest(), mock.NewResponse(), e)
}

func IsMockContext(c echo.Context) bool {
	_, ok := c.Request().(*mock.Request)
	if !ok {
		ok = testing.IsMock(c)
	}
	return ok
}

var poolMockContextIniters = []func(ctx echo.Context){}

func RegisterPoolMockContextIniter(init func(echo.Context)) {
	poolMockContextIniters = append(poolMockContextIniters, init)
}

var poolMockContext = sync.Pool{
	New: func() interface{} {
		c := echo.NewContext(mock.NewRequest(), mock.NewResponse(), Default)
		for _, f := range poolMockContextIniters {
			f(c)
		}
		return c
	},
}

func AcquireMockContext() echo.Context {
	return poolMockContext.Get().(echo.Context)
}

func ReleaseMockContext(ctx echo.Context) {
	if v, y := ctx.(echo.ContextReseter); y {
		v.Reset(nil, nil)
	}
	poolMockContext.Put(ctx)
}
