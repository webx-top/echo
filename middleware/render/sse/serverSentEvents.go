package sse

import (
	"io"

	"github.com/admpub/sse"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
)

const MIMEEventStream = "text/event-stream"

func init() {
	render.Reg(`sse`, func(_ string) driver.Driver {
		return New()
	})
}

func New() *ServerSentEvents {
	return &ServerSentEvents{
		NopRenderer: &driver.NopRenderer{},
	}
}

type ServerSentEvents struct {
	*driver.NopRenderer
}

func (s *ServerSentEvents) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	hdr := c.Response().Header()
	hdr.Set(echo.HeaderContentType, MIMEEventStream)
	hdr.Set(`Cache-Control`, `no-cache`)
	hdr.Set(`Connection`, `keep-alive`)
	if v, y := data.(sse.Event); y {
		return sse.Encode(w, v)
	}
	return sse.Encode(w, sse.Event{
		Event: name,
		Data:  data,
	})
}
