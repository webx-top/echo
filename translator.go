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

// E creates a new error using the given format string and optional arguments.
// If arguments are provided, it uses fmt.Errorf to format the error message.
// If no arguments are provided, it creates a new error with the format string as-is.
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
	Reset(langCode string, separator ...string)
	String() string                                          // all lowercase characters
	Normalize() string                                       // lowercase(language)-uppercase(region)
	Format(regionUppercase bool, separator ...string) string // language-region
	Language() string                                        // language
	Region(regionUppercase bool) string                      // region
}

type Multilingual interface {
	LangDefault() string             // default language
	LangList() []string              // language list
	LangExists(langCode string) bool // language code exists
	SetLang(lang string)             // set language
}

// NewLangCode creates a new LangCode from a language string with optional separator.
// The language string is split into language and region parts using the separator (default is "-").
// The language part is converted to lowercase, and the region part (if exists) is converted to lowercase and uppercase variants.
func NewLangCode(language string, separator ...string) LangCode {
	l := &langCode{}
	l.Reset(language, separator...)
	return l
}

type langCode struct {
	language        string
	region          string
	regionLowercase string
}

// Reset sets the language and region codes by parsing the input language string.
// The language string is split using the provided separator (default is "-") into language and region parts.
// The language part is converted to lowercase, and if present, the region part is converted to lowercase and uppercase variants.
func (l *langCode) Reset(language string, separator ...string) {
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
}

// String returns the language code in string format, combining language and region with a hyphen if region exists.
func (l langCode) String() string {
	if len(l.regionLowercase) > 0 {
		return l.language + `-` + l.regionLowercase
	}
	return l.language
}

// Normalize returns the normalized language code string in the format "language-region" if region is specified,
// otherwise returns just the language code.
func (l langCode) Normalize() string {
	if len(l.region) > 0 {
		return l.language + `-` + l.region
	}
	return l.language
}

// Language returns the language code string
func (l langCode) Language() string {
	return l.language
}

// Region returns the region part of the language code, either in uppercase or lowercase format depending on the regionUppercase parameter.
func (l langCode) Region(regionUppercase bool) string {
	if regionUppercase {
		return l.region
	}
	return l.regionLowercase
}

// Format formats the language code with optional region code.
// If regionUppercase is true, uses uppercase region code, otherwise uses lowercase.
// separator specifies the delimiter between language and region (defaults to "-").
// Returns the formatted language-region string or just language if no region.
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
	code: &langCode{
		language: `en`,
	},
}

type NopTranslate struct {
	code LangCode
}

// trimTranslatorGroupPrefix removes the translator group prefix from the format string.
// If the format starts with '#', it splits the string by '#' and returns the part after the second '#'.
// Otherwise, it returns the original format string unchanged.
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

// T translates the given format string after removing translator group prefix.
// If args are provided, it uses fmt.Sprintf to format the string with the arguments.
// Returns the translated or formatted string.
func (n *NopTranslate) T(format string, args ...interface{}) string {
	format = trimTranslatorGroupPrefix(format)
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

// E wraps the E function to return an error with formatted message
// format: the error message format string
// args: arguments to format into the message
func (n *NopTranslate) E(format string, args ...interface{}) error {
	return E(format, args...)
}

// Lang returns the language code stored in the NopTranslate instance.
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

// SetLang sets the language code for NopTranslate instance
func (n *NopTranslate) SetLang(langCode string) {
	n.code.Reset(langCode)
}
