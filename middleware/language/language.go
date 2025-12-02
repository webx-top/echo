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

// New creates a new Language instance with optional configuration.
// If configuration is provided, it initializes the Language with the first config.
// The returned Language contains a sync.Pool for Translate instances and default settings.
func New(c ...*Config) *Language {
	lang := &Language{
		Default: DefaultLang,
		translatePool: sync.Pool{
			New: func() interface{} {
				return &Translate{_pool: true}
			},
		},
		langsMapGetter: LangsMapGetter,
	}
	if len(c) > 0 {
		lang.Init(c[0])
	}
	return lang
}

type Language struct {
	List           map[string]bool
	Default        string //默认语种
	I18n           *I18n
	translatePool  sync.Pool
	langsMapGetter func(echo.Context, *Language) (map[string]bool, error) //语种列表
}

var LangsMapGetter = func(c echo.Context, a *Language) (map[string]bool, error) {
	return a.List, nil
}

// Close closes the I18n instance if it exists.
func (a *Language) Close() {
	if a.I18n != nil {
		a.I18n.Close()
	}
}

// Init initializes the Language instance with the provided configuration.
// It sets up the language list based on the configuration, including the default language.
// If AllList is provided in the config, it registers all listed languages.
// Otherwise, it registers the default language and optionally English ('en') if not the default.
// It also initializes the I18n instance and starts monitoring for changes if Reload is enabled.
func (a *Language) Init(c *Config) {
	c.Init()
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
	a.List[lang] = on
	if on && setDefault {
		a.Default = lang
	}
	return a
}

// DetectURI detects the language code from the URI path.
// It checks if the detected language is valid according to the provided list.
// Returns the detected language code if valid, otherwise returns an empty string.
// The function also updates the dispatch path by removing the language prefix.
func (a *Language) DetectURI(c echo.Context, list map[string]bool) string {
	dispatchPath := c.DispatchPath()
	p := strings.TrimPrefix(dispatchPath, `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if !a.Valid(lang, list) {
		return ``
	}
	c.SetDispatchPath(strings.TrimPrefix(p, lang))
	return lang
}

// Valid checks if the specified language is valid and enabled in the provided list.
// Returns true if the language is non-empty and marked as enabled in the list, false otherwise.
func (a *Language) Valid(lang string, list map[string]bool) bool {
	if len(lang) == 0 {
		return false
	}
	if list == nil {
		list = a.List
	}
	if on, ok := list[lang]; ok {
		return on
	}
	return false
}

// ParseHeader parses the Accept-Language header string into a slice of language tags.
// It removes any quality values (e.g. ";q=0.8") and splits the string by commas.
// The n parameter controls the maximum number of languages to return (similar to strings.SplitN).
// Returns an empty slice if the input string is empty.
func ParseHeader(al string, n int) []string {
	if len(al) == 0 {
		return []string{}
	}
	al = headerAcceptRemove.ReplaceAllString(al, ``)
	return strings.SplitN(al, `,`, n)
}

// DetectHeader detects the preferred language from the Accept-Language header.
// It checks each language in the header against the provided list of valid languages,
// and returns the first valid match (including base language without region code).
// If no valid language is found, it returns the default language.
// Parameters:
//   - r: the request containing the Accept-Language header
//   - list: map of valid languages (keys are language codes, values are ignored)
//
// Returns:
//   - string: the detected language code or the default language
func (a *Language) DetectHeader(r engine.Request, list map[string]bool) string {
	lg := ParseHeader(r.Header().Get(`Accept-Language`), 5)
	for _, lang := range lg {
		if a.Valid(lang, list) {
			return lang
		}
		lang = strings.SplitN(lang, `-`, 2)[0]
		if a.Valid(lang, list) {
			return lang
		}
	}
	return a.Default
}

// AcquireTranslator gets a Translate instance from the pool and initializes it with the specified language code.
// The returned translator is ready to use for translation operations.
func (a *Language) AcquireTranslator(langCode string, list map[string]bool) *Translate {
	tr := a.translatePool.Get().(*Translate)
	tr.Reset(langCode, a, list)
	return tr
}

// ReleaseTranslator releases the resources associated with the given translator.
// It calls the Release method on the provided Translate instance.
func (a *Language) ReleaseTranslator(tr *Translate) {
	tr.Release()
}

// release releases the translator resources associated with the context
func (a *Language) release(c echo.Context) {
	c.Translator().(*Translate).Release()
}

// GetLangsMap retrieves a map of available languages with their enabled status.
// It delegates the actual retrieval to the configured langsMapGetter function.
// Returns the language map or an error if retrieval fails.
func (a *Language) GetLangsMap(c echo.Context) (map[string]bool, error) {
	if a.langsMapGetter == nil {
		return LangsMapGetter(c, a)
	}
	return a.langsMapGetter(c, a)
}

// Middleware returns an echo middleware function that handles language detection and translation.
// It detects the language from query parameters, URI, cookies, or Accept-Language header in order.
// If a valid language is found, it sets the language cookie and attaches a translator to the context.
// The middleware will return an error if no valid languages are available or if language detection fails.
func (a *Language) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			list, err := a.GetLangsMap(c)
			if err != nil || len(list) == 0 {
				return err
			}
			lang := c.Query(LangVarName)
			var hasCookie bool
			if !a.Valid(lang, list) {
				lang = a.DetectURI(c, list)
				if !a.Valid(lang, list) {
					lang = c.GetCookie(LangVarName)
					if !a.Valid(lang, list) {
						lang = a.DetectHeader(c.Request(), list)
					} else {
						hasCookie = true
					}
				}
			}
			if !hasCookie {
				c.SetCookie(LangVarName, lang)
			}
			tr := a.AcquireTranslator(lang, list)
			c.OnRelease(a.release)
			c.SetTranslator(tr)
			return h.Handle(c)
		})
	})
}

// Config returns the I18n configuration associated with the Language instance.
func (a *Language) Config() *Config {
	return a.I18n.config
}

// Handler registers HTTP routes for serving i18n messages in JSON and JavaScript formats.
// It provides two endpoints:
// - /i18n.json: returns messages as JSON
// - /i18n.js: returns messages as JavaScript with optional variable assignment
// i18nJSVarName specifies the JavaScript variable name to assign the messages to (optional)
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
