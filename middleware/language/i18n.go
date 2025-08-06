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
	"github.com/webx-top/echo"
	"golang.org/x/sync/singleflight"
)

var defaultInstance *I18n

type I18n struct {
	*i18n.TranslatorFactory
	lock        sync.RWMutex
	translators map[string]*i18n.Translator
	config      *Config
	sg          singleflight.Group
}

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

	if defaultInstance == nil {
		SetDefault(ins)
	}
	return ins
}

func SetDefault(ins *I18n) {
	defaultInstance = ins
}

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

func (a *I18n) Get(langCode string) *i18n.Translator {
	a.lock.RLock()
	t, ok := a.translators[langCode]
	a.lock.RUnlock()
	if !ok {
		t = a.GetAndCache(langCode)
	}
	return t
}

func (a *I18n) Translate(langCode, key string, args map[string]string) string {
	t := a.Get(langCode)
	translation, errs := t.Translate(key, args)
	if errs != nil {
		return i18n.TrimGroupPrefix(key)
	}
	return translation
}

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

// T 多语言翻译
func T(langCode, key string, args ...interface{}) (t string) {
	if defaultInstance == nil {
		t = i18n.TrimGroupPrefix(key)
		if len(args) > 0 {
			t = fmt.Sprintf(t, args...)
		}
		return
	}
	t = defaultInstance.T(langCode, key, args...)
	return
}
