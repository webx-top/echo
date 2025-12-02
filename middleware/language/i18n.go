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
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/admpub/i18n"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type I18n struct {
	*i18n.TranslatorFactory
	lock        sync.RWMutex
	translators map[string]*i18n.Translator
	config      *Config
	monitor     *com.MonitorEvent
}

// fixedPath resolves a relative path to an absolute path if needed.
// If ppath is empty or already absolute, returns it unchanged.
// If the path exists as-is (checked using the provided open function), returns it unchanged.
// Otherwise, joins the path with the current working directory (echo.Wd()).
// The open function is used to check path existence, defaulting to http.Dir if nil.
func fixedPath(ppath string, open func(string) http.FileSystem) string {
	if len(ppath) == 0 || filepath.IsAbs(ppath) {
		return ppath
	}
	if open == nil {
		open = func(p string) http.FileSystem {
			return http.Dir(p)
		}
	}
	fs := open(ppath)
	var exists bool
	file, err := fs.Open(`.`)
	if err == nil {
		_, err = file.Stat()
		exists = err == nil
		file.Close()
	}
	if exists {
		return ppath
	}
	return filepath.Join(echo.Wd(), ppath)
}

// NewI18n creates a new I18n instance with the given configuration.
// It initializes translator paths, creates a translator factory, and caches the default translator.
// If any errors occur during initialization, it will panic with the aggregated error messages.
// The first created instance will be set as the default instance if none exists.
func NewI18n(c *Config) *I18n {
	for index, value := range c.RulesPath {
		c.RulesPath[index] = fixedPath(value, c.FSFunc())
	}
	for index, value := range c.MessagesPath {
		c.MessagesPath[index] = fixedPath(value, c.FSFunc())
	}
	f, errs := i18n.NewTranslatorFactoryWith(c.Project, c.RulesPath, c.MessagesPath, c.Fallback, c.FSFunc())
	if len(errs) > 0 {
		var errMsg string
		for idx, err := range errs {
			if idx > 0 {
				errMsg += "\n"
			}
			errMsg += err.Error()
		}
		if len(errMsg) > 0 {
			panic("== i18n error: " + errMsg + "\n")
		}
	}
	ins := &I18n{
		TranslatorFactory: f,
		translators:       make(map[string]*i18n.Translator),
		config:            c,
	}
	ins.GetAndCache(c.Default)
	return ins
}

// Close stops the I18n monitor if it exists.
func (a *I18n) Close() {
	if a.monitor == nil {
		return
	}
	a.monitor.Close()
}

// GetAndCache retrieves a translator for the specified language code and caches it.
// If the translator for the requested language cannot be loaded, it falls back to the default language.
// Returns the cached translator instance.
// Panics if errors occur while loading both the requested and default language translators.
func (a *I18n) GetAndCache(langCode string) *i18n.Translator {
	var (
		t    *i18n.Translator
		errs []error
	)
	t, errs = a.TranslatorFactory.GetTranslator(langCode)
	if len(errs) > 0 {
		if a.config.Default != langCode {
			t, errs = a.TranslatorFactory.GetTranslator(a.config.Default)
		}
	}
	if len(errs) > 0 {
		var errMsg string
		for idx, err := range errs {
			if idx > 0 {
				errMsg += "\n"
			}
			errMsg += err.Error()
		}
		if len(errMsg) > 0 {
			panic("== i18n error: " + errMsg + "\n")
		}
	}
	a.lock.Lock()
	a.translators[langCode] = t
	a.lock.Unlock()
	return t
}

// Get returns the translator for the specified language code.
// If the translator is not cached, it will be loaded and cached using GetAndCache.
// The returned translator is safe for concurrent use.
func (a *I18n) Get(langCode string) *i18n.Translator {
	a.lock.RLock()
	t, ok := a.translators[langCode]
	a.lock.RUnlock()
	if !ok {
		t = a.GetAndCache(langCode)
	}
	return t
}

// Translate translates the given key for the specified language code, using the provided args for variable substitution.
// If translation fails, returns the key with any group prefix removed.
// langCode: target language code
// key: translation key
// args: variables to substitute in the translation
// Returns: translated string or cleaned key if translation fails
func (a *I18n) Translate(langCode, key string, args map[string]string) string {
	t := a.Get(langCode)
	translation, errs := t.Translate(key, args)
	if errs != nil {
		return i18n.TrimGroupPrefix(key)
	}
	return translation
}

// T translates the given key for the specified language code, optionally formatting the result with provided arguments.
// If args[0] is a map[string]string, it's used as translation variables. Otherwise, args are used for fmt.Sprintf formatting.
// Returns the translated string, formatted if arguments were provided.
func (a *I18n) T(langCode, key string, args ...interface{}) (t string) {
	if len(args) > 0 {
		if v, ok := args[0].(map[string]string); ok {
			t = a.Translate(langCode, key, v)
			return
		}
		t = a.Translate(langCode, key, nil)
		t = fmt.Sprintf(t, args...)
		return
	}
	t = a.Translate(langCode, key, nil)
	return
}
