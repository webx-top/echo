package mockcontext

import (
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/mock"
)

func Acquire() echo.Context {
	c := poolMockContext.Get().(echo.Context)
	c.SetAuto(true)
	c.Request().Form().Set(echo.HeaderAccept, echo.MIMETextPlain)
	return c
}

func Release(c echo.Context) {
	c.Reset(mock.NewRequest(), mock.NewResponse())
	poolMockContext.Put(c)
}

func AcquireRequest() engine.Request {
	return poolMockRequest.Get().(engine.Request)
}

func AcquireResponse() engine.Response {
	return poolMockRequest.Get().(engine.Response)
}

func ReleaseRequest(c engine.Request) {
	poolMockRequest.Put(c)
}

func ReleaseResponse(c engine.Response) {
	poolMockResponse.Put(c)
}

var (
	poolMockContext = sync.Pool{
		New: func() interface{} {
			return echo.NewContext(mock.NewRequest(), mock.NewResponse(), defaults.Default)
		},
	}

	poolMockRequest = sync.Pool{
		New: func() interface{} {
			return mock.NewRequest()
		},
	}

	poolMockResponse = sync.Pool{
		New: func() interface{} {
			return mock.NewResponse()
		},
	}
)
