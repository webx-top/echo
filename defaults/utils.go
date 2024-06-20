package defaults

import (
	"context"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
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
	return ok
}
