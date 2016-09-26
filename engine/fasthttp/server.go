// +build !appengine

package fasthttp

import (
	"net"
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

func New(addr string) *Server {
	c := &engine.Config{Address: addr}
	return NewWithConfig(c)
}

func NewWithTLS(addr, certFile, keyFile string) *Server {
	c := &engine.Config{
		Address:     addr,
		TLSCertFile: certFile,
		TLSKeyFile:  keyFile,
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
		handler: engine.ClearHandler(engine.HandlerFunc(func(req engine.Request, res engine.Response) {
			s.logger.Error("handler not set, use `SetHandler()` to set it.")
		})),
		logger: log.GetLogger("echo"),
	}
	s.Handler = s.ServeHTTP
	return
}

func (s *Server) SetHandler(h engine.Handler) {
	s.handler = engine.ClearHandler(h)
}

func (s *Server) SetLogger(l logger.Logger) {
	s.logger = l
}

// Start implements `engine.Server#Start` function.
func (s *Server) Start() error {
	if s.config.Listener == nil {
		return s.startDefaultListener()
	}
	return s.startCustomListener()

}

// Stop implements `engine.Server#Stop` function.
func (s *Server) Stop() error {
	if s.config.Listener == nil {
		return nil
	}
	return s.config.Listener.Close()
}

func (s *Server) startDefaultListener() error {
	/*
		c := s.config
		if len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0 {
			s.logger.Info(`FastHTTP is running at `, c.Address, ` [TLS]`)
			return s.ListenAndServeTLS(c.Address, c.TLSCertFile, c.TLSKeyFile)
		}
		s.logger.Info(`FastHTTP is running at `, c.Address)
		return s.ListenAndServe(c.Address)
	*/
	ln, err := net.Listen("tcp4", s.config.Address)
	if err != nil {
		return err
	}
	s.config.Listener = ln
	return s.startCustomListener()
}

func (s *Server) startCustomListener() error {
	c := s.config
	if len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0 {
		s.logger.Info(`FastHTTP is running at `, c.Listener.Addr(), ` [TLS]`)
		return s.ServeTLS(c.Listener, c.TLSCertFile, c.TLSKeyFile)
	}
	s.logger.Info(`FastHTTP is running at `, c.Listener.Addr())
	return s.Serve(c.Listener)
}

func (s *Server) ServeHTTP(c *fasthttp.RequestCtx) {
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
