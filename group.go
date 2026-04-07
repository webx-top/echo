package echo

type Group struct {
	parent     *Group
	host       *host
	prefix     string
	middleware []Middleware
	echo       *Echo
	meta       H
}

func (g *Group) URL(h any, params ...any) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetAlias(alias string) *Group {
	if g.host != nil {
		g.host.alias = alias
		for a, v := range g.echo.hostAlias {
			if v == g.host.name {
				delete(g.echo.hostAlias, a)
			}
		}
		if len(alias) > 0 {
			g.echo.hostAlias[alias] = g.host.name
		}
	}
	return g
}

func (g *Group) Alias(alias string) Hoster {
	if name, ok := g.echo.hostAlias[alias]; ok {
		hs, ok := g.echo.hosts[name]
		if !ok || hs == nil || hs.group == nil {
			return nil
		}
		return hs.group.host
	}
	return nil
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
}

func (g *Group) Use(middleware ...any) {
	for _, m := range middleware {
		g.middleware = append(g.middleware, g.echo.WrapMiddleware(m))
		if g.echo.MiddlewareDebug {
			g.echo.logger.Debugf(`Middleware[Use](%p): [%s] -> %s`, m, g.prefix, HandlerName(m))
		}
	}
}

// Pre adds handler to the middleware chain.
func (g *Group) Pre(middleware ...any) {
	var middlewares []Middleware
	for _, m := range middleware {
		middlewares = append(middlewares, g.echo.WrapMiddleware(m))
		if g.echo.MiddlewareDebug {
			g.echo.logger.Debugf(`Middleware[Pre](%p): [%s] -> %s`, m, g.prefix, HandlerName(m))
		}
	}
	g.middleware = append(middlewares, g.middleware...)
}

func (g *Group) Connect(path string, h any, m ...any) IRouter {
	return g.Add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h any, m ...any) IRouter {
	return g.Add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h any, m ...any) IRouter {
	return g.Add(GET, path, h, m...)
}

func (g *Group) Head(path string, h any, m ...any) IRouter {
	return g.Add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h any, m ...any) IRouter {
	return g.Add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h any, m ...any) IRouter {
	return g.Add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h any, m ...any) IRouter {
	return g.Add(POST, path, h, m...)
}

func (g *Group) Put(path string, h any, m ...any) IRouter {
	return g.Add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h any, m ...any) IRouter {
	return g.Add(TRACE, path, h, m...)
}

func (g *Group) Any(path string, h any, middleware ...any) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, g.Add(m, path, h, middleware...))
	}
	return routes
}

func (g *Group) Route(methods string, path string, h any, middleware ...any) IRouter {
	return g.Match(splitHTTPMethod.Split(methods, -1), path, h, middleware...)
}

func (g *Group) Match(methods []string, path string, h any, middleware ...any) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, g.Add(m, path, h, middleware...))
	}
	return routes
}

func (g *Group) getMiddlewares() []Middleware {
	middlewares := []Middleware{}
	if g.parent != nil {
		middlewares = append(middlewares, g.parent.getMiddlewares()...)
		middlewares = append(middlewares, g.middleware...)
	} else {
		middlewares = g.middleware
	}
	return middlewares
}

func (g *Group) Group(prefix string, middleware ...any) *Group {
	if g.host != nil {
		subG, y := g.echo.hosts[g.host.name].groups[prefix]
		if !y {
			subG = &Group{parent: g, host: g.host, prefix: prefix, echo: g.echo, meta: H{}}
			g.echo.hosts[g.host.name].groups[prefix] = subG
			if len(g.meta) > 0 {
				subG.meta.DeepMerge(g.meta)
			}
		}
		if len(middleware) > 0 {
			subG.Use(middleware...)
		}
		return subG
	}
	return g.echo.subgroup(g, prefix, middleware...)
}

// Static implements `Echo#Static()` for sub-routes within the Group.
func (g *Group) Static(prefix, root string) {
	static(g, prefix, root)
}

// File implements `Echo#File()` for sub-routes within the Group.
func (g *Group) File(path, file string) {
	g.echo.File(g.prefix+path, file)
}

func (g *Group) Prefix() string {
	return g.prefix
}

func (g *Group) Echo() *Echo {
	return g.echo
}

// MetaHandler Add meta information about endpoint
func (g *Group) MetaHandler(m H, handler any, requests ...any) Handler {
	return g.echo.MetaHandler(m, handler, requests...)
}

func (g *Group) Add(method, path string, h any, middleware ...any) *Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	var host string
	if g.host != nil {
		host = g.host.name
	}
	r := g.echo.addWithGroup(g, host, method, g.prefix, g.prefix+path, h, middleware...)
	if len(g.meta) > 0 {
		r.Meta = H{}
		r.Meta.DeepMerge(g.meta)
	}
	return r
}

func (g *Group) SetMeta(meta H) *Group {
	g.meta = meta
	return g
}

func (g *Group) SetMetaKV(key string, value any) *Group {
	if g.meta == nil {
		g.meta = H{}
	}
	g.meta[key] = value
	return g
}
