package standard

import (
	"html/template"
)

func NewCache(t *template.Template) *CacheData {
	return &CacheData{
		template: t,
		blocks:   map[string]struct{}{},
	}
}

type CacheData struct {
	template *template.Template
	blocks   map[string]struct{}
}

func (c *CacheData) setFunc(funcMap template.FuncMap) template.FuncMap {
	if funcMap == nil {
		funcMap = template.FuncMap{}
	}
	funcMap["hasBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := c.blocks[blockName]; !ok {
				return false
			}
		}
		return true
	}
	funcMap["hasAnyBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := c.blocks[blockName]; ok {
				return true
			}
		}
		return false
	}
	return funcMap
}
