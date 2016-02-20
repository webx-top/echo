// +build !appengine

package fasthttp

import (
	"github.com/admpub/fasthttp"
	"github.com/labstack/gommon/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Server struct {
		*fasthttp.Server
		config  *engine.Config
		handler engine.HandlerFunc
		logger  logger.Logger
	}
)

func New(addr string) *Server {
	c := &engine.Config{Address: addr}
	return NewConfig(c)
}

func NewTLS(addr, certfile, keyfile string) *Server {
	c := &engine.Config{
		Address:     addr,
		TLSCertfile: certfile,
		TLSKeyfile:  keyfile,
	}
	return NewConfig(c)
}

func NewConfig(c *engine.Config) (s *Server) {
	fastHTTPServer := &fasthttp.Server{
		ReadTimeout:        c.ReadTimeout,
		WriteTimeout:       c.WriteTimeout,
		MaxConnsPerIP:      c.MaxConnsPerIP,
		MaxRequestsPerConn: c.MaxRequestsPerConn,
		MaxRequestBodySize: c.MaxRequestBodySize,
	}
	s = &Server{
		Server: fastHTTPServer,
		config: c,
		handler: engine.ClearHandler(func(req engine.Request, res engine.Response) {
			s.logger.Info("handler not set")
		}),
		logger: log.New("echo"),
	}
	//s.Server.Logger = s.logger
	return
}

func (s *Server) SetHandler(h engine.HandlerFunc) {
	s.handler = engine.ClearHandler(h)
}

func (s *Server) SetLogger(l logger.Logger) {
	s.logger = l
}

func (s *Server) Start() {
	s.Server.Handler = func(c *fasthttp.RequestCtx) {
		req := NewRequest(c)
		res := NewResponse(c)
		s.handler(req, res)
	}
	s.logger.Fatal(s.Server.ListenAndServe(s.config.Address))
}
