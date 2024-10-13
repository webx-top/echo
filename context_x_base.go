package echo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/admpub/events"

	pkgCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/poolx/bufferpool"
)

type XContext struct {
	Translator
	events.Emitterer
	transaction         *BaseTransaction
	validator           Validator
	sessioner           Sessioner
	cookier             Cookier
	request             engine.Request
	response            engine.Response
	path                string
	pnames              []string
	pvalues             []string
	hnames              []string // host
	hvalues             []string // host
	store               *param.SafeStore
	internal            *param.SafeMap
	handler             Handler
	route               *Route
	rid                 int
	echo                *Echo
	funcs               map[string]interface{}
	renderer            Renderer
	renderDataWrapper   DataWrapper
	sessionOptions      *SessionOptions
	withFormatExtension bool
	defaultExtension    string
	format              string
	code                int
	preResponseHook     []func() error
	dataEngine          Data
	accept              *Accepts
	auto                bool
	onHostFound         func(Context) (bool, error)
	realIP              string
}

var _ context.Context = (*XContext)(nil)

// NewContext creates a Context object.
func NewContext(req engine.Request, res engine.Response, e *Echo) Context {
	c := &XContext{
		validator:         e.Validator,
		Translator:        DefaultNopTranslate,
		Emitterer:         events.Default,
		transaction:       DefaultNopTransaction,
		request:           req,
		response:          res,
		echo:              e,
		pvalues:           make([]string, *e.maxParam),
		internal:          param.NewMap(),
		store:             param.NewSafeStore(),
		handler:           NotFoundHandler,
		sessioner:         DefaultSession,
		onHostFound:       e.onHostFound,
		renderDataWrapper: e.renderDataWrapper,
	}
	c.cookier = NewCookier(c)
	c.dataEngine = NewData(c)
	c.ResetFuncs(e.FuncMap)
	//c.SetValue(ContextKey, c)
	return c
}

func (c *XContext) StdContext() context.Context {
	return c.request.Context()
}

func (c *XContext) WithContext(ctx context.Context) *http.Request {
	return c.request.WithContext(ctx)
}

func (c *XContext) SetValue(key string, value interface{}) {
	c.request.SetValue(key, value)
}

func (c *XContext) Internal() *param.SafeMap {
	return c.internal
}

func (c *XContext) SetEmitterer(emitterer events.Emitterer) {
	c.Emitterer = emitterer
}

func (c *XContext) Handler() Handler {
	return c.handler
}

func (c *XContext) Deadline() (deadline time.Time, ok bool) {
	return c.StdContext().Deadline()
}

func (c *XContext) Done() <-chan struct{} {
	return c.StdContext().Done()
}

func (c *XContext) Err() error {
	return c.StdContext().Err()
}

func (c *XContext) Value(key interface{}) interface{} {
	return c.StdContext().Value(key)
}

func (c *XContext) Handle(ctx Context) error {
	return c.handler.Handle(ctx)
}

func (c *XContext) Route() *Route {
	if c.route == nil {
		if c.rid < 0 || c.rid >= len(c.echo.router.routes) {
			c.route = defaultRoute
		} else {
			c.route = c.echo.router.routes[c.rid]
		}
	}
	return c.route
}

func (c *XContext) SetAuto(on bool) Context {
	c.auto = on
	return c
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *XContext) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

func (c *XContext) NewError(code pkgCode.Code, msg string, args ...interface{}) *Error {
	if len(msg) > 0 {
		msg = c.T(msg, args...)
	}
	return NewError(msg, code).NoClone()
}

func (c *XContext) NewErrorWith(err error, code pkgCode.Code, args ...interface{}) *Error {
	var msg string
	if len(args) > 0 {
		msg = param.AsString(args[0])
		if len(msg) > 0 {
			if len(args) > 1 {
				msg = c.T(msg, args[1:]...)
			} else {
				msg = c.T(msg)
			}
		}
	}
	return NewErrorWith(err, msg, code).NoClone()
}

// Logger returns the `Logger` instance.
func (c *XContext) Logger() logger.Logger {
	return c.echo.logger
}

// Object returns the `context` object.
func (c *XContext) Object() *XContext {
	return c
}

// Echo returns the `Echo` instance.
func (c *XContext) Echo() *Echo {
	return c.echo
}

func (c *XContext) SetTranslator(t Translator) {
	c.Translator = t
}

func (c *XContext) SetDefaultExtension(ext string) {
	c.defaultExtension = ext
}

func (c *XContext) DefaultExtension() string {
	if c.withFormatExtension {
		return `.` + c.Format()
	}
	if len(c.defaultExtension) > 0 {
		return c.defaultExtension
	}
	return c.echo.defaultExtension
}

func (c *XContext) Reset(req engine.Request, res engine.Response) {
	if req != nil {
		req.SetMaxSize(c.echo.MaxRequestBodySize())
	}
	c.validator = c.echo.Validator
	c.Emitterer = events.Default
	c.Translator = DefaultNopTranslate
	c.transaction = DefaultNopTransaction
	c.sessioner = DefaultSession
	c.cookier = NewCookier(c)
	c.request = req
	c.response = res
	c.internal = param.NewMap()
	c.store = param.NewSafeStore()
	c.path = ""
	c.pnames = nil
	c.hnames = nil
	c.hvalues = nil
	c.renderer = nil
	c.handler = NotFoundHandler
	c.route = nil
	c.rid = -1
	c.sessionOptions = nil
	c.withFormatExtension = false
	c.defaultExtension = ""
	c.format = ""
	c.code = 0
	c.auto = false
	c.preResponseHook = nil
	c.accept = nil
	c.dataEngine = NewData(c)
	c.onHostFound = c.echo.onHostFound
	c.renderDataWrapper = c.echo.renderDataWrapper
	c.ResetFuncs(c.echo.FuncMap)
	c.realIP = ""
	// NOTE: Don't reset because it has to have length c.echo.maxParam at all times
	for i := 0; i < *c.echo.maxParam; i++ {
		c.pvalues[i] = ""
	}
}

func (c *XContext) GetFunc(key string) interface{} {
	return c.funcs[key]
}

func (c *XContext) SetFunc(key string, val interface{}) {
	if ctxFunc, ok := val.(func(Context) interface{}); ok {
		val = ctxFunc(c)
	}
	c.funcs[key] = val
}

func (c *XContext) ResetFuncs(funcs map[string]interface{}) {
	c.funcs = map[string]interface{}{}
	for name, fn := range funcs {
		c.SetFunc(name, fn)
	}
}

func (c *XContext) Funcs() map[string]interface{} {
	return c.funcs
}

func (c *XContext) Renderer() Renderer {
	if c.renderer != nil {
		return c.renderer
	}
	return c.echo.renderer
}

func (c *XContext) getRenderData(data interface{}) interface{} {
	if data == nil {
		data = c.dataEngine.GetData()
		if c.renderDataWrapper == nil {
			return data
		}
		rdata := c.Internal().Get(`wrappedNilRenderData`)
		if rdata != nil {
			return rdata
		}
		data = c.renderDataWrapper(c, data)
		c.Internal().Set(`wrappedNilRenderData`, data)
		return data
	}
	if c.renderDataWrapper != nil {
		data = c.renderDataWrapper(c, data)
	}
	return data
}

func (c *XContext) Fetch(name string, data interface{}) (b []byte, err error) {
	name, err = c.echo.Template(c, name, data)
	if err != nil {
		return
	}
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := bufferpool.Get()
	defer bufferpool.Release(buf)
	data = c.getRenderData(data)
	err = c.renderer.Render(buf, name, data, c)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

func (c *XContext) Validate(item interface{}, args ...interface{}) error {
	return Validate(c, item, args...)
}

func (c *XContext) Validator() Validator {
	return c.validator
}

func (c *XContext) SetValidator(v Validator) {
	c.validator = v
}

// SetRenderer registers an HTML template renderer.
func (c *XContext) SetRenderer(r Renderer) {
	c.renderer = r
}

// SetRenderDataWrapper .
func (c *XContext) SetRenderDataWrapper(dataWrapper DataWrapper) {
	c.renderDataWrapper = dataWrapper
}

// RenderDataWrapper .
func (c *XContext) RenderDataWrapper() DataWrapper {
	return c.renderDataWrapper
}

func (c *XContext) SetSessioner(s Sessioner) {
	c.sessioner = s
}

func (c *XContext) Atop(v string) param.String {
	return param.String(v)
}

func (c *XContext) ToParamString(v string) param.String {
	return param.String(v)
}

func (c *XContext) ToStringSlice(v []string) param.StringSlice {
	return param.StringSlice(v)
}

func (c *XContext) SetFormat(format string) {
	c.format = format
}

func (c *XContext) WithFormatExtension(on bool) {
	c.withFormatExtension = on
}

func (c *XContext) SetCode(code int) {
	c.code = code
}

func (c *XContext) Code() int {
	return c.code
}

func (c *XContext) SetData(data Data) {
	c.dataEngine = data
}

func (c *XContext) Data() Data {
	return c.dataEngine
}

// MapData 映射数据到结构体
func (c *XContext) MapData(i interface{}, data map[string][]string, names ...string) error {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	return FormToStruct(c.echo, i, data, name)
}

func (c *XContext) AddPreResponseHook(hook func() error) Context {
	if c.preResponseHook == nil {
		c.preResponseHook = []func() error{hook}
	} else {
		c.preResponseHook = append(c.preResponseHook, hook)
	}
	return c
}

func (c *XContext) SetPreResponseHook(hook ...func() error) Context {
	c.preResponseHook = hook
	return c
}

func (c *XContext) OnHostFound(onHostFound func(Context) (bool, error)) Context {
	c.onHostFound = onHostFound
	return c
}

func (c *XContext) FireHostFound() (bool, error) {
	if c.onHostFound == nil {
		return true, nil
	}
	return c.onHostFound(c)
}

func (c *XContext) preResponse() error {
	if c.preResponseHook == nil {
		c.cookier.Send()
		return nil
	}
	for _, hook := range c.preResponseHook {
		if err := hook(); err != nil {
			return err
		}
	}
	c.cookier.Send()
	return nil
}

func (c *XContext) PrintFuncs() {
	for key, fn := range c.Funcs() {
		fmt.Printf("[Template Func](%p) %-15s -> %s \n", fn, key, HandlerName(fn))
	}
}

func (c *XContext) Dispatch(route string) Handler {
	u, err := url.Parse(route)
	if err != nil {
		return ErrorHandler(err)
	}
	c.Request().URL().SetRawQuery(u.RawQuery)
	for key, values := range u.Query() {
		for index, value := range values {
			if index == 0 {
				c.Request().URL().Query().Set(key, value)
			} else {
				c.Request().URL().Query().Add(key, value)
			}
		}
	}
	c.handler = NotFoundHandler
	return c.Echo().Router().Dispatch(c, u.Path)
}
