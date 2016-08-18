package middleware

import (
	"net"
	"time"

	"github.com/webx-top/echo"
)

func Log() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger := c.Logger()

			start := time.Now()
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}

			remoteAddr := req.RealIP()
			stop := time.Now()
			method := req.Method()
			uri := req.URI()
			size := res.Size()
			code := res.Status()
			logger.Infof("%s %s %s %v %s %d", remoteAddr, method, uri, code, stop.Sub(start), size)
			return nil
		})
	}
}
