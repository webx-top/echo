package middleware

import (
	"github.com/webx-top/echo"
)

var DefaultAJAXConfig = AJAXConfig{
	Skipper:   echo.DefaultSkipper,
	Handler:   echo.HandlerFuncs{},
	ParamName: `op`,
	OnlyAJAX:  true,
}

type AJAXConfig struct {
	Skipper   echo.Skipper
	Handler   echo.HandlerFuncs
	ParamName string
	OnlyAJAX  bool
}

func AJAX(handler echo.HandlerFuncs) echo.MiddlewareFunc {
	config := DefaultAJAXConfig
	config.Handler = handler
	return AJAXWithConfig(config)
}

func AJAXWithConfig(config AJAXConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultAJAXConfig.Skipper
	}
	if config.Handler == nil {
		config.Handler = DefaultAJAXConfig.Handler
	}
	if len(config.ParamName) == 0 {
		config.ParamName = DefaultAJAXConfig.ParamName
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) || (config.OnlyAJAX && !c.IsAjax()) {
				return h.Handle(c)
			}
			operate := c.Form(config.ParamName)
			if len(operate) > 0 {
				return config.Handler.Call(c, operate)
			}
			return h.Handle(c)
		})
	}
}
