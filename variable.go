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
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/webx-top/echo/param"
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

	ErrUnsupportedMediaType         error = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrBadRequest                   error = NewHTTPError(http.StatusBadRequest)
	ErrPaymentRequired              error = NewHTTPError(http.StatusPaymentRequired)
	ErrNotAcceptable                error = NewHTTPError(http.StatusNotAcceptable)
	ErrProxyAuthRequired            error = NewHTTPError(http.StatusProxyAuthRequired)
	ErrRequestTimeout               error = NewHTTPError(http.StatusRequestTimeout)
	ErrConflict                     error = NewHTTPError(http.StatusConflict)
	ErrGone                         error = NewHTTPError(http.StatusGone)
	ErrLengthRequired               error = NewHTTPError(http.StatusLengthRequired)
	ErrPreconditionFailed           error = NewHTTPError(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge        error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong            error = NewHTTPError(http.StatusRequestURITooLong)
	ErrRequestedRangeNotSatisfiable error = NewHTTPError(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed            error = NewHTTPError(http.StatusExpectationFailed)
	ErrUnprocessableEntity          error = NewHTTPError(http.StatusUnprocessableEntity)
	ErrLocked                       error = NewHTTPError(http.StatusLocked)
	ErrFailedDependency             error = NewHTTPError(http.StatusFailedDependency)
	ErrTooEarly                     error = NewHTTPError(http.StatusTooEarly)
	ErrUpgradeRequired              error = NewHTTPError(http.StatusUpgradeRequired)
	ErrPreconditionRequired         error = NewHTTPError(http.StatusPreconditionRequired)
	ErrTooManyRequests              error = NewHTTPError(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge  error = NewHTTPError(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons   error = NewHTTPError(http.StatusUnavailableForLegalReasons)
	ErrNotImplemented               error = NewHTTPError(http.StatusNotImplemented)
	ErrNotFound                     error = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                 error = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                    error = NewHTTPError(http.StatusForbidden)
	ErrStatusRequestEntityTooLarge  error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrMethodNotAllowed             error = NewHTTPError(http.StatusMethodNotAllowed)
	ErrRendererNotRegistered              = errors.New("renderer not registered")
	ErrInvalidRedirectCode                = errors.New("invalid redirect status code")
	ErrNotFoundFileInput                  = errors.New("the specified name file input was not found")

	//----------------
	// Error handlers
	//----------------

	NotFoundHandler = HandlerFunc(func(c Context) error {
		return ErrNotFound
	})

	ErrorHandler = func(err error) Handler {
		return HandlerFunc(func(c Context) error {
			return err
		})
	}

	MethodNotAllowedHandler = HandlerFunc(func(c Context) error {
		return ErrMethodNotAllowed
	})

	_ MiddlewareFuncd = func(h Handler) HandlerFunc {
		return func(c Context) error {
			return h.Handle(c)
		}
	}

	//----------------
	// Shortcut
	//----------------

	StringerMapStart = param.StringerMapStart
	StoreStart       = param.StoreStart
	HStart           = param.StoreStart

	//Custom global variable
	globalVars = param.NewMap()
)

func Set(key, value any) {
	globalVars.Set(key, value)
}

func Get(key any, defaults ...any) any {
	return globalVars.Get(key, defaults...)
}

func GetStoreByKeys(key any, keys ...string) H {
	st, ok := Get(key).(H)
	if !ok {
		if st == nil {
			st = H{}
		}
		return st
	}
	return st.GetStoreByKeys(keys...)
}

func GetOk(key any) (any, bool) {
	return globalVars.GetOk(key)
}

func Has(key any) bool {
	return globalVars.Has(key)
}

func Delete(key any) {
	globalVars.Delete(key)
}

func Range(f func(key, value any) bool) {
	globalVars.Range(f)
}

func GetOrSet(key, value any) (actual any, loaded bool) {
	return globalVars.GetOrSet(key, value)
}

func String(key any, defaults ...any) string {
	return globalVars.String(key, defaults...)
}

func Split(key any, sep string, limit ...int) param.StringSlice {
	return globalVars.Split(key, sep, limit...)
}

func Trim(key any, defaults ...any) param.String {
	return globalVars.Trim(key, defaults...)
}

func HTML(key any, defaults ...any) template.HTML {
	return globalVars.HTML(key, defaults...)
}

func HTMLAttr(key any, defaults ...any) template.HTMLAttr {
	return globalVars.HTMLAttr(key, defaults...)
}

func JS(key any, defaults ...any) template.JS {
	return globalVars.JS(key, defaults...)
}

func CSS(key any, defaults ...any) template.CSS {
	return globalVars.CSS(key, defaults...)
}

func Bool(key any, defaults ...any) bool {
	return globalVars.Bool(key, defaults...)
}

func Float64(key any, defaults ...any) float64 {
	return globalVars.Float64(key, defaults...)
}

func Float32(key any, defaults ...any) float32 {
	return globalVars.Float32(key, defaults...)
}

func Int8(key any, defaults ...any) int8 {
	return globalVars.Int8(key, defaults...)
}

func Int16(key any, defaults ...any) int16 {
	return globalVars.Int16(key, defaults...)
}

func Int(key any, defaults ...any) int {
	return globalVars.Int(key, defaults...)
}

func Int32(key any, defaults ...any) int32 {
	return globalVars.Int32(key, defaults...)
}

func Int64(key any, defaults ...any) int64 {
	return globalVars.Int64(key, defaults...)
}

func Decr(key any, n int64, defaults ...any) int64 {
	return globalVars.Decr(key, n, defaults...)
}

func Incr(key any, n int64, defaults ...any) int64 {
	return globalVars.Incr(key, n, defaults...)
}

func Uint8(key any, defaults ...any) uint8 {
	return globalVars.Uint8(key, defaults...)
}

func Uint16(key any, defaults ...any) uint16 {
	return globalVars.Uint16(key, defaults...)
}

func Uint(key any, defaults ...any) uint {
	return globalVars.Uint(key, defaults...)
}

func Uint32(key any, defaults ...any) uint32 {
	return globalVars.Uint32(key, defaults...)
}

func Uint64(key any, defaults ...any) uint64 {
	return globalVars.Uint64(key, defaults...)
}

func Timestamp(key any, defaults ...any) time.Time {
	return globalVars.Timestamp(key, defaults...)
}

func DateTime(key any, layouts ...string) time.Time {
	return globalVars.DateTime(key, layouts...)
}

func Children(key any, keys ...any) Store {
	r := GetStore(key)
	for _, key := range keys {
		r = GetStore(key)
	}
	return r
}

func GetStore(key any, defaults ...any) Store {
	return AsStore(globalVars.Get(key, defaults...))
}
