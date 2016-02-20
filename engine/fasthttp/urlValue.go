// +build !appengine

package fasthttp

import (
	"net/url"

	"github.com/admpub/fasthttp"
)

type UrlValue struct {
	*fasthttp.Args
}

func (u *UrlValue) Add(key string, value string) {
	u.Args.Set(key, value)
}

func (u *UrlValue) Del(key string) {
	u.Args.Del(key)
}

func (u *UrlValue) Get(key string) string {
	return string(u.Args.Peek(key))
}

func (u *UrlValue) Set(key string, value string) {
	u.Args.Set(key, value)
}

func (u *UrlValue) Encode() string {
	return u.Args.String()
}

func (u *UrlValue) All() map[string][]string {
	r := make(map[string][]string)
	u.Args.VisitAll(func(k, v []byte) {
		key := string(k)
		if _, ok := r[key]; !ok {
			r[key] = make([]string, 0)
		}
		r[key] = append(r[key], string(v))
	})
	return r
}

func NewValue(c *fasthttp.RequestCtx) *Value {
	v := &Value{
		queryArgs: &UrlValue{Args: c.QueryArgs()},
		postArgs:  &UrlValue{Args: c.PostArgs()},
		context:   c,
	}
	return v
}

type Value struct {
	queryArgs *UrlValue
	postArgs  *UrlValue
	form      url.Values
	context   *fasthttp.RequestCtx
}

func (v *Value) Add(key string, value string) {
	if v.form == nil {
		return
	}
	v.form.Set(key, value)
}

func (v *Value) Del(key string) {
	if v.form == nil {
		return
	}
	v.form.Del(key)
}

func (v *Value) Get(key string) string {
	if v.form == nil {
		return ``
	}
	return v.form.Get(key)
}

func (v *Value) Set(key string, value string) {
	if v.form == nil {
		return
	}
	v.form.Set(key, value)
}

func (v *Value) Encode() string {
	if v.form == nil {
		return ``
	}
	return v.form.Encode()
}

func (v *Value) All() map[string][]string {
	if v.form == nil {
		v.form = url.Values(v.postArgs.All())
		for key, vals := range v.queryArgs.All() {
			v.form[key] = vals
		}
		/*
			mf, err := v.context.MultipartForm()
			if err == nil && mf.Value != nil {
				for key, vals := range mf.Value {
					v.form[key] = vals
				}
			}
		*/
	}
	return v.form
}
