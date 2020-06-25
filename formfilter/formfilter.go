package formfilter

import (
	"strings"

	"github.com/webx-top/echo"
)

type (
	Filter  func(*Data)
	Filters map[string][]Filter
	Options func() (key string, filter Filter)
	Data    struct {
		Key           string
		Value         []string
		normalizedKey string
	}
)

func (d *Data) NormalizedKey() string {
	return d.normalizedKey
}

func Build(options ...Options) echo.FormDataFilter {
	filters := Filters{}
	for _, opt := range options {
		key, filter := opt()
		key = strings.Title(key)
		if _, ok := filters[key]; !ok {
			filters[key] = []Filter{}
		}
		filters[key] = append(filters[key], filter)
	}
	return echo.FormDataFilter(func(k string, v []string) (string, []string) {
		key := strings.Title(k)
		filters, ok := filters[key]
		if !ok {
			return k, v
		}
		data := &Data{Key: k, Value: v, normalizedKey: key}
		for _, filter := range filters {
			filter(data)
			if len(data.Key) == 0 {
				break
			}
		}
		return data.Key, data.Value
	})
}
