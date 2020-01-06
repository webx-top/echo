package mock

import "net/url"

type Value struct {
	form url.Values
}

func (v *Value) Add(key string, value string) {
	v.form.Add(key, value)
}

func (v *Value) Del(key string) {
	v.form.Del(key)
}

func (v *Value) Get(key string) string {
	return v.form.Get(key)
}

func (v *Value) Gets(key string) []string {
	if v, ok := v.form[key]; ok {
		return v
	}
	return []string{}
}

func (v *Value) Set(key string, value string) {
	v.form.Set(key, value)
}

func (v *Value) Encode() string {
	return v.form.Encode()
}

func (v *Value) All() map[string][]string {
	return v.form
}

func (v *Value) Reset(data url.Values) {
	v.form = data
}

func (v *Value) Merge(data url.Values) {
	for key, values := range data {
		for index, value := range values {
			if index == 0 {
				v.form.Set(key, value)
			} else {
				v.form.Add(key, value)
			}
		}
	}
}
