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
		translatePool: &sync.Pool{
			New: func() interface{} {
				return TranslatorNew()
			},
		},
		langsGetter: LangsGetter,
	}
	if len(c) > 0 {
		lang.Init(c[0])
	}
	return lang
}

type Language struct {
	List          map[string]bool                                                                      // 语种列表
	Default       string                                                                               // 默认语种
	I18n          *I18n                                                                                // I18n 实例
	translatePool *sync.Pool                                                                           // 翻译实例池
	langsGetter   func(echo.Context, *Language) (langs map[string]bool, langDefault string, err error) // 获取可用语言列表
}

var (
	LangsGetter = func(c echo.Context, a *Language) (map[string]bool, string, error) {
		return a.List, a.Default, nil
	}
	TranslatorNew = func() Translator {
		return &Translate{_pool: true}
	}
)

// Close closes the I18n instance if it exists.
func (a *Language) Close() {
	if a.I18n != nil {
		a.I18n.Close()
	}
}

func (a *Language) SetTranslatePool(pool *sync.Pool) *Language {
	a.translatePool = pool
	return a
}

func (a *Language) SetTranslatePoolNew(newFunc func() Translator) *Language {
	a.translatePool.New = func() interface{} {
		return newFunc()
	}
	return a
}

// Init initializes the Language instance with the provided configuration.
// It sets up the language list based on the configuration, including the default language.
// If AllList is provided in the config, it registers all listed languages.
// Otherwise, it registers the default language and optionally English ('en') if not the default.
// It also initializes the I18n instance and starts monitoring for changes if Reload is enabled.
func (a *Language) Init(c *Config) {
	c.Init()
	if len(c.AllList) > 0 {
		for _, lang := range c.AllList {
			a.Set(lang, true, lang == c.Default, true)
		}
	} else { // 如果没有指定语言列表，则注册默认语言(DefaultLang:zh-CN)和en
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
	var lang string
	dispatchPath := c.DispatchPath()
	p := strings.TrimPrefix(dispatchPath, `/`)
	if before, _, found := strings.Cut(p, `/`); found {
		lang = before
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
func (a *Language) Valid(lang string, langs map[string]bool) bool {
	if len(lang) == 0 {
		return false
	}
	if langs == nil {
		langs = a.List
	}
	if on, ok := langs[lang]; ok {
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
	return ``
}

// AcquireTranslator gets a Translate instance from the pool and initializes it with the specified language code.
// The returned translator is ready to use for translation operations.
func (a *Language) AcquireTranslator(ctx echo.Context, langCode string, langs map[string]bool, langDefault string) Translator {
	tr := a.translatePool.Get().(Translator)
	tr.Reset(ctx, langCode, a, langs, langDefault)
	return tr
}

// ReleaseTranslator releases the resources associated with the given translator.
// It calls the Release method on the provided Translate instance.
func (a *Language) ReleaseTranslator(tr Translator) {
	tr.Release()
}

// release releases the translator resources associated with the context
func (a *Language) release(c echo.Context) {
	c.Translator().(Translator).Release()
}

// GetLangs retrieves available languages and default language from context.
// Returns:
//   - map of available languages (key: language code, value: enabled status)
//   - default language code
//   - error if any occurred during retrieval
//
// If no languages are found, it creates a default map with the default language and 'en' as fallback.
func (a *Language) GetLangs(c echo.Context) (map[string]bool, string, error) {
	langs, langDefault, err := a.getLangs(c)
	if err != nil {
		return nil, ``, err
	}
	if len(langs) == 0 { // 如果没有指定语言列表，则注册默认语言(DefaultLang:zh-CN)和en
		if len(langDefault) == 0 {
			langDefault = DefaultLang
		}
		langs = map[string]bool{langDefault: true}
		if langDefault != `en` {
			langs[`en`] = true
		}
	}
	return langs, langDefault, err
}

// getLangs retrieves the available languages and the current language from the context.
// If langsGetter is nil, it uses the default LangsGetter function.
// Returns a map of available languages, the current language code, and any error encountered.
func (a *Language) getLangs(c echo.Context) (map[string]bool, string, error) {
	if a.langsGetter == nil {
		return LangsGetter(c, a)
	}
	return a.langsGetter(c, a)
}

// Middleware returns an echo middleware function that handles language detection and translation.
// It detects the language from query parameters, URI, cookies, or Accept-Language header in order.
// If a valid language is found, it sets the language cookie and attaches a translator to the context.
// The middleware will return an error if no valid languages are available or if language detection fails.
func (a *Language) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			langs, langDefault, err := a.GetLangs(c)
			if err != nil {
				return err
			}
			lang := c.Query(LangVarName)
			var hasCookie bool
			if !a.Valid(lang, langs) {
				lang = a.DetectURI(c, langs)
				if !a.Valid(lang, langs) {
					lang = c.GetCookie(LangVarName)
					if !a.Valid(lang, langs) {
						lang = a.DetectHeader(c.Request(), langs)
					} else {
						hasCookie = true
					}
				}
			}
			if len(lang) == 0 {
				lang = langDefault
			}
			if !hasCookie {
				c.SetCookie(LangVarName, lang)
			}
			tr := a.AcquireTranslator(c, lang, langs, langDefault)
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

// Handler registers routes for i18n language data retrieval.
// It registers two routes: "/i18n.json" and "/i18n.js".
// The first route returns the i18n language data as JSON.
// The second route returns the i18n language data as a JavaScript file.
// The function takes the following parameters:
//   - e: the echo.RouteRegister instance.
//   - i18nJSVarName: the JavaScript variable name for the i18n language data.
//   - uriLangVarName: the URI query parameter name for retrieving the language code.
//     If not provided, it defaults to "lng".
func (a *Language) Handler(e echo.RouteRegister, i18nJSVarName string, uriLangVarName ...string) {
	var langVarName string
	if len(uriLangVarName) > 0 && len(uriLangVarName[0]) > 0 {
		langVarName = uriLangVarName[0]
	} else {
		langVarName = `lng`
	}
	getLangCode := func(c echo.Context) string {
		langCode := c.Form(langVarName)
		if len(langCode) == 0 || !a.Valid(langCode, nil) {
			langCode = c.Lang().String()
		}
		return langCode
	}
	e.Get(`/i18n.json`, func(c echo.Context) error {
		t := a.I18n.Get(getLangCode(c))
		if t != nil {
			messages := t.Messages()
			return c.JSON(messages)
		}
		return c.JSONBlob([]byte(`{}`))
	})
	e.Get(`/i18n.js`, func(c echo.Context) error {
		t := a.I18n.Get(getLangCode(c))
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
