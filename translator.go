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

package echo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/admpub/events"
)

// T 标记为多语言文本(fake)
func T(format string, _ ...interface{}) string {
	return format
}

func E(format string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf(format, args...)
	}
	return errors.New(format)
}

type (
	translator      = Translator
	eventsEmitterer = events.Emitterer
)

type Translator interface {
	T(format string, args ...interface{}) string
	E(format string, args ...interface{}) error
	Lang() LangCode
	Multilingual
}

type LangCode interface {
	String() string    // all lowercase characters
	Normalize() string // lowercase(language)-uppercase(region)
	Format(regionUppercase bool, separator ...string) string
	Language() string
	Region(regionUppercase bool) string
}

type Multilingual interface {
	LangDefault() string             // default language
	LangList() []string              // language list
	LangExists(langCode string) bool // language code exists
}

func NewLangCode(language string, separator ...string) LangCode {
	l := langCode{}
	sep := `-`
	if len(separator) > 0 && len(separator[0]) > 0 {
		sep = separator[0]
	}
	lg := strings.SplitN(language, sep, 2)
	switch len(lg) {
	case 2:
		l.regionLowercase = strings.ToLower(lg[1])
		l.region = strings.ToUpper(lg[1])
		fallthrough
	case 1:
		l.language = strings.ToLower(lg[0])
	}
	return l
}

type langCode struct {
	language        string
	region          string
	regionLowercase string
}

func (l langCode) String() string {
	if len(l.regionLowercase) > 0 {
		return l.language + `-` + l.regionLowercase
	}
	return l.language
}

func (l langCode) Normalize() string {
	if len(l.region) > 0 {
		return l.language + `-` + l.region
	}
	return l.language
}

func (l langCode) Language() string {
	return l.language
}

func (l langCode) Region(regionUppercase bool) string {
	if regionUppercase {
		return l.region
	}
	return l.regionLowercase
}

func (l langCode) Format(regionUppercase bool, separator ...string) string {
	var region string
	if regionUppercase {
		region = l.region
	} else {
		region = l.regionLowercase
	}
	if len(region) > 0 {
		if len(separator) > 0 {
			return l.language + separator[0] + region
		}
		return l.language + `-` + region
	}
	return l.language
}

var DefaultNopTranslate Translator = &NopTranslate{
	code: langCode{
		language: `en`,
	},
}

type NopTranslate struct {
	code LangCode
}

func trimTranslatorGroupPrefix(format string) string {
	if len(format) > 1 && format[0] == '#' {
		s := format[1:]
		parts := strings.SplitN(s, `#`, 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return format
}

func (n *NopTranslate) T(format string, args ...interface{}) string {
	format = trimTranslatorGroupPrefix(format)
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

func (n *NopTranslate) E(format string, args ...interface{}) error {
	return E(format, args...)
}

func (n *NopTranslate) Lang() LangCode {
	return n.code
}

// LangDefault default language
func (n *NopTranslate) LangDefault() string {
	return n.code.Normalize()
}

// LangList language list
func (n *NopTranslate) LangList() []string {
	return []string{n.code.Normalize()}
}

// LangExists language code exists
func (n *NopTranslate) LangExists(langCode string) bool {
	return langCode == n.code.Normalize()
}
