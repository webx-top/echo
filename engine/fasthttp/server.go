//go:build !appengine
// +build !appengine

package fasthttp

import (
	"context"
	"sync"

	"github.com/admpub/fasthttp"
	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Server struct {
		*fasthttp.Server
		config  *engine.Config
		handler engine.Handler
		logger  logger.Logger
		pool    *pool
	}

	pool struct {
		request        sync.Pool
		response       sync.Pool
		requestHeader  sync.Pool
		responseHeader sync.Pool
		url            sync.Pool
	}
)

func New(addr string, opts ...engine.ConfigSetter) *Server {
	c := &engine.Config{Address: addr}
	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

func NewWithTLS(addr, certFile, keyFile string, opts ...engine.ConfigSetter) *Server {
	c := &engine.Config{
		Address:     addr,
		TLSCertFile: certFile,
		TLSKeyFile:  keyFile,
	}
	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

func NewWithConfig(c *engine.Config) (s *Server) {
	s = &Server{
		Server: &fasthttp.Server{
			ReadTimeout:        c.ReadTimeout,
			WriteTimeout:       c.WriteTimeout,
			MaxConnsPerIP:      c.MaxConnsPerIP,
			MaxRequestsPerConn: c.MaxRequestsPerConn,
			MaxRequestBodySize: c.MaxRequestBodySize,
		},
		config: c,
		pool: &pool{
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
		handler: engine.HandlerFunc(func(req engine.Request, res engine.Response) {
			s.logger.Error("handler not set, use `SetHandler()` to set it.")
		}),
		logger: log.GetLogger("echo"),
	}
	s.Handler = s.ServeHTTP
	return
}

func (s *Server) SetHandler(h engine.Handler) {
	s.handler = h
}

func (s *Server) SetLogger(l logger.Logger) {
	s.logger = l
}

// Start implements `engine.Server#Start` function.
func (s *Server) Start() error {
	if s.config.Listener == nil {
		s.config.DisableHTTP2 = true
		err := s.config.InitListener(func() error {
			if s.config.TLSConfig == nil {
				return nil
			}
			s.config.TLSConfig.PreferServerCipherSuites = true
			return nil
		})
		if err != nil {
			return err
		}
	}
	s.config.Print(`fast`)
	return s.Serve(s.config.Listener)

}

// Stop implements `engine.Server#Stop` function.
func (s *Server) Stop() error {
	if s.config.Listener == nil {
		return nil
	}
	return s.config.Listener.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown()
}

func (s *Server) Config() *engine.Config {
	return s.config
}

func (s *Server) ServeHTTP(c *fasthttp.RequestCtx) {
	// Request
	req := s.pool.request.Get().(*Request)
	reqHdr := s.pool.requestHeader.Get().(*RequestHeader)
	reqURL := s.pool.url.Get().(*URL)
	reqHdr.reset(&c.Request.Header)
	reqURL.reset(c.URI())
	req.reset(c, reqHdr, reqURL)
	req.config = s.config

	// Response
	res := s.pool.response.Get().(*Response)
	resHdr := s.pool.responseHeader.Get().(*ResponseHeader)
	resHdr.reset(&c.Response.Header)
	res.reset(req, resHdr)

	s.handler.ServeHTTP(req, res)

	s.pool.request.Put(req)
	s.pool.requestHeader.Put(reqHdr)
	s.pool.url.Put(reqURL)
	s.pool.response.Put(res)
	s.pool.responseHeader.Put(resHdr)
}
