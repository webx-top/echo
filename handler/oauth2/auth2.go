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
package oauth2

import (
	"net/http"

	"github.com/admpub/goth"
	"github.com/webx-top/echo"
)

// OAuth is a plugin which helps you to use OAuth/OAuth2 apis from famous websites
type OAuth struct {
	Config              *Config
	HostURL             string
	successHandlers     []interface{}
	failHandler         echo.HTTPErrorHandler
	beginAuthHandler    echo.Handler
	completeAuthHandler func(ctx echo.Context) (goth.User, error)
}

// New returns a new OAuth plugin
// receives one parameter of type 'Config'
func New(hostURL string, cfg *Config) *OAuth {
	c := DefaultConfig().MergeSingle(cfg)
	c.Host = hostURL
	return &OAuth{
		Config:              &c,
		beginAuthHandler:    echo.HandlerFunc(BeginAuthHandler),
		completeAuthHandler: CompleteUserAuth,
	}
}

// SetSuccessHandler registers handler(s) which fires when the user logged in successfully
func (p *OAuth) SetSuccessHandler(handlersFn ...interface{}) *OAuth {
	p.successHandlers = handlersFn
	return p
}

// AddSuccessHandler registers handler(s) which fires when the user logged in successfully
func (p *OAuth) AddSuccessHandler(handlersFn ...interface{}) *OAuth {
	p.successHandlers = append(p.successHandlers, handlersFn...)
	return p
}

// SetFailHandler registers handler which fires when the user failed to logged in
// underhood it justs registers an error handler to the StatusUnauthorized(400 status code), same as 'iris.OnError(400,handler)'
func (p *OAuth) SetFailHandler(handler echo.HTTPErrorHandler) *OAuth {
	p.failHandler = handler
	return p
}

func (p *OAuth) SetBeginAuthHandler(handler echo.Handler) *OAuth {
	p.beginAuthHandler = handler
	return p
}

func (p *OAuth) SetCompleteAuthHandler(handler func(ctx echo.Context) (goth.User, error)) *OAuth {
	p.completeAuthHandler = handler
	return p
}

// User returns the user for the particular client
// if user is not validated  or not found it returns nil
// same as 'ctx.Get(config's ContextKey field).(goth.User)'
func (p *OAuth) User(ctx echo.Context) (u goth.User) {
	u, _ = ctx.Internal().Get(p.Config.ContextKey).(goth.User)
	return u
}

func MiddlewareVerifyProvider(config *Config) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			providerName, err := GetProviderName(ctx)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}
			account := config.GetAccount(providerName)
			if account == nil {
				return echo.ErrNotFound
			}
			return h.Handle(ctx)
		}
	}
}

func (p *OAuth) MiddlewareAuth(h echo.Handler) echo.Handler {
	return echo.HandlerFunc(func(ctx echo.Context) error {
		user, err := p.completeAuthHandler(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error()).SetRaw(err)
		}
		ctx.Internal().Set(p.Config.ContextKey, user)
		return h.Handle(ctx)
	})
}

// Wrapper register the oauth route
func (p *OAuth) Wrapper(e *echo.Echo, middlewares ...interface{}) {
	p.Config.GenerateProviders()

	g := e.Group(p.Config.Path, append([]interface{}{MiddlewareVerifyProvider(p.Config)}, middlewares...)...)

	successHandlers := append([]interface{}{p.MiddlewareAuth}, p.successHandlers...)
	lastIndex := len(successHandlers) - 1
	var callbackHandler echo.Handler
	var callbackMiddlewares []interface{}
	if lastIndex == 0 {
		callbackHandler = echo.HandlerFunc(func(ctx echo.Context) error {
			return ctx.String(`Success Handler is not set`)
		})
		callbackMiddlewares = successHandlers
	} else {
		callbackHandler = echo.WrapHandler(successHandlers[lastIndex])
		callbackMiddlewares = successHandlers[0:lastIndex]
	}

	// set the mux path to handle the registered providers
	g.Get("/login/:provider", func(ctx echo.Context) error {
		// try to get the user without re-authenticating
		user, err := fetchUser(ctx)
		if err != nil {
			return p.beginAuthHandler.Handle(ctx)
		}
		ctx.Internal().Set(p.Config.ContextKey, user)
		return callbackHandler.Handle(ctx)
	})

	g.Get("/callback/:provider", callbackHandler, callbackMiddlewares...)

	// register the error handler
	if p.failHandler != nil {
		e.SetHTTPErrorHandler(p.failHandler)
	}
}
