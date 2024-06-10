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
	"regexp"
	"strconv"
	"strings"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c Context) bool
)

func CaptureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindStringSubmatch(input)
	if len(groups) == 0 {
		return nil
	}
	values := groups[1:]
	kv := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		if i != 0 && len(name) > 0 {
			kv[name] = groups[i]
		}
	}
	return CaptureTokensByValues(values, kv)
}

func CaptureTokensByValues(values []string, kv map[string]string, quoted ...bool) *strings.Replacer {
	var prefix string
	if len(quoted) > 0 && quoted[0] {
		prefix = `\`
	}
	replace := make([]string, 0, 2*(len(values)+len(kv)))
	for i, v := range values {
		replace = append(replace, prefix+"$"+strconv.Itoa(i+1), v)
	}
	for k, v := range kv {
		replace = append(replace, prefix+"{"+k+prefix+"}", v)
	}
	return strings.NewReplacer(replace...)
}

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(c Context) bool {
	return false
}
