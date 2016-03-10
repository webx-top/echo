package echo

//import "fmt"

type (
	Group struct {
		prefix     string
		middleware []Middleware
		echo       *Echo
	}
)

func (g *Group) URL(h Handler, params ...interface{}) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
}

func (g *Group) Any(path string, h Handler, middleware ...Middleware) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Match(methods []string, path string, h Handler, middleware ...Middleware) {
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

func (g *Group) Connect(path string, h Handler, m ...Middleware) {
	g.add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h Handler, m ...Middleware) {
	g.add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h Handler, m ...Middleware) {
	g.add(GET, path, h, m...)
}

func (g *Group) Head(path string, h Handler, m ...Middleware) {
	g.add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h Handler, m ...Middleware) {
	g.add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h Handler, m ...Middleware) {
	g.add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h Handler, m ...Middleware) {
	g.add(POST, path, h, m...)
}

func (g *Group) Put(path string, h Handler, m ...Middleware) {
	g.add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h Handler, m ...Middleware) {
	g.add(TRACE, path, h, m...)
}

func (g *Group) Group(prefix string, m ...Middleware) *Group {
	return g.echo.Group(g.prefix+prefix, m...)
}

func (g *Group) add(method, path string, handler Handler, middleware ...Middleware) {
	path = g.prefix + path
	name := handlerName(handler)
	middleware = append(g.middleware, middleware...)
	// for k, mw := range g.echo.middleware {
	// 	fmt.Printf("%v. %+v\n", k, handlerName(mw))
	// }
	// fmt.Printf("=========================%+v\n", handlerName(h))
	for _, m := range middleware {
		handler = m.Handle(handler)
	}
	g.echo.router.Add(method, path, HandlerFunc(func(c Context) error {
		return handler.Handle(c)
	}), g.echo)
	r := Route{
		Method:  method,
		Path:    path,
		Handler: name,
	}
	g.echo.router.routes = append(g.echo.router.routes, r)
}
