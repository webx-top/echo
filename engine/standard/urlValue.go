package standard

import (
	"net/url"
)

type UrlValue struct {
	Args *url.Values
}

func (u *UrlValue) Add(key string, value string) {
	u.Args.Set(key, value)
}

func (u *UrlValue) Del(key string) {
	u.Args.Del(key)
}

func (u *UrlValue) Get(key string) string {
	return u.Args.Get(key)
}

func (u *UrlValue) Set(key string, value string) {
	u.Args.Set(key, value)
}

func (u *UrlValue) Encode() string {
	return u.Args.Encode()
}

func (u *UrlValue) All() map[string][]string {
	return *u.Args
}

func (u *UrlValue) Reset(data url.Values) {
	*u.Args = data
}
