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
	"bytes"
	"regexp"
	"strings"
	"sync"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

var (
	LangVarName        = `lang`
	DefaultLang        = `zh-CN`
	headerAcceptRemove = regexp.MustCompile(`;q=[0-9.]+`)
)

func New(c ...*Config) *Language {
	lang := &Language{
		List:    make(map[string]bool),
		Index:   make([]string, 0),
		Default: DefaultLang,
		translatePool: sync.Pool{
			New: func() interface{} {
				return &Translate{_pool: true}
			},
		},
	}
	if len(c) > 0 {
		lang.Init(c[0])
	}
	return lang
}

type Language struct {
	List          map[string]bool //语种列表
	Index         []string        //索引
	Default       string          //默认语种
	I18n          *I18n
	translatePool sync.Pool
}

func (a *Language) Init(c *Config) {
	if len(c.Default) > 0 {
		c.Default = echo.NewLangCode(c.Default).Normalize()
	}
	if len(c.Fallback) > 0 {
		c.Fallback = echo.NewLangCode(c.Fallback).Normalize()
	}
	if c.AllList != nil {
		for index, lang := range c.AllList {
			lang = echo.NewLangCode(lang).Normalize()
			a.Set(lang, true, lang == c.Default, true)
			c.AllList[index] = lang
		}
	} else {
		a.Set(c.Default, true, true, true)
		if c.Default != `en` {
			a.Set(`en`, true)
		}
	}
	a.I18n = NewI18n(c)
	if c.Reload {
		a.I18n.Monitor()
	}
}

// Set 记录语言
//   - on: 使用启用此语言
//   - args[0]: setDefault 是否设置为默认语言
//   - args[1]: normalized 是否已经标准格式化 lang 值
func (a *Language) Set(lang string, on bool, args ...bool) *Language {
	var setDefault, normalized bool
	com.ExtractSlicex(args, &setDefault, &normalized)
	if a.List == nil {
		a.List = make(map[string]bool)
	}
	if !normalized {
		lang = echo.NewLangCode(lang).Normalize()
	}
	if _, ok := a.List[lang]; !ok {
		a.Index = append(a.Index, lang)
	}
	a.List[lang] = on
	if on && setDefault {
		a.Default = lang
	}
	return a
}

func (a *Language) DetectURI(c echo.Context) string {
	dispatchPath := c.DispatchPath()
	p := strings.TrimPrefix(dispatchPath, `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if len(lang) == 0 {
		return lang
	}
	on, ok := a.List[lang]
	if !ok {
		return ``
	}
	c.SetDispatchPath(strings.TrimPrefix(p, lang))
	if !on {
		return ``
	}
	return lang
}

func (a *Language) Valid(lang string) bool {
	if len(lang) == 0 {
		return false
	}
	if on, ok := a.List[lang]; ok {
		return on
	}
	return false
}

func ParseHeader(al string, n int) []string {
	if len(al) == 0 {
		return []string{}
	}
	al = headerAcceptRemove.ReplaceAllString(al, ``)
	return strings.SplitN(al, `,`, n)
}

func (a *Language) DetectHeader(r engine.Request) string {
	lg := ParseHeader(r.Header().Get(`Accept-Language`), 5)
	for _, lang := range lg {
		if a.Valid(lang) {
			return lang
		}
		lang = strings.SplitN(lang, `-`, 2)[0]
		if a.Valid(lang) {
			return lang
		}
	}
	return a.Default
}

func (a *Language) AcquireTranslator(langCode string) *Translate {
	tr := a.translatePool.Get().(*Translate)
	tr.Reset(langCode, a)
	return tr
}

func (a *Language) ReleaseTranslator(tr *Translate) {
	tr.Release()
}

func (a *Language) release(c echo.Context) {
	c.Translator().(*Translate).Release()
}

func (a *Language) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			lang := c.Query(LangVarName)
			var hasCookie bool
			if !a.Valid(lang) {
				lang = a.DetectURI(c)
				if !a.Valid(lang) {
					lang = c.GetCookie(LangVarName)
					if !a.Valid(lang) {
						lang = a.DetectHeader(c.Request())
					} else {
						hasCookie = true
					}
				}
			}
			if !hasCookie {
				c.SetCookie(LangVarName, lang)
			}
			tr := a.AcquireTranslator(lang)
			c.OnRelease(a.release)
			c.SetTranslator(tr)
			return h.Handle(c)
		})
	})
}

func (a *Language) Config() *Config {
	return a.I18n.config
}

func (a *Language) Handler(e echo.RouteRegister, i18nJSVarName string) {
	e.Get(`/i18n.json`, func(c echo.Context) error {
		t := a.I18n.Get(c.Lang().String())
		if t != nil {
			messages := t.Messages()
			return c.JSON(messages)
		}
		return c.JSONBlob([]byte(`{}`))
	})
	e.Get(`/i18n.js`, func(c echo.Context) error {
		t := a.I18n.Get(c.Lang().String())
		buf := bytes.NewBuffer(nil)
		if t != nil {
			messages := t.Messages()
			if messages != nil {
				if len(i18nJSVarName) > 0 {
					buf.WriteString(i18nJSVarName + `=`)
				}
				b, _ := com.JSONEncode(messages)
				if len(b) > 0 {
					buf.Write(b)
				} else {
					buf.WriteString(`{}`)
				}
			}
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextJavaScriptCharsetUTF8)
		return c.Blob(buf.Bytes())
	})
}
