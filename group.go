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

func (g *Group) Any(path string, h interface{}, middleware ...Middleware) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Match(methods []string, path string, h interface{}, middleware ...Middleware) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Use(m ...Middleware) {
	g.middleware = append(g.middleware, m...)
}

func (g *Group) PreUse(m ...Middleware) {
	g.middleware = append(m, g.middleware...)
}

func (g *Group) Connect(path string, h interface{}, m ...Middleware) {
	g.add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h interface{}, m ...Middleware) {
	g.add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h interface{}, m ...Middleware) {
	g.add(GET, path, h, m...)
}

func (g *Group) Head(path string, h interface{}, m ...Middleware) {
	g.add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h interface{}, m ...Middleware) {
	g.add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h interface{}, m ...Middleware) {
	g.add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h interface{}, m ...Middleware) {
	g.add(POST, path, h, m...)
}

func (g *Group) Put(path string, h interface{}, m ...Middleware) {
	g.add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h interface{}, m ...Middleware) {
	g.add(TRACE, path, h, m...)
}

func (g *Group) Group(prefix string, m ...Middleware) *Group {
	return g.echo.Group(g.prefix+prefix, m...)
}

func (g *Group) add(method, path string, h interface{}, middleware ...Middleware) {
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

	middleware = append(g.middleware, middleware...)
	for _, m := range middleware {
		handler = m.Handle(handler)
	}
	fpath, pnames := g.echo.router.Add(method, path, HandlerFunc(func(c Context) error {
		return handler.Handle(c)
	}), g.echo)
	g.echo.logger.Debugf(`ROUTE|[%v]%v -> %v`+"\n", method, fpath, name)
	r := Route{
		Method:  method,
		Path:    path,
		Handler: name,
		Format:  fpath,
		Params:  pnames,
	}
	if _, ok := g.echo.router.nroute[name]; !ok {
		g.echo.router.nroute[name] = []int{len(g.echo.router.routes)}
	} else {
		g.echo.router.nroute[name] = append(g.echo.router.nroute[name], len(g.echo.router.routes))
	}
	g.echo.router.routes = append(g.echo.router.routes, r)
}
