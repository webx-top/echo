package echo

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"github.com/admpub/events"

	pkgCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
)

// Context represents context for the current request. It holds request and
// response objects, path parameters, data and registered handler.
type Context interface {
	context.Context
	eventsEmitterer
	IRouteDispatchPath
	URLGenerator
	SetEmitterer(events.Emitterer)
	Emitterer() events.Emitterer
	Handler() Handler

	//Transaction
	SetTransaction(t Transaction)
	Transaction() Transaction
	Begin() error
	Rollback() error
	Commit() error
	End(succeed bool) error

	//Standard Context
	StdContext() context.Context
	WithContext(ctx context.Context) *http.Request
	SetValue(key string, value interface{})

	SetValidator(Validator)
	Validator() Validator
	Validate(item interface{}, args ...interface{}) error
	translator
	SetTranslator(Translator)
	Translator() Translator
	Request() engine.Request
	Response() engine.Response
	Handle(Context) error
	Logger() logger.Logger
	Object() *XContext
	Echo() *Echo
	Route() *Route
	Reset(engine.Request, engine.Response)
	Dispatch(route string) Handler

	//----------------
	// Param
	//----------------

	Path() string
	P(int, ...string) string
	Param(string, ...string) string
	// ParamNames returns path parameter names.
	ParamNames() []string
	ParamValues() []string
	SetParamNames(names ...string)
	SetParamValues(values ...string)
	// Host
	HostParamNames() []string
	HostParamValues() []string
	HostParam(string, ...string) string
	HostP(int, ...string) string
	SetHostParamNames(names ...string)
	SetHostParamValues(values ...string)

	// Queries returns the query parameters as map. It is an alias for `engine.URL#Query()`.
	Queries() map[string][]string
	QueryValues(string) []string
	QueryxValues(string) param.StringSlice
	QueryAnyValues(name string, other ...string) []string
	QueryAnyxValues(name string, other ...string) param.StringSlice
	Query(name string, defaults ...string) string
	QueryAny(name string, other ...string) string
	QueryLast(name string, defaults ...string) string
	QueryAnyLast(name string, other ...string) string

	//----------------
	// Form data
	//----------------

	Form(name string, defaults ...string) string
	FormAny(name string, other ...string) string
	FormLast(name string, defaults ...string) string
	FormAnyLast(name string, other ...string) string
	FormValues(string) []string
	FormxValues(string) param.StringSlice
	FormAnyValues(name string, other ...string) []string
	FormAnyxValues(name string, other ...string) param.StringSlice
	// Forms returns the form parameters as map. It is an alias for `engine.Request#Form().All()`.
	Forms() map[string][]string

	// Param+
	Px(int, ...string) param.String
	Paramx(string, ...string) param.String
	Queryx(name string, defaults ...string) param.String
	QueryLastx(name string, defaults ...string) param.String
	QueryAnyx(name string, other ...string) param.String
	QueryAnyLastx(name string, other ...string) param.String
	Formx(name string, defaults ...string) param.String
	FormLastx(name string, defaults ...string) param.String
	FormAnyx(name string, other ...string) param.String
	FormAnyLastx(name string, other ...string) param.String
	// string to param.String
	Atop(string) param.String
	ToParamString(string) param.String
	ToStringSlice([]string) param.StringSlice

	//----------------
	// Context data
	//----------------

	Set(string, interface{})
	Get(string, ...interface{}) interface{}
	Incr(key string, n interface{}, defaults ...interface{}) int64
	Decr(key string, n interface{}, defaults ...interface{}) int64
	Delete(...string)
	Stored() Store
	Internal() *param.SafeMap

	//----------------
	// Bind
	//----------------

	Bind(interface{}, ...FormDataFilter) error
	BindAndValidate(interface{}, ...FormDataFilter) error
	MustBind(interface{}, ...FormDataFilter) error
	MustBindAndValidate(interface{}, ...FormDataFilter) error
	BindWithDecoder(interface{}, BinderValueCustomDecoders, ...FormDataFilter) error
	BindAndValidateWithDecoder(interface{}, BinderValueCustomDecoders, ...FormDataFilter) error
	MustBindWithDecoder(interface{}, BinderValueCustomDecoders, ...FormDataFilter) error
	MustBindAndValidateWithDecoder(interface{}, BinderValueCustomDecoders, ...FormDataFilter) error

	//----------------
	// Response data
	//----------------

	Render(string, interface{}, ...int) error
	RenderBy(string, func(string) ([]byte, error), interface{}, ...int) ([]byte, error)
	HTML(string, ...int) error
	String(string, ...int) error
	Blob([]byte, ...int) error
	JSON(interface{}, ...int) error
	JSONBlob([]byte, ...int) error
	JSONP(string, interface{}, ...int) error
	XML(interface{}, ...int) error
	XMLBlob([]byte, ...int) error
	Stream(func(io.Writer) (bool, error)) error
	SSEvent(string, chan interface{}) error
	File(string, ...http.FileSystem) error
	CacheableFile(string, time.Duration, ...http.FileSystem) error
	Attachment(io.Reader, string, time.Time, ...bool) error
	CacheableAttachment(io.Reader, string, time.Time, time.Duration, ...bool) error
	NotModified() error
	NoContent(...int) error
	Redirect(string, ...int) error
	Error(err error)
	NewError(code pkgCode.Code, msg string, args ...interface{}) *Error
	NewErrorWith(err error, code pkgCode.Code, args ...interface{}) *Error
	SetCode(int)
	Code() int
	SetData(Data)
	Data() Data

	IsValidCache(modifiedAt time.Time) bool
	SetCacheHeader(modifiedAt time.Time, maxAge ...time.Duration)

	// ServeContent sends static content from `io.Reader` and handles caching
	// via `If-Modified-Since` request header. It automatically sets `Content-Type`
	// and `Last-Modified` response headers.
	ServeContent(io.Reader, string, time.Time, ...time.Duration) error
	ServeCallbackContent(func(Context) (io.Reader, error), string, time.Time, ...time.Duration) error

	//----------------
	// FuncMap
	//----------------

	SetFunc(string, interface{})
	GetFunc(string) interface{}
	ResetFuncs(map[string]interface{})
	Funcs() map[string]interface{}
	PrintFuncs()

	//----------------
	// Render
	//----------------
	SetAuto(on bool) Context
	Fetch(string, interface{}) ([]byte, error)
	SetRenderer(Renderer)
	SetRenderDataWrapper(DataWrapper)
	Renderer() Renderer
	RenderDataWrapper() DataWrapper

	//----------------
	// Cookie
	//----------------

	SetCookieOptions(*CookieOptions)
	CookieOptions() *CookieOptions
	NewCookie(string, string) *http.Cookie
	Cookie() Cookier
	GetCookie(string) string
	// SetCookie set cookie
	//  @param: key, value, maxAge(seconds), path(/), domain, secure(false), httpOnly(false), sameSite(lax/strict/default)
	SetCookie(string, string, ...interface{})

	//----------------
	// Session
	//----------------

	SetSessionOptions(*SessionOptions)
	SessionOptions() *SessionOptions
	SetSessioner(Sessioner)
	Session() Sessioner
	Flash(...string) interface{}

	//----------------
	// Request data
	//----------------

	Header(string) string
	IsAjax() bool
	IsPjax() bool
	PjaxContainer() string
	Method() string
	Format() string
	SetFormat(string)
	IsMethod(method string) bool
	IsPost() bool
	IsGet() bool
	IsPut() bool
	IsDel() bool
	IsHead() bool
	IsPatch() bool
	IsOptions() bool
	IsSecure() bool
	IsWebsocket() bool
	IsUpload() bool
	ResolveContentType() string
	WithFormatExtension(bool)
	SetDefaultExtension(string)
	DefaultExtension() string
	ResolveFormat() string
	Accept() *Accepts
	Protocol() string
	Site() string
	FullRequestURI() string
	RequestURI() string
	Scheme() string
	Domain() string
	Host() string
	Proxy() []string
	Referer() string
	Port() int
	RealIP() string
	HasAnyRequest() bool

	MapForm(i interface{}, names ...string) error
	MapData(i interface{}, data map[string][]string, names ...string) error
	SaveUploadedFile(fieldName string, saveAbsPath string, saveFileName ...func(*multipart.FileHeader) (string, error)) (*multipart.FileHeader, error)
	SaveUploadedFileToWriter(string, io.Writer) (*multipart.FileHeader, error)
	//Multiple file upload
	SaveUploadedFiles(fieldName string, savePath func(*multipart.FileHeader) (string, error)) error
	SaveUploadedFilesToWriter(fieldName string, writer func(*multipart.FileHeader) (io.Writer, error)) error

	//----------------
	// Hook
	//----------------

	AddPreResponseHook(func() error) Context
	SetPreResponseHook(...func() error) Context
	OnHostFound(func(Context) (bool, error)) Context
	FireHostFound() (bool, error)
	OnRelease(...func(Context)) Context
	FireRelease()
}

type eCtxKey struct{}

var contextKey eCtxKey

func FromStdContext(c context.Context) (Context, bool) {
	ctx, ok := c.Value(contextKey).(Context)
	return ctx, ok
}

func ToStdContext(ctx context.Context, eCtx Context) context.Context {
	return context.WithValue(ctx, contextKey, eCtx)
}

func AsStdContext(eCtx Context) context.Context {
	return ToStdContext(eCtx, eCtx)
}

var typeOfContext = reflect.TypeOf((*Context)(nil)).Elem()

func IsContext(t reflect.Type) bool {
	return t.Implements(typeOfContext)
}

type ContextReseter interface {
	Reset(req engine.Request, res engine.Response)
}

type Releaseable interface {
	Release()
}
