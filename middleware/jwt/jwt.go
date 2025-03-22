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
package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/webx-top/echo"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Signing key to validate token.
		// Required.
		SigningKey interface{} `json:"signing_key"`

		// Signing method, used to check token signing method.
		// Optional. Default value HS256.
		SigningMethod string `json:"signing_method"`

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey string `json:"context_key"`

		// Claims are extendable claims data defining token content.
		// Optional. Default value jwt.MapClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "param:<name>"
		// - "cookie:<name>"
		TokenLookup string `json:"token_lookup"`

		OnErrorAbort bool `json:"on_error_abort"`

		errorHandler      func(c echo.Context, err error)
		fallbackExtractor func(c echo.Context) (string, error)
		tokenPreprocessor func(c echo.Context, token string) (string, error)
		keyFunc           jwt.Keyfunc
	}

	jwtExtractor func(echo.Context) (string, error)
)

func (j *JWTConfig) SetErrorHandler(errorHandler func(c echo.Context, err error)) *JWTConfig {
	j.errorHandler = errorHandler
	return j
}

func (j *JWTConfig) SetTokenPreprocessor(tokenPreprocessor func(c echo.Context, token string) (string, error)) *JWTConfig {
	j.tokenPreprocessor = tokenPreprocessor
	return j
}

func (j *JWTConfig) ErrorHandler() func(c echo.Context, err error) {
	return j.errorHandler
}

func (j *JWTConfig) SetFallbackExtractor(extractor func(c echo.Context) (string, error)) *JWTConfig {
	j.fallbackExtractor = extractor
	return j
}

func (j *JWTConfig) FallbackExtractor() func(c echo.Context) (string, error) {
	return j.fallbackExtractor
}

const (
	bearer = "Bearer"
)

// Algorithims
const (
	AlgorithmHS256 = "HS256"
)

// Errors
var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:       echo.DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "jwtUser",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		Claims:        jwt.MapClaims{},
		OnErrorAbort:  true,
		tokenPreprocessor: func(_ echo.Context, token string) (string, error) {
			return token, nil
		},
	}
)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For empty token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func JWT(key []byte) echo.MiddlewareFuncd {
	c := DefaultJWTConfig
	c.SigningKey = key
	return JWTWithConfig(c)
}

// JWTWithConfig returns a JWT auth middleware with config.
// See: `JWT()`.
func JWTWithConfig(config JWTConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil {
		panic("jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.tokenPreprocessor == nil {
		config.tokenPreprocessor = DefaultJWTConfig.tokenPreprocessor
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.SigningKey, nil
	}

	// Initialize
	parts := strings.SplitN(config.TokenLookup, ":", 2)
	extractor := jwtFromHeader(parts[1])
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "param":
		extractor = jwtFromParam(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	case "any":
		extractor = jwtFromFromAny(parts[1])
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			auth, err := extractor(c)
			if err == ErrJWTMissing {
				if config.fallbackExtractor != nil {
					auth, err = config.fallbackExtractor(c)
				}
			}
			if err == nil {
				auth, err = config.tokenPreprocessor(c, auth)
			}
			if err != nil {
				if config.errorHandler != nil {
					config.errorHandler(c, err)
				}
				if config.OnErrorAbort {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
				}
				return next.Handle(c)
			}
			var token *jwt.Token
			// Issue #647, #656
			if _, ok := config.Claims.(jwt.MapClaims); ok {
				token, err = jwt.Parse(auth, config.keyFunc)
			} else {
				claims := reflect.ValueOf(config.Claims).Interface().(jwt.Claims)
				token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
			}
			if err == nil && token.Valid {
				// Store user information from token into context.
				c.Internal().Set(config.ContextKey, token)
				return next.Handle(c)
			}
			if config.errorHandler != nil {
				config.errorHandler(c, err)
			}
			if config.OnErrorAbort {
				return echo.ErrUnauthorized
			}
			return next.Handle(c)
		}
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from request header.
func jwtFromHeader(header string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header().Get(header)
		l := len(bearer)
		if len(auth) > l+1 && auth[:l] == bearer {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Query(param)
		var err error
		if len(token) == 0 {
			return "", ErrJWTMissing
		}
		return token, err
	}
}

// jwtFromParam returns a `jwtExtractor` that extracts token from the url param string.
func jwtFromParam(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Param(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

func jwtFromFromAny(key string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Request().Header().Get(key)
		if len(token) > 0 {
			return token, nil
		}
		token = c.Param(key)
		if len(token) > 0 {
			return token, nil
		}
		token = c.Query(key)
		if len(token) == 0 {
			return token, ErrJWTMissing
		}
		return token, nil
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.GetCookie(name)
		if len(token) == 0 {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

// BuildSignedString example: github.com/golang-jwt/jwt/example_test.go
func BuildSignedString(claims jwt.Claims, mySigningKey interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// BuildRegisteredSignedString example: github.com/golang-jwt/jwt/example_test.go
func BuildRegisteredSignedString(claims *jwt.RegisteredClaims, mySigningKey interface{}) (string, error) {
	return BuildSignedString(claims, mySigningKey)
}

// BuildMapSignedString example: github.com/golang-jwt/jwt/example_test.go
func BuildMapSignedString(claims jwt.MapClaims, mySigningKey interface{}) (string, error) {
	return BuildSignedString(claims, mySigningKey)
}
