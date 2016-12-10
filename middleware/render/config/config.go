package config

import "strings"

type Template struct {
	Theme        string
	Engine       string
	Style        string
	Reload       bool
	ParseStrings map[string]string
}

func (t *Template) Parser() func([]byte) []byte {
	if t.ParseStrings == nil {
		return nil
	}
	return func(b []byte) []byte {
		s := string(b)
		for oldVal, newVal := range t.ParseStrings {
			s = strings.Replace(s, oldVal, newVal, -1)
		}
		return []byte(s)
	}
}
