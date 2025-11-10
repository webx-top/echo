package echo

type (
	Host struct {
		head   Handler
		group  *Group
		groups map[string]*Group
		Router *Router
	}
	TypeHost struct {
		host   string
		router *Router
		echo   *Echo
	}
)

func (t TypeHost) URI(handler interface{}, params ...interface{}) string {
	if t.router == nil || t.echo == nil {
		return ``
	}
	return t.host + t.echo.URI(handler, params...)
}

func (t TypeHost) URIWithContext(c Context, handler interface{}, params ...interface{}) string {
	if t.router == nil || t.echo == nil {
		return ``
	}
	return t.host + t.echo.URIWithContext(c, handler, params...)
}

func (t TypeHost) String() string {
	return t.host
}

func (h *Host) Host(args ...interface{}) (r TypeHost) {
	if h.group == nil || h.group.host == nil {
		return
	}
	r.echo = h.group.echo
	r.router = h.Router
	if len(args) != 1 {
		r.host = h.group.host.Format(args...)
		return
	}
	switch v := args[0].(type) {
	case map[string]interface{}:
		r.host = h.group.host.FormatMap(v)
	case H:
		r.host = h.group.host.FormatMap(v)
	default:
		r.host = h.group.host.Format(args...)
	}
	return
}
