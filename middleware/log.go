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

			remoteAddr := req.RemoteAddress()
			if ip := req.Header().Get(echo.HeaderXRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header().Get(echo.HeaderXForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}
			start := time.Now()
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}
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
