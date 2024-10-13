package echo

import "net/http"

func (c *XContext) Session() Sessioner {
	return c.sessioner
}

func (c *XContext) Flash(names ...string) (r interface{}) {
	if v := c.sessioner.Flashes(names...); len(v) > 0 {
		r = v[len(v)-1]
	}
	return r
}

func (c *XContext) SetCookieOptions(opts *CookieOptions) {
	c.SessionOptions().CookieOptions = opts
}

func (c *XContext) CookieOptions() *CookieOptions {
	return c.SessionOptions().CookieOptions
}

func (c *XContext) SetSessionOptions(opts *SessionOptions) {
	c.sessionOptions = opts
}

func (c *XContext) SessionOptions() *SessionOptions {
	if c.sessionOptions == nil {
		c.sessionOptions = DefaultSessionOptions
	}
	return c.sessionOptions
}

func (c *XContext) NewCookie(key string, value string) *http.Cookie {
	return NewCookie(key, value, c.CookieOptions())
}

func (c *XContext) Cookie() Cookier {
	return c.cookier
}

func (c *XContext) GetCookie(key string) string {
	return c.cookier.Get(key)
}

func (c *XContext) SetCookie(key string, val string, args ...interface{}) {
	c.cookier.Set(key, val, args...)
}
