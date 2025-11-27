/*
Copyright 2016 Wenhui Shen <www.webx.top>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package language

import (
	"net/http"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Config struct {
	Project  string
	Default  string
	Fallback string
	AllList  []string

	// key: language, value: map[key]value.
	//
	// Example:
	// {
	// 	"en":{
	// 		"label":"English",
	// 		"flag":"🇺🇸"
	// 	},
	// 	"zh-CN":{
	// 		"label":"简体中文",
	// 		"flag":"🇨🇳"
	// 	},
	// }
	Extra map[string]param.Store

	RulesPath    []string
	MessagesPath []string

	// Reload indicates whether to reload the language file each time it is modified.
	Reload bool
	fsFunc func(string) http.FileSystem
	kvList echo.KVList
}

func (c *Config) Init() {
	c.kvList = make(echo.KVList, len(c.AllList))
	for index, lang := range c.AllList {
		extra := c.ExtraBy(lang)
		label := extra.String(`label`)
		if len(label) == 0 {
			label = lang
		}
		flag := extra.String(`flag`)
		var kv *echo.KV
		if len(flag) > 0 {
			kv = echo.NewKV(lang, flag+` `+label)
		} else {
			kv = echo.NewKV(lang, label)
		}
		kv.SetH(extra)
		c.kvList[index] = kv
	}
}

func (c Config) KVList() echo.KVList {
	return c.kvList
}

func (c *Config) SetFSFunc(fsFunc func(string) http.FileSystem) *Config {
	c.fsFunc = fsFunc
	return c
}

// FSFunc returns the configured http.FileSystem function.
func (c Config) FSFunc() func(string) http.FileSystem {
	return c.fsFunc
}

// ExtraBy returns the extra parameters for the specified language.
// If no extra parameters exist for the language or Extra is nil, returns an empty Store.
func (c Config) ExtraBy(lang string) param.Store {
	if c.Extra == nil {
		return param.Store{}
	}
	if extra, ok := c.Extra[lang]; ok {
		return extra
	}
	return param.Store{}
}

func (c Config) Clone() Config {
	cfg := Config{
		Project:      c.Project,
		Default:      c.Default,
		Fallback:     c.Fallback,
		AllList:      make([]string, len(c.AllList)),
		Extra:        make(map[string]param.Store, len(c.Extra)),
		RulesPath:    make([]string, len(c.RulesPath)),
		MessagesPath: make([]string, len(c.MessagesPath)),
		Reload:       c.Reload,
		fsFunc:       c.fsFunc,
		kvList:       make(echo.KVList, len(c.kvList)),
	}
	copy(cfg.AllList, c.AllList)
	copy(cfg.RulesPath, c.RulesPath)
	copy(cfg.MessagesPath, c.MessagesPath)
	for k, v := range c.Extra {
		cfg.Extra[k] = v.Clone()
	}
	for k, v := range c.kvList {
		cloned := v.Clone()
		cfg.kvList[k] = &cloned
	}
	return cfg
}
