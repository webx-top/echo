package opentracing

import (
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/webx-top/echo"
)

const defaultComponentName = "webx-top/echo"

type (
	// TraceConfig defines the config for Trace middleware.
	TraceConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper

		// OpenTracing Tracer instance which should be got before
		tracer opentracing.Tracer

		// componentName used for describing the tracing component name
		componentName string
	}
)

var (
	// DefaultTraceConfig is the default Trace middleware config.
	DefaultTraceConfig = TraceConfig{
		Skipper:       echo.DefaultSkipper,
		componentName: defaultComponentName,
	}
	ErrSpanContextInject = errors.New("SpanContext Inject Error!")
)

// Trace returns a Trace middleware.
//
// Trace middleware traces http requests and reporting errors.
func Trace(tracer opentracing.Tracer) echo.MiddlewareFunc {
	c := DefaultTraceConfig
	c.tracer = tracer
	c.componentName = defaultComponentName
	return TraceWithConfig(c)
}

// TraceWithConfig returns a Trace middleware with config.
// See: `Trace()`.
func TraceWithConfig(config TraceConfig) echo.MiddlewareFunc {
	if config.tracer == nil {
		panic("echo: trace middleware requires opentracing tracer")
	}
	if config.Skipper == nil {
		config.Skipper = echo.DefaultSkipper
	}
	if len(config.componentName) == 0 {
		config.componentName = defaultComponentName
	}

	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}

			req := c.Request()
			opname := "HTTP " + req.Method()
			var sp opentracing.Span
			tr := config.tracer
			carrier := opentracing.HTTPHeadersCarrier(req.Header().Std())
			if ctx, err := tr.Extract(opentracing.HTTPHeaders, carrier); err != nil {
				sp = tr.StartSpan(opname)
			} else {
				sp = tr.StartSpan(opname, ext.RPCServerOption(ctx))
			}

			ext.HTTPMethod.Set(sp, req.Method())
			ext.HTTPUrl.Set(sp, req.URL().String())
			ext.Component.Set(sp, config.componentName)

			*req.StdRequest() = *c.WithContext(opentracing.ContextWithSpan(c.StdContext(), sp))

			err := tr.Inject(sp.Context(), opentracing.HTTPHeaders, carrier)
			if err != nil {
				return ErrSpanContextInject
			}

			defer func() {
				status := c.Response().Status()
				committed := c.Response().Committed()
				ext.HTTPStatusCode.Set(sp, uint16(status))
				if status >= http.StatusInternalServerError || !committed {
					ext.Error.Set(sp, true)
				}
				sp.Finish()
			}()

			return h.Handle(c)
		})
	}
}
