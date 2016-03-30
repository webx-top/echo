// +build !appengine

package fasthttp

import (
	"net/url"

	"github.com/admpub/fasthttp"
)

type (
	URL struct {
		url   *fasthttp.URI
		query url.Values
	}
)

func (u *URL) SetPath(path string) {
	u.url.SetPath(path)
}

func (u *URL) Path() string {
	return string(u.url.Path())
}

func (u *URL) QueryValue(name string) string {
	return string(u.url.QueryArgs().Peek(name))
}

func (u *URL) QueryValues(name string) []string {
	u.Query()
	if v, ok := u.query[name]; ok {
		return v
	}
	return []string{}
}

func (u *URL) Query() url.Values {
	if u.query == nil {
		u.query = url.Values{}
		u.url.QueryArgs().VisitAll(func(key []byte, value []byte) {
			u.query.Set(string(key), string(value))
		})
	}
	return u.query
}

func (u *URL) RawQuery() string {
	return string(u.url.QueryString())
}

func (u *URL) Object() interface{} {
	return u.url
}

func (u *URL) reset(url *fasthttp.URI) {
	u.url = url
}
