// +build !appengine

package fasthttp

import "github.com/admpub/fasthttp"

type (
	URL struct {
		url *fasthttp.URI
	}
)

func (u *URL) Scheme() string {
	return string(u.url.Scheme())
}

func (u *URL) Host() string {
	return string(u.url.Host())
}

func (u *URL) SetPath(path string) {
	u.url.SetPath(path)
}

func (u *URL) Path() string {
	return string(u.url.Path())
}

func (u *URL) QueryValue(name string) string {
	return string(u.url.QueryArgs().Peek(name))
}

func (u *URL) RawQuery() string {
	return string(u.url.QueryString())
}
