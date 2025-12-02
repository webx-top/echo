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
	"errors"

	"github.com/webx-top/echo"
)

// NewTranslate creates a new Translate instance and initializes it with the given language and language object.
// Returns a pointer to the initialized Translate struct.
func NewTranslate(ctx echo.Context, language string, langObject *Language, langs map[string]bool, langDefault string) *Translate {
	tr := &Translate{}
	tr.Reset(ctx, language, langObject, langs, langDefault)
	return tr
}

type Translator interface {
	Reset(ctx echo.Context, language string, langObject *Language, langs map[string]bool, langDefault string) Translator
	echo.Translator
	echo.Releaseable
}

var _ echo.Translator = (*Translate)(nil)
var _ echo.Releaseable = (*Translate)(nil)

type Translate struct {
	code        echo.LangCode
	lang        *Language
	langs       map[string]bool
	langDefault string
	_pool       bool
}

// Release releases the Translate instance resources and returns it to the pool if it was created from a pool.
// It sets both code and lang references to nil and optionally returns the instance to the translatePool.
func (t *Translate) Release() {
	t.code = nil
	if t._pool {
		t.lang.translatePool.Put(t)
	}
	t.lang = nil
	t.langs = nil
	t.langDefault = ``
}

// Reset sets the language code and language object for the translator
// language: the language code to set
// langObject: the language object containing translations
// Returns the modified Translate instance for method chaining
func (t *Translate) Reset(_ echo.Context, language string, langObject *Language, langs map[string]bool, langDefault string) Translator {
	t.code = echo.NewLangCode(language)
	t.lang = langObject
	t.langs = langs
	t.langDefault = langDefault
	return t
}

// T translates the given format string using the current language code and optional arguments.
// Returns the translated string.
func (t *Translate) T(format string, args ...interface{}) string {
	return t.lang.I18n.T(t.code.String(), format, args...)
}

// E returns a new error with the translated message using the given format string and arguments.
func (t *Translate) E(format string, args ...interface{}) error {
	return errors.New(t.T(format, args...))
}

// Lang returns the language code of the Translate instance
func (t *Translate) Lang() echo.LangCode {
	return t.code
}

// SetLang sets the language code for translation using the specified language string.
func (t *Translate) SetLang(lang string) {
	t.code.Reset(lang)
}

// LangDefault default language
func (t *Translate) LangDefault() string {
	return t.langDefault
}

// LangList language list
func (t *Translate) LangList() []string {
	return t.lang.Config().AllList
}

// LangExists language code exists
func (t *Translate) LangExists(langCode string) bool {
	return t.lang.Valid(langCode, t.langs)
}
