package mock

import (
	"net/http"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

func NewRequest(reqs ...*http.Request) engine.Request {
	var req *http.Request
	if len(reqs) > 0 {
		req = reqs[0]
	}
	if req == nil {
		req = &http.Request{}
	}
	return &Request{
		Request: standard.NewRequest(req),
	}
}

type Request struct {
	*standard.Request
}
