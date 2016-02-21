package standard

import (
	"net/url"
)

type (
	URL struct {
		url   *url.URL
		query url.Values
	}
)

func (u *URL) URL() *url.URL {
	return u.url
}

func (u *URL) Scheme() string {
	return u.url.Scheme
}

func (u *URL) Host() string {
	return u.url.Host
}

func (u *URL) SetPath(path string) {
	u.url.Path = path
}

func (u *URL) Path() string {
	return u.url.Path
}

func (u *URL) QueryValue(name string) string {
	if u.query == nil {
		u.query = u.url.Query()
	}
	return u.query.Get(name)
}

func (u *URL) reset(url *url.URL) {
	u.url = url
	u.query = nil
}

func (u *URL) RawQuery() string {
	return u.url.RawQuery
}

func (u *URL) Object() interface{} {
	return u.url
}
