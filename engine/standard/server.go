package standard

import (
	"net/http"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Server struct {
		server  *http.Server
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
		server: &http.Server{
			ReadTimeout: c.ReadTimeout,
			Addr:        c.Address,
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
					return &Header{}
				},
			},
			responseHeader: sync.Pool{
				New: func() interface{} {
					return &Header{}
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
	s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request
		req := s.pool.request.Get().(*Request)
		reqHdr := s.pool.requestHeader.Get().(*Header)
		reqHdr.reset(r.Header)
		reqURL := s.pool.url.Get().(*URL)
		reqURL.reset(r.URL)
		req.reset(r, reqHdr, reqURL)
		req.config = s.config

		// Response
		res := s.pool.response.Get().(*Response)
		resHdr := s.pool.responseHeader.Get().(*Header)
		resHdr.reset(w.Header())
		res.reset(w, r, resHdr)
		res.config = s.config

		s.handler.ServeHTTP(req, res)

		s.pool.request.Put(req)
		s.pool.requestHeader.Put(reqHdr)
		s.pool.url.Put(reqURL)
		s.pool.response.Put(res)
		s.pool.responseHeader.Put(resHdr)
	})

	certfile := s.config.TLSCertfile
	keyfile := s.config.TLSKeyfile
	if certfile != "" && keyfile != "" {
		s.logger.Fatal(s.server.ListenAndServeTLS(certfile, keyfile))
	} else {
		s.logger.Fatal(s.server.ListenAndServe())
	}
}
