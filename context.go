package echo

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"

	"bytes"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type (
	// Context represents context for the current request. It holds request and
	// response objects, path parameters, data and registered handler.
	Context interface {
		context.Context
		Request() *http.Request
		Response() *Response
		Socket() *websocket.Conn
		Path() string
		P(int) string
		Param(string) string
		Query(string) string
		Form(string) string
		Set(string, interface{})
		Get(string) interface{}
		Bind(interface{}) error
		Render(int, string, interface{}) error
		HTML(int, string) error
		String(int, string) error
		JSON(int, interface{}) error
		JSONIndent(int, interface{}, string, string) error
		JSONP(int, string, interface{}) error
		XML(int, interface{}) error
		XMLIndent(int, interface{}, string, string) error
		File(string, string, bool) error
		NoContent(int) error
		Redirect(int, string) error
		Error(error)
		SetFunc(string, interface{})
		GetFunc(string) interface{}
		Funcs() template.FuncMap
		X() *xContext
	}

	xContext struct {
		context.Context
		request  *http.Request
		response *Response
		socket   *websocket.Conn
		path     string
		pnames   []string
		pvalues  []string
		query    url.Values
		store    store
		echo     *Echo
		funcs    template.FuncMap
	}
	store map[string]interface{}
)

// NewContext creates a Context object.
func NewContext(req *http.Request, res *Response, e *Echo) Context {
	return &xContext{
		request:  req,
		response: res,
		echo:     e,
		pvalues:  make([]string, *e.maxParam),
		store:    make(store),
		funcs:    make(template.FuncMap),
	}
}

// Request returns *http.Request.
func (c *xContext) Request() *http.Request {
	return c.request
}

// Response returns *Response.
func (c *xContext) Response() *Response {
	return c.response
}

// Socket returns *websocket.Conn.
func (c *xContext) Socket() *websocket.Conn {
	return c.socket
}

func (c *xContext) SetSocket(socket *websocket.Conn) {
	c.socket = socket
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

// Query returns query parameter by name.
func (c *xContext) Query(name string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query.Get(name)
}

// Form returns form parameter by name.
func (c *xContext) Form(name string) string {
	return c.request.FormValue(name)
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

func (c *xContext) GetFunc(key string) interface{} {
	return c.funcs[key]
}

func (c *xContext) SetFunc(key string, val interface{}) {
	if c.funcs == nil {
		c.funcs = make(template.FuncMap)
	}
	c.funcs[key] = val
}

func (c *xContext) Funcs() template.FuncMap {
	return c.funcs
}

func (c *xContext) X() *xContext {
	return c
}

// Bind binds the request body into specified type `i`. The default binder does
// it based on Content-Type header.
func (c *xContext) Bind(i interface{}) error {
	return c.echo.binder.Bind(c.request, i)
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *xContext) Render(code int, name string, data interface{}) (err error) {
	b, err := c.Fetch(name, data)
	if err != nil {
		return
	}
	c.response.Header().Set(ContentType, TextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write(b)
	return
}

func (c *xContext) Fetch(name string, data interface{}) (b []byte, err error) {
	if c.echo.renderer == nil {
		return nil, RendererNotRegistered
	}
	buf := new(bytes.Buffer)
	err = c.echo.renderer.Render(buf, name, data, c.funcs)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// HTML sends an HTTP response with status code.
func (c *xContext) HTML(code int, html string) (err error) {
	c.response.Header().Set(ContentType, TextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(html))
	return
}

// String sends a string response with status code.
func (c *xContext) String(code int, s string) (err error) {
	c.response.Header().Set(ContentType, TextPlainCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(s))
	return
}

// JSON sends a JSON response with status code.
func (c *xContext) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.Json(code, b)
	return
}

// JSONIndent sends a JSON response with status code, but it applies prefix and indent to format the output.
func (c *xContext) JSONIndent(code int, i interface{}, prefix string, indent string) (err error) {
	b, err := json.MarshalIndent(i, prefix, indent)
	if err != nil {
		return err
	}
	c.Json(code, b)
	return
}

func (c *xContext) Json(code int, b []byte) {
	c.response.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write(b)
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *xContext) JSONP(code int, callback string, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.Jsonp(code, callback, b)
	return
}

func (c *xContext) Jsonp(code int, callback string, b []byte) {
	c.response.Header().Set(ContentType, ApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(callback + "("))
	c.response.Write(b)
	c.response.Write([]byte(");"))
}

// XML sends an XML response with status code.
func (c *xContext) XML(code int, i interface{}) (err error) {
	b, err := xml.Marshal(i)
	if err != nil {
		return err
	}
	c.Xml(code, b)
	return
}

// XMLIndent sends an XML response with status code, but it applies prefix and indent to format the output.
func (c *xContext) XMLIndent(code int, i interface{}, prefix string, indent string) (err error) {
	b, err := xml.MarshalIndent(i, prefix, indent)
	if err != nil {
		return err
	}
	c.Xml(code, b)
	return
}

func (c *xContext) Xml(code int, b []byte) {
	c.response.Header().Set(ContentType, ApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(xml.Header))
	c.response.Write(b)
}

// File sends a response with the content of the file. If `attachment` is set
// to true, the client is prompted to save the file with provided `name`,
// name can be empty, in that case name of the file is used.
func (c *xContext) File(path, name string, attachment bool) (err error) {
	dir, file := filepath.Split(path)
	if attachment {
		c.response.Header().Set(ContentDisposition, "attachment; filename="+name)
	}
	if err = c.echo.serveFile(dir, file, c); err != nil {
		c.response.Header().Del(ContentDisposition)
	}
	return
}

// NoContent sends a response with no body and a status code.
func (c *xContext) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// Redirect redirects the request using http.Redirect with status code.
func (c *xContext) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return InvalidRedirectCode
	}
	http.Redirect(c.response, c.request, url, code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *xContext) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

// Echo returns the `Echo` instance.
func (c *xContext) Echo() *Echo {
	return c.echo
}

func (c *xContext) reset(r *http.Request, w http.ResponseWriter, e *Echo) {
	c.request = r
	c.response.reset(w, e)
	c.query = nil
	c.store = nil
	c.echo = e
	c.funcs = nil
}
