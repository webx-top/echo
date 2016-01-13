package echo

import "net/http"

type (
	Router struct {
		tree   *node
		static map[string]*nodeLite
		routes []Route
		echo   *Echo
	}
	node struct {
		kind          kind
		label         byte
		prefix        string
		parent        *node
		children      children
		ppath         string
		pnames        []string
		methodHandler *methodHandler
		echo          *Echo
		isFileServer  bool
	}
	nodeLite struct {
		methodHandler *methodHandler
		echo          *Echo
		isFileServer  bool
	}
	kind          uint8
	children      []*node
	methodHandler struct {
		connect HandlerFunc
		delete  HandlerFunc
		get     HandlerFunc
		head    HandlerFunc
		options HandlerFunc
		patch   HandlerFunc
		post    HandlerFunc
		put     HandlerFunc
		trace   HandlerFunc
	}
)

const (
	skind kind = iota
	pkind
	mkind
)

func (m *methodHandler) addHandler(method string, h HandlerFunc) {
	switch method {
	case GET:
		m.get = h
	case POST:
		m.post = h
	case PUT:
		m.put = h
	case DELETE:
		m.delete = h
	case PATCH:
		m.patch = h
	case OPTIONS:
		m.options = h
	case HEAD:
		m.head = h
	case CONNECT:
		m.connect = h
	case TRACE:
		m.trace = h
	}
}

func (m *methodHandler) findHandler(method string) HandlerFunc {
	switch method {
	case GET:
		return m.get
	case POST:
		return m.post
	case PUT:
		return m.put
	case DELETE:
		return m.delete
	case PATCH:
		return m.patch
	case OPTIONS:
		return m.options
	case HEAD:
		return m.head
	case CONNECT:
		return m.connect
	case TRACE:
		return m.trace
	default:
		return nil
	}
}

func NewRouter(e *Echo) *Router {
	return &Router{
		tree: &node{
			methodHandler: new(methodHandler),
		},
		static: make(map[string]*nodeLite),
		routes: []Route{},
		echo:   e,
	}
}

func (r *Router) Add(method, path string, h HandlerFunc, e *Echo, args ...bool) {
	ppath := path        // Pristine path
	pnames := []string{} // Param names
	isFileServer := false
	if len(args) > 0 {
		isFileServer = args[0]
	}
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, skind, "", nil, e, isFileServer)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, pkind, ppath, pnames, e, isFileServer)
				return
			}
			r.insert(method, path[:i], nil, pkind, ppath, pnames, e, isFileServer)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, skind, "", nil, e, isFileServer)
			pnames = append(pnames, "_*")
			r.insert(method, path[:i+1], h, mkind, ppath, pnames, e, isFileServer)
			return
		}
	}

	//static route
	if m, ok := r.static[path]; ok {
		m.methodHandler.addHandler(method, h)
	} else {
		m = &nodeLite{
			methodHandler: &methodHandler{},
			echo:          e,
			isFileServer:  isFileServer,
		}
		m.methodHandler.addHandler(method, h)
		r.static[path] = m
	}
	//r.insert(method, path, h, skind, ppath, pnames, e,isFileServer)
}

func (r *Router) insert(method, path string, h HandlerFunc, t kind, ppath string, pnames []string, e *Echo, isFileServer bool) {
	// Adjust max param
	l := len(pnames)
	if *e.maxParam < l {
		*e.maxParam = l
	}

	cn := r.tree // Current node as root
	if cn == nil {
		panic("echo => invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
				cn.echo = e
			}
		} else if l < pl {
			// Split node
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.children, cn.methodHandler, cn.ppath, cn.pnames, cn.echo, isFileServer)

			// Reset parent node
			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.methodHandler = new(methodHandler)
			cn.ppath = ""
			cn.pnames = nil
			cn.echo = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.kind = t
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
				cn.echo = e
			} else {
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames, e, isFileServer)
				n.addHandler(method, h)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			n := newNode(t, search, cn, nil, new(methodHandler), ppath, pnames, e, isFileServer)
			n.addHandler(method, h)
			cn.addChild(n)
		} else {
			// Node already exists
			if h != nil {
				cn.addHandler(method, h)
				cn.ppath = path
				cn.pnames = pnames
				cn.echo = e
			}
		}
		return
	}
}

func newNode(t kind, pre string, p *node, c children, mh *methodHandler, ppath string, pnames []string, e *Echo, isFileServer bool) *node {
	return &node{
		kind:          t,
		label:         pre[0],
		prefix:        pre,
		parent:        p,
		children:      c,
		ppath:         ppath,
		pnames:        pnames,
		methodHandler: mh,
		echo:          e,
		isFileServer:  isFileServer,
	}
}

func (n *node) addChild(c *node) {
	n.children = append(n.children, c)
}

func (n *node) findChild(l byte, t kind) *node {
	for _, c := range n.children {
		if c.label == l && c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) addHandler(method string, h HandlerFunc) {
	n.methodHandler.addHandler(method, h)
}

func (n *node) findHandler(method string) HandlerFunc {
	return n.methodHandler.findHandler(method)
}

func (n *node) check405() HandlerFunc {
	for _, m := range methods {
		if h := n.findHandler(m); h != nil {
			return methodNotAllowedHandler
		}
	}
	return notFoundHandler
}

func (r *Router) Find(method, path string, ctx Context) (h HandlerFunc, e *Echo) {
	x := ctx.X()
	h = notFoundHandler
	e = r.echo

	if m, ok := r.static[path]; ok {
		h = m.methodHandler.findHandler(method)
		x.path = path
		x.pnames = []string{}
		x.isFileServer = m.isFileServer
		if m.echo != nil {
			e = m.echo
		}
		if h == nil {
			h = methodNotAllowedHandler
		}
		return
	}

	cn := r.tree // Current node as root

	var (
		search = path
		c      *node  // Child node
		n      int    // Param counter
		nk     kind   // Next kind
		nn     *node  // Next node
		ns     string // Next search
	)

	// Search order static > param > match-any
	for {
		if search == "" {
			goto End
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// Continue search
			search = search[l:]
		} else {
			cn = nn
			search = ns
			if nk == pkind {
				goto Param
			} else if nk == mkind {
				goto MatchAny
			} else {
				// Not found
				return
			}
		}

		if search == "" {
			goto End
		}

		// Static node
		c = cn.findChild(search[0], skind)
		if c != nil {
			// Save next
			if cn.label == '/' {
				nk = pkind
				nn = cn
				ns = search
			}
			cn = c
			continue
		}

		// Param node
	Param:
		c = cn.findChildByKind(pkind)
		if c != nil {
			// Save next
			if cn.label == '/' {
				nk = mkind
				nn = cn
				ns = search
			}
			cn = c
			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}
			x.pvalues[n] = search[:i]
			n++
			search = search[i:]
			continue
		}

		// Match-any node
	MatchAny:
		// c = cn.getChild()
		if cn = cn.findChildByKind(mkind); cn == nil {
			// Not found
			return
		}
		x.pvalues[len(cn.pnames)-1] = search
		goto End
	}

End:
	x.path = cn.ppath
	x.pnames = cn.pnames
	x.isFileServer = cn.isFileServer
	h = cn.findHandler(method)
	if cn.echo != nil {
		e = cn.echo
	}

	// NOTE: Slow zone...
	if h == nil {
		h = cn.check405()

		// Dig further for match-any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if cn = cn.findChildByKind(mkind); cn == nil {
			return
		}
		x.pvalues[len(cn.pnames)-1] = ""
		if h = cn.findHandler(method); h == nil {
			h = cn.check405()
		}
	}
	return
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := r.echo.pool.Get().(Context)
	h, _ := r.Find(req.Method, req.URL.Path, c)
	c.Reset(req, w, r.echo)
	if err := h(c); err != nil {
		r.echo.httpErrorHandler(err, c)
	}
	r.echo.pool.Put(c)
}
