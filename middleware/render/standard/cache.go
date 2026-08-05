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

func (c *CacheData) hasBlock(blocks ...string) bool {
	for _, blockName := range blocks {
		if _, ok := c.blocks[blockName]; !ok {
			return false
		}
	}
	return true
}

func (c *CacheData) hasAnyBlock(blocks ...string) bool {
	for _, blockName := range blocks {
		if _, ok := c.blocks[blockName]; ok {
			return true
		}
	}
	return false
}

func (c *CacheData) setFunc(funcMap template.FuncMap) template.FuncMap {
	result := template.FuncMap{}
	for k, v := range funcMap {
		result[k] = v
	}
	result["hasBlock"] = c.hasBlock
	result["hasAnyBlock"] = c.hasAnyBlock
	return result
}
