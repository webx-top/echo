package oauth2

import (
	"net/http"

	"github.com/markbates/goth"
	"github.com/webx-top/echo"
)

// OAuth is a plugin which helps you to use OAuth/OAuth2 apis from famous websites
type OAuth struct {
	Config          Config
	URLPrefix       string
	successHandlers []interface{}
	failHandler     echo.HTTPErrorHandler
}

// New returns a new OAuth plugin
// receives one parameter of type 'Config'
func New(cfg Config) *OAuth {
	c := DefaultConfig().MergeSingle(cfg)
	return &OAuth{Config: c}
}

// Success registers handler(s) which fires when the user logged in successfully
func (p *OAuth) Success(handlersFn ...interface{}) {
	p.successHandlers = append(p.successHandlers, handlersFn...)
}

// Fail registers handler which fires when the user failed to logged in
// underhood it justs registers an error handler to the StatusUnauthorized(400 status code), same as 'iris.OnError(400,handler)'
func (p *OAuth) Fail(handler echo.HandlerFunc) {
	p.failHandler = handler
}

// User returns the user for the particular client
// if user is not validated  or not found it returns nil
// same as 'ctx.Get(config's ContextKey field).(goth.User)'
func (p *OAuth) User(ctx echo.Context) (u goth.User) {
	return ctx.Get(p.Config.ContextKey).(goth.User)
}

// PreBuild plugin in order to register the oauth route
func (p *OAuth) Wrapper(e *echo.Echo) {
	oauthProviders := p.Config.GenerateProviders(p.URLPrefix)
	if len(oauthProviders) > 0 {
		goth.UseProviders(oauthProviders...)
		// set the mux path to handle the registered providers
		e.Get(p.Config.Path+"/:provider", func(ctx echo.Context) {
			err := BeginAuthHandler(ctx)
			if err != nil {
				e.Logger.Error("[OAUTH2] Error: " + err.Error())
			}
		})
		//println("registered " + p.Config.Path + "/:provider")

		authMiddleware := func(h echo.Handler) echo.Handler {
			return echo.HandlerFunc(func(ctx echo.Context) {
				user, err := CompleteUserAuth(ctx)
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
				}
				ctx.Set(p.Config.ContextKey, user)
				return h.Handle(ctx)
			})
		}

		p.successHandlers = append([]echo.HandlerFunc{authMiddleware}, p.successHandlers...)
		var handlers []interface{}
		for i := len(p.successHandlers) - 1; i >= 0; i-- {
			handlers = append(handlers, p.successHandlers[i])
		}
		e.Get(p.Config.Path+"/:provider/callback", handlers...)
		// register the error handler
		if p.failHandler != nil {
			s.SetHTTPErrorHandler(p.failHandler)
		}
	}
}
