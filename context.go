package echo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
	"golang.org/x/net/context"
)

type (
	// Context represents context for the current request. It holds request and
	// response objects, path parameters, data and registered handler.
	Context interface {
		context.Context
		Translator
		Request() engine.Request
		Response() engine.Response
		Path() string
		P(int) string
		Param(string) string

		// ParamNames returns path parameter names.
		ParamNames() []string
		SetParamNames(...string)
		ParamValues() []string
		SetParamValues(values ...string)

		// Queries returns the query parameters as map. It is an alias for `engine.URL#Query()`.
		Queries() map[string][]string
		QueryValues(string) []string
		Query(string) string

		Form(string) string
		FormValues(string) []string

		// Forms returns the form parameters as map. It is an alias for `engine.Request#Form().All()`.
		Forms() map[string][]string

		Set(string, interface{})
		Get(string) interface{}
		Bind(interface{}) error
		MustBind(interface{}) error
		Render(string, interface{}, ...int) error
		HTML(string, ...int) error
		String(string, ...int) error
		JSON(interface{}, ...int) error
		JSONBlob([]byte, ...int) error
		JSONP(string, interface{}, ...int) error
		XML(interface{}, ...int) error
		XMLBlob([]byte, ...int) error
		File(string) error
		Attachment(io.ReadSeeker, string) error
		NoContent(...int) error
		Redirect(string, ...int) error
		Error(err error)
		Handle(Context) error
		Logger() logger.Logger
		Object() *xContext

		// ServeContent sends static content from `io.Reader` and handles caching
		// via `If-Modified-Since` request header. It automatically sets `Content-Type`
		// and `Last-Modified` response headers.
		ServeContent(io.ReadSeeker, string, time.Time) error

		SetFunc(string, interface{})
		GetFunc(string) interface{}
		ResetFuncs(map[string]interface{})
		Funcs() map[string]interface{}
		Reset(engine.Request, engine.Response)
		Fetch(string, interface{}) ([]byte, error)
		SetRenderer(Renderer)

		// Cookie
		SetCookieOptions(*CookieOptions)
		CookieOptions() *CookieOptions
		NewCookie(string, string) *Cookie
		Cookie() Cookier
		GetCookie(string) string
		SetCookie(string, string, ...interface{})

		SetSessioner(Sessioner)
		Session() Sessioner
		Flash(string) interface{}

		//with type action
		Px(int) param.String
		Paramx(string) param.String
		Queryx(string) param.String
		Formx(string) param.String
		//string to param.String
		Atop(string) param.String

		SetTranslator(Translator)
	}

	xContext struct {
		Translator
		sessioner     Sessioner
		cookier       Cookier
		context       context.Context
		request       engine.Request
		response      engine.Response
		path          string
		pnames        []string
		pvalues       []string
		store         store
		handler       Handler
		echo          *Echo
		funcs         map[string]interface{}
		renderer      Renderer
		cookieOptions *CookieOptions
	}

	store map[string]interface{}
)

// NewContext creates a Context object.
func NewContext(req engine.Request, res engine.Response, e *Echo) Context {
	c := &xContext{
		Translator: DefaultNopTranslate,
		context:    context.Background(),
		request:    req,
		response:   res,
		echo:       e,
		pvalues:    make([]string, *e.maxParam),
		store:      make(store),
		handler:    notFoundHandler,
		funcs:      make(map[string]interface{}),
		sessioner:  DefaultNopSession,
	}
	c.cookier = NewCookier(c)
	return c
}

func (c *xContext) Context() context.Context {
	return c.context
}

func (c *xContext) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *xContext) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

func (c *xContext) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *xContext) Err() error {
	return c.context.Err()
}

func (c *xContext) Value(key interface{}) interface{} {
	return c.context.Value(key)
}

func (c *xContext) Handle(ctx Context) error {
	return c.handler.Handle(ctx)
}

// Request returns *http.Request.
func (c *xContext) Request() engine.Request {
	return c.request
}

// Response returns *Response.
func (c *xContext) Response() engine.Response {
	return c.response
}

// Path returns the registered path for the handler.
func (c *xContext) Path() string {
	return c.path
}

// P returns path parameter by index.
func (c *xContext) P(i int) (value string) {
	l := len(c.pnames)
	if i < l {
		value = c.pvalues[i]
	}
	return
}

// Param returns path parameter by name.
func (c *xContext) Param(name string) (value string) {
	l := len(c.pnames)
	for i, n := range c.pnames {
		if n == name && i < l {
			value = c.pvalues[i]
			break
		}
	}
	return
}

func (c *xContext) ParamNames() []string {
	return c.pnames
}

func (c *xContext) SetParamNames(names ...string) {
	c.pnames = names
}

func (c *xContext) ParamValues() []string {
	return c.pvalues
}

func (c *xContext) SetParamValues(values ...string) {
	c.pvalues = values
}

// Query returns query parameter by name.
func (c *xContext) Query(name string) string {
	return c.request.URL().QueryValue(name)
}

func (c *xContext) QueryValues(name string) []string {
	return c.request.URL().QueryValues(name)
}

func (c *xContext) Queries() map[string][]string {
	return c.request.URL().Query()
}

// Form returns form parameter by name.
func (c *xContext) Form(name string) string {
	return c.request.FormValue(name)
}

func (c *xContext) FormValues(name string) []string {
	return c.request.Form().Gets(name)
}

func (c *xContext) Forms() map[string][]string {
	return c.request.Form().All()
}

// Get retrieves data from the context.
func (c *xContext) Get(key string) interface{} {
	return c.store[key]
}

// Set saves data in the context.
func (c *xContext) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(store)
	}
	c.store[key] = val
}

// Bind binds the request body into specified type `i`. The default binder does
// it based on Content-Type header.
func (c *xContext) Bind(i interface{}) error {
	return c.echo.binder.Bind(i, c)
}

func (c *xContext) MustBind(i interface{}) error {
	return c.echo.binder.MustBind(i, c)
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *xContext) Render(name string, data interface{}, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	b, err := c.Fetch(name, data)
	if err != nil {
		return
	}
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

// HTML sends an HTTP response with status code.
func (c *xContext) HTML(html string, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write([]byte(html))
	return
}

// String sends a string response with status code.
func (c *xContext) String(s string, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	c.response.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write([]byte(s))
	return
}

// JSON sends a JSON response with status code.
func (c *xContext) JSON(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = json.MarshalIndent(i, "", "  ")
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(b, codes...)
}

// JSONBlob sends a JSON blob response with status code.
func (c *xContext) JSONBlob(b []byte, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *xContext) JSONP(callback string, i interface{}, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)

	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.response.Write(b); err != nil {
		return
	}
	_, err = c.response.Write([]byte(");"))
	return
}

// XML sends an XML response with status code.
func (c *xContext) XML(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = xml.MarshalIndent(i, "", "  ")
	} else {
		b, err = xml.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.XMLBlob(b, codes...)
}

// XMLBlob sends a XML blob response with status code.
func (c *xContext) XMLBlob(b []byte, codes ...int) (err error) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(xml.Header)); err != nil {
		return
	}
	_, err = c.response.Write(b)
	return
}

func (c *xContext) File(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, "index.html")
		f, err = os.Open(file)
		if err != nil {
			return ErrNotFound
		}
		fi, _ = f.Stat()
	}
	return c.ServeContent(f, fi.Name(), fi.ModTime())
}

func (c *xContext) Attachment(r io.ReadSeeker, name string) (err error) {
	c.response.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	c.response.Header().Set(HeaderContentDisposition, "attachment; filename="+name)
	c.response.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.response, r)
	return
}

// NoContent sends a response with no body and a status code.
func (c *xContext) NoContent(codes ...int) error {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	c.response.WriteHeader(code)
	return nil
}

// Redirect redirects the request with status code.
func (c *xContext) Redirect(url string, codes ...int) error {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	c.response.Redirect(url, code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *xContext) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

// Logger returns the `Logger` instance.
func (c *xContext) Logger() logger.Logger {
	return c.echo.logger
}

// Object returns the `context` object.
func (c *xContext) Object() *xContext {
	return c
}

func (c *xContext) ServeContent(content io.ReadSeeker, name string, modtime time.Time) error {
	rq := c.Request()
	rs := c.Response()

	if t, err := time.Parse(http.TimeFormat, rq.Header().Get(HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
		rs.Header().Del(HeaderContentType)
		rs.Header().Del(HeaderContentLength)
		return c.NoContent(http.StatusNotModified)
	}

	rs.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	rs.Header().Set(HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
	rs.WriteHeader(http.StatusOK)
	_, err := io.Copy(rs, content)
	return err
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = MIMEOctetStream
	}
	return
}

// Echo returns the `Echo` instance.
func (c *xContext) Echo() *Echo {
	return c.echo
}

func (c *xContext) SetTranslator(t Translator) {
	c.Translator = t
}

func (c *xContext) Reset(req engine.Request, res engine.Response) {
	c.Translator = DefaultNopTranslate
	c.sessioner = DefaultNopSession
	c.context = context.Background()
	c.request = req
	c.response = res
	c.store = nil
	c.funcs = make(map[string]interface{})
	c.renderer = nil
	c.handler = notFoundHandler
	c.cookieOptions = nil

	c.SetFunc(`Lang`, c.Lang)
	c.SetFunc(`T`, c.T)
	c.SetFunc(`Cookie`, c.Cookie)
	c.SetFunc(`Session`, c.Session)
	c.SetFunc(`Query`, c.Query)
	c.SetFunc(`Form`, c.Form)
	c.SetFunc(`QueryValues`, c.QueryValues)
	c.SetFunc(`FormValues`, c.FormValues)
	c.SetFunc(`Param`, c.Param)
	c.SetFunc(`Atop`, c.Atop)
	c.SetFunc(`URL`, req.URL)
	c.SetFunc(`Header`, req.Header)
	c.SetFunc(`Flash`, c.Flash)
	for name, function := range DefaultFuncMap {
		c.SetFunc(name, function)
	}
	if c.echo.FuncMap != nil {
		for name, function := range c.echo.FuncMap {
			c.SetFunc(name, function)
		}
	}
}

func (c *xContext) GetFunc(key string) interface{} {
	return c.funcs[key]
}

func (c *xContext) SetFunc(key string, val interface{}) {
	c.funcs[key] = val
}

func (c *xContext) ResetFuncs(funcs map[string]interface{}) {
	c.funcs = funcs
}

func (c *xContext) Funcs() map[string]interface{} {
	return c.funcs
}

func (c *xContext) Fetch(name string, data interface{}) (b []byte, err error) {
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := new(bytes.Buffer)
	err = c.renderer.Render(buf, name, data, c)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// SetRenderer registers an HTML template renderer.
func (c *xContext) SetRenderer(r Renderer) {
	c.renderer = r
}

func (c *xContext) SetSessioner(s Sessioner) {
	c.sessioner = s
}

func (c *xContext) Session() Sessioner {
	return c.sessioner
}

func (c *xContext) Flash(name string) (r interface{}) {
	if v := c.sessioner.Flashes(name); len(v) > 0 {
		r = v[0]
	}
	return r
}

func (c *xContext) SetCookieOptions(opts *CookieOptions) {
	c.cookieOptions = opts
}

func (c *xContext) CookieOptions() *CookieOptions {
	if c.cookieOptions == nil {
		c.cookieOptions = &CookieOptions{}
	}
	return c.cookieOptions
}

func (c *xContext) NewCookie(key string, value string) *Cookie {
	return NewCookie(key, value, c.CookieOptions())
}

func (c *xContext) Cookie() Cookier {
	return c.cookier
}

func (c *xContext) GetCookie(key string) string {
	return c.cookier.Get(key)
}

func (c *xContext) SetCookie(key string, val string, args ...interface{}) {
	c.cookier.Set(key, val, args...)
}

func (c *xContext) Px(n int) param.String {
	return param.String(c.P(n))
}

func (c *xContext) Paramx(name string) param.String {
	return param.String(c.Param(name))
}

func (c *xContext) Queryx(name string) param.String {
	return param.String(c.Query(name))
}

func (c *xContext) Formx(name string) param.String {
	return param.String(c.Form(name))
}

func (c *xContext) Atop(v string) param.String {
	return param.String(v)
}
