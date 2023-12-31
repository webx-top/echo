package middleware

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/random"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string

		// RequestIDHandler defines a function which is executed for a request id.
		RequestIDHandler func(echo.Context, string)
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:   echo.DefaultSkipper,
		Generator: generator,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFuncd {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			rid := c.Header(echo.HeaderXRequestID)
			if len(rid) == 0 {
				rid = config.Generator()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, rid)
			if config.RequestIDHandler != nil {
				config.RequestIDHandler(c, rid)
			}

			return next.Handle(c)
		}
	}
}

func generator() string {
	return random.String(32)
}
