// +build !appengine

package fasthttp

import (
	"sync"

	"github.com/admpub/fasthttp"
	"github.com/labstack/gommon/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Server struct {
		server  *fasthttp.Server
		config  *engine.Config
		handler engine.Handler
		logger  logger.Logger
		pool    *Pool
	}

	Pool struct {
		request        sync.Pool
		response       sync.Pool
		requestHeader  sync.Pool
		responseHeader sync.Pool
		url            sync.Pool
	}
)

func New(addr string) *Server {
	c := &engine.Config{Address: addr}
	return NewWithConfig(c)
}

func NewWithTLS(addr, certfile, keyfile string) *Server {
	c := &engine.Config{
		Address:     addr,
		TLSCertfile: certfile,
		TLSKeyfile:  keyfile,
	}
	return NewWithConfig(c)
}

func NewWithConfig(c *engine.Config) (s *Server) {
	s = &Server{
		server: &fasthttp.Server{
			ReadTimeout:        c.ReadTimeout,
			WriteTimeout:       c.WriteTimeout,
			MaxConnsPerIP:      c.MaxConnsPerIP,
			MaxRequestsPerConn: c.MaxRequestsPerConn,
			MaxRequestBodySize: c.MaxRequestBodySize,
		},
		config: c,
		pool: &Pool{
			request: sync.Pool{
				New: func() interface{} {
					return &Request{}
				},
			},
			response: sync.Pool{
				New: func() interface{} {
					return &Response{logger: s.logger}
				},
			},
			requestHeader: sync.Pool{
				New: func() interface{} {
					return &RequestHeader{}
				},
			},
			responseHeader: sync.Pool{
				New: func() interface{} {
					return &ResponseHeader{}
				},
			},
			url: sync.Pool{
				New: func() interface{} {
					return &URL{}
				},
			},
		},
		handler: engine.ClearHandler(engine.HandlerFunc(func(req engine.Request, res engine.Response) {
			s.logger.Warn("handler not set")
		})),
		logger: log.New("echo"),
	}
	return
}

func (s *Server) SetHandler(h engine.Handler) {
	s.handler = engine.ClearHandler(h)
}

func (s *Server) SetLogger(l logger.Logger) {
	s.logger = l
}

func (s *Server) Start() {
	//s.server.Logger = s.logger
	s.server.Handler = func(c *fasthttp.RequestCtx) {
		// Request
		req := s.pool.request.Get().(*Request)
		reqHdr := s.pool.requestHeader.Get().(*RequestHeader)
		reqURL := s.pool.url.Get().(*URL)
		reqHdr.reset(&c.Request.Header)
		reqURL.reset(c.URI())
		req.reset(c, reqHdr, reqURL)

		// Response
		res := s.pool.response.Get().(*Response)
		resHdr := s.pool.responseHeader.Get().(*ResponseHeader)
		resHdr.reset(&c.Response.Header)
		res.reset(c, resHdr)

		s.handler.ServeHTTP(req, res)

		s.pool.request.Put(req)
		s.pool.requestHeader.Put(reqHdr)
		s.pool.url.Put(reqURL)
		s.pool.response.Put(res)
		s.pool.responseHeader.Put(resHdr)
	}
	addr := s.config.Address
	certfile := s.config.TLSCertfile
	keyfile := s.config.TLSKeyfile
	if certfile != "" && keyfile != "" {
		s.logger.Fatal(s.server.ListenAndServeTLS(addr, certfile, keyfile))
	} else {
		s.logger.Fatal(s.server.ListenAndServe(addr))
	}
}
