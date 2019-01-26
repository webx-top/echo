package subdomains

import (
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
)

var Default = New()

func New() *Subdomains {
	s := &Subdomains{
		Hosts:   map[string]string{},
		Alias:   map[string]*Info{},
		Default: ``,
	}
	return s
}

type Info struct {
	Name string
	Host string
	*echo.Echo
}

type Subdomains struct {
	Hosts    map[string]string //{host:name}
	Alias    map[string]*Info
	Default  string //default name
	Protocol string //http/https
}

func (s *Subdomains) Add(name string, e *echo.Echo) *Subdomains {
	r := strings.SplitN(name, `@`, 2) //blog@www.blog.com
	var host string
	if len(r) > 1 {
		name = r[0]
		host = r[1]
	}
	s.Hosts[host] = name
	s.Alias[name] = &Info{Name: name, Host: host, Echo: e}
	return s
}

func (s *Subdomains) Get(args ...string) *Info {
	name := s.Default
	if len(args) > 0 {
		name = args[0]
	}
	if e, ok := s.Alias[name]; ok {
		return e
	}
	return nil
}

func (s *Subdomains) URL(purl string, args ...string) string {
	info := s.Get(args...)
	if info == nil {
		return purl
	}
	if len(info.Host) < 1 {
		if s.Default == info.Name {
			return purl
		}
		return `/` + info.Name + purl
	}
	if len(s.Protocol) < 1 {
		return `http://` + info.Host + purl
	}
	return s.Protocol + `://` + info.Host + purl
}

func (s *Subdomains) FindByDomain(host string) (*echo.Echo, bool) {
	name, exists := s.Hosts[host]
	if !exists {
		if p := strings.LastIndexByte(host, ':'); p > -1 {
			name, exists = s.Hosts[host[0:p]]
		}
		if !exists {
			name = s.Default
		}
	}
	var info *Info
	info, exists = s.Alias[name]
	if exists {
		return info.Echo, true
	}
	return nil, false
}

func (s *Subdomains) ServeHTTP(r engine.Request, w engine.Response) {
	domain := r.Host()
	handler, exists := s.FindByDomain(domain)
	if exists && handler != nil {
		handler.ServeHTTP(r, w)
	} else {
		w.NotFound()
	}
}

func (s *Subdomains) Run(args ...interface{}) {
	var eng engine.Engine
	var arg interface{}
	size := len(args)
	if size > 0 {
		arg = args[0]
	}
	if size > 1 {
		if conf, ok := arg.(*engine.Config); ok {
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.NewWithConfig(conf)
				} else {
					eng = standard.NewWithConfig(conf)
				}
			} else {
				eng = fasthttp.NewWithConfig(conf)
			}
		} else {
			addr := `:80`
			if v, ok := arg.(string); ok && len(v) > 0 {
				addr = v
			}
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.New(addr)
				} else {
					eng = standard.New(addr)
				}
			} else {
				eng = fasthttp.New(addr)
			}
		}
	} else {
		switch v := arg.(type) {
		case string:
			eng = fasthttp.New(v)
		case engine.Engine:
			eng = v
		default:
			eng = fasthttp.New(`:80`)
		}
	}
	e := s.Get()
	if e == nil {
		for _, info := range s.Alias {
			e = info
			break
		}
	}
	e.Logger().Info(`Server has been launched.`)
	e.Run(eng, s)
	e.Logger().Info(`Server has been closed.`)
}
