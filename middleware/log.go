package middleware

import (
	"fmt"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/log"

	"github.com/webx-top/echo"
)

type VisitorInfo struct {
	RealIP       string
	Time         time.Time
	Elapsed      time.Duration
	Scheme       string
	Host         string
	URI          string
	Method       string
	UserAgent    string
	Referer      string
	RequestSize  int64
	ResponseSize int64
	ResponseCode int
}

var (
	terminalColors = map[StatusColor]*color.Color{
		`green`:  color.New(color.FgHiGreen),
		`red`:    color.New(color.FgHiRed),
		`yellow`: color.New(color.FgHiYellow),
		`cyan`:   color.New(color.FgHiCyan),
	}
)

// StatusColor 状态色
type StatusColor string

func (s StatusColor) String() string {
	return string(s)
}

// Terminal 控制台样式
func (s StatusColor) Terminal() *color.Color {
	return terminalColors[s]
}

// HTTPStatusColor HTTP状态码相应颜色
func HTTPStatusColor(httpCode int) StatusColor {
	s := `green`
	switch {
	case httpCode >= 500:
		s = `red`
	case httpCode >= 400:
		s = `yellow`
	case httpCode >= 300:
		s = `cyan`
	}
	return StatusColor(s)
}

func Log(recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	var logging func(*VisitorInfo)
	if len(recv) > 0 {
		logging = recv[0]
	}
	if logging == nil {
		logger := log.GetLogger(`HTTP`)
		logging = func(v *VisitorInfo) {
			colorSprint := HTTPStatusColor(v.ResponseCode).Terminal().SprintFunc()
			logger.Info(" " + colorSprint(v.ResponseCode) + " " + v.RealIP + " " + v.Method + " " + v.Scheme + " " + v.Host + " " + v.URI + " " + v.Elapsed.String() + " " + fmt.Sprint(v.ResponseSize))
		}
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			info := &VisitorInfo{Time: time.Now()}
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}
			info.RealIP = req.RealIP()
			info.UserAgent = req.UserAgent()
			info.Referer = req.Referer()
			info.RequestSize = req.Size()
			info.Elapsed = time.Now().Sub(info.Time)
			info.Method = req.Method()
			info.Host = req.Host()
			info.Scheme = req.Scheme()
			info.URI = req.URI()
			info.ResponseSize = res.Size()
			info.ResponseCode = res.Status()
			logging(info)
			return nil
		})
	}
}
