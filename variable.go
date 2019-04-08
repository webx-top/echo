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
	"net/http"
	"regexp"
	"sync"
)

var (
	splitHTTPMethod = regexp.MustCompile(`[^A-Z]+`)

	methods = []string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}

	//--------
	// Errors
	//--------

	ErrUnsupportedMediaType        error = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    error = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                error = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   error = NewHTTPError(http.StatusForbidden)
	ErrStatusRequestEntityTooLarge error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrMethodNotAllowed            error = NewHTTPError(http.StatusMethodNotAllowed)
	ErrRendererNotRegistered             = errors.New("renderer not registered")
	ErrInvalidRedirectCode               = errors.New("invalid redirect status code")
	ErrNotFoundFileInput                 = errors.New("The specified name file input was not found")

	//----------------
	// Error handlers
	//----------------

	NotFoundHandler = HandlerFunc(func(c Context) error {
		return ErrNotFound
	})

	MethodNotAllowedHandler = HandlerFunc(func(c Context) error {
		return ErrMethodNotAllowed
	})

	_ MiddlewareFuncd = func(h Handler) HandlerFunc {
		return func(c Context) error {
			return h.Handle(c)
		}
	}

	globalVars = sync.Map{} //Custom global variable
)

func Set(key, value interface{}) {
	globalVars.Store(key, value)
}

func Get(key interface{}, defaults ...interface{}) interface{} {
	value, ok := globalVars.Load(key)
	if !ok && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() interface{}); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func GetOk(key interface{}) (interface{}, bool) {
	return globalVars.Load(key)
}

func Exists(key interface{}) bool {
	_, ok := globalVars.Load(key)
	return ok
}

func Delete(key interface{}) {
	globalVars.Delete(key)
}

func Range(f func(key, value interface{}) bool) {
	globalVars.Range(f)
}

func GetOrSet(key, value interface{}) (actual interface{}, loaded bool) {
	return globalVars.LoadOrStore(key, value)
}
