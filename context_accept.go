package echo

import "github.com/webx-top/echo/param"

type Accept struct {
	Type   string
	Params param.StringMap
}

func NewAccept() *Accept {
	return &Accept{
		Params: param.StringMap{},
	}
}
