package echo

type (
	Group struct {
		prefix     string
		middleware []Middleware
		echo       *Echo
	}
)

func (g *Group) URL(h interface{}, params ...interface{}) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
}

func (g *Group) Any(path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Route(methods string, path string, h interface{}, middleware ...interface{}) {
	g.Match(httpMethodRegexp.Split(methods, -1), path, h, middleware...)
}

func (g *Group) Match(methods []string, path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Use(middleware ...interface{}) {
	for _, m := range middleware {
		g.middleware = append(g.middleware, WrapMiddleware(m))
	}
}

func (g *Group) PreUse(middleware ...interface{}) {
	middlewares := make([]Middleware, 0)
	for _, m := range middleware {
		middlewares = append(middlewares, WrapMiddleware(m))
	}
	g.middleware = append(middlewares, g.middleware...)
}

func (g *Group) Connect(path string, h interface{}, m ...interface{}) {
	g.add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h interface{}, m ...interface{}) {
	g.add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h interface{}, m ...interface{}) {
	g.add(GET, path, h, m...)
}

func (g *Group) Head(path string, h interface{}, m ...interface{}) {
	g.add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h interface{}, m ...interface{}) {
	g.add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h interface{}, m ...interface{}) {
	g.add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h interface{}, m ...interface{}) {
	g.add(POST, path, h, m...)
}

func (g *Group) Put(path string, h interface{}, m ...interface{}) {
	g.add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h interface{}, m ...interface{}) {
	g.add(TRACE, path, h, m...)
}

func (g *Group) Group(prefix string, m ...interface{}) *Group {
	return g.echo.Group(g.prefix+prefix, m...)
}

func (g *Group) add(method, path string, h interface{}, middleware ...interface{}) {
	var handler Handler = WrapHandler(h)
	if handler == nil {
		return
	}
	path = g.prefix + path

	var name string
	if hn, ok := handler.(HandleNamer); ok {
		name = hn.HandleName()
	} else {
		name = handlerName(handler)
	}

	for _, m := range g.middleware {
		handler = m.Handle(handler)
	}

	for _, m := range middleware {
		mw := WrapMiddleware(m)
		handler = mw.Handle(handler)
	}
	hdl := HandlerFunc(func(c Context) error {
		return handler.Handle(c)
	})
	fpath, pnames := g.echo.router.Add(method, path, hdl, g.echo)
	g.echo.logger.Debugf(`Route: %7v %-30v -> %v`, method, fpath, name)
	r := &Route{
		Method:      method,
		Path:        path,
		Handler:     hdl,
		HandlerName: name,
		Format:      fpath,
		Params:      pnames,
		Prefix:      g.prefix,
	}
	if _, ok := g.echo.router.nroute[name]; !ok {
		g.echo.router.nroute[name] = []int{len(g.echo.router.routes)}
	} else {
		g.echo.router.nroute[name] = append(g.echo.router.nroute[name], len(g.echo.router.routes))
	}
	g.echo.router.routes = append(g.echo.router.routes, r)
}
