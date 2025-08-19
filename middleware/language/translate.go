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

func NewTranslate(language string, langObject *Language) *Translate {
	tr := &Translate{}
	return tr.Reset(language, langObject)
}

type Translate struct {
	code  echo.LangCode
	lang  *Language
	_pool bool
}

func (t *Translate) Release() {
	if t._pool {
		t.lang.translatePool.Put(t)
	}
	t.code = nil
	t.lang = nil
}

func (t *Translate) Reset(language string, langObject *Language) *Translate {
	t.code = echo.NewLangCode(language)
	t.lang = langObject
	return nil
}

func (t *Translate) T(format string, args ...interface{}) string {
	return t.lang.I18n.T(t.code.String(), format, args...)
}

func (t *Translate) E(format string, args ...interface{}) error {
	return errors.New(t.T(format, args...))
}

func (t *Translate) Lang() echo.LangCode {
	return t.code
}

func (t *Translate) SetLang(lang string) {
	t.code = echo.NewLangCode(lang)
}

// LangDefault default language
func (t *Translate) LangDefault() string {
	return t.lang.Default
}

// LangList language list
func (t *Translate) LangList() []string {
	return t.lang.Config().AllList
}

// LangExists language code exists
func (t *Translate) LangExists(langCode string) bool {
	return t.lang.Valid(langCode)
}
