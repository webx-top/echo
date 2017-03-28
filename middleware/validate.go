package middleware

import "github.com/webx-top/echo"

func Validate(validator echo.Validator, skipper ...echo.Skipper) echo.MiddlewareFunc {
	var skip echo.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = echo.DefaultSkipper
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skip(c) {
				return h.Handle(c)
			}
			c.SetValidator(validator)
			return h.Handle(c)
		})
	}
}
