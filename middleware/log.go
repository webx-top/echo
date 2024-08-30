package middleware

import (
	"io"
	std "log"
	"strconv"
	"sync"
	"time"

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

var emptyTime = time.Time{}

func (v *VisitorInfo) reset() {
	v.RealIP = ``
	v.Time = emptyTime
	v.Elapsed = 0
	v.Scheme = ``
	v.Host = ``
	v.URI = ``
	v.Method = ``
	v.UserAgent = ``
	v.Referer = ``
	v.RequestSize = 0
	v.ResponseSize = 0
	v.ResponseCode = 0
}

func (v *VisitorInfo) SetFromContext(c echo.Context) {
	req := c.Request()
	res := c.Response()
	v.RealIP = c.RealIP()
	v.UserAgent = req.UserAgent()
	v.Referer = req.Referer()
	v.RequestSize = req.Size()
	v.Elapsed = time.Since(v.Time)
	v.Method = req.Method()
	v.Host = req.Host()
	v.Scheme = req.Scheme()
	v.URI = req.URI()
	v.ResponseSize = res.Size()
	v.ResponseCode = res.Status()
}

var DefaultLogWriter = GetDefaultLogWriter()
var visitorInfoPool = sync.Pool{
	New: func() interface{} {
		return &VisitorInfo{}
	},
}

func Log(recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	return LogWithWriter(nil, recv...)
}

func AcquireVisitorInfo() *VisitorInfo {
	return visitorInfoPool.Get().(*VisitorInfo)
}

func ReleaseVisitorInfo(v *VisitorInfo) {
	v.reset()
	visitorInfoPool.Put(v)
}

func GetDefaultLogWriter() io.Writer {
	return std.Writer()
}

func LogWithWriter(writer io.Writer, recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	var logging func(*VisitorInfo)
	if len(recv) > 0 {
		logging = recv[0]
	}
	return LogWithConfig(LogConfig{
		Writer:  writer,
		Execute: logging,
	})
}

type LogConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper echo.Skipper       `json:"-"`
	Writer  io.Writer          `json:"-"`
	Execute func(*VisitorInfo) `json:"-"`
}

func LogWithConfig(config LogConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = echo.DefaultSkipper
	}
	if config.Writer == nil {
		config.Writer = DefaultLogWriter
	}
	if config.Execute == nil {
		logger := std.New(config.Writer, ``, 0)
		config.Execute = func(v *VisitorInfo) {
			logger.Println(":" + strconv.Itoa(v.ResponseCode) + ": " + v.Time.Format(time.RFC3339) + " " + v.RealIP + " " + v.Method + " " + v.Scheme + " " + v.Host + " " + v.URI + " " + v.Elapsed.String() + " " + strconv.FormatInt(v.ResponseSize, 10))
		}
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}
			info := AcquireVisitorInfo()
			info.Time = time.Now()
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}
			info.SetFromContext(c)
			config.Execute(info)
			ReleaseVisitorInfo(info)
			return nil
		})
	}
}
