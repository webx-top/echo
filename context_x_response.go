package echo

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

	json "github.com/admpub/xencoding/json/standard"
	xml "github.com/admpub/xencoding/xml/standard"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/poolx/bufferpool"
)

// Response returns *Response.
func (c *XContext) Response() engine.Response {
	return c.response
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *XContext) Render(name string, data interface{}, codes ...int) error {
	if c.auto {
		if ok, err := c.echo.AutoDetectRenderFormat(c, data); ok {
			return err
		}
	}
	c.dataEngine.SetTmplFuncs()
	b, err := c.Fetch(name, data)
	if err != nil {
		return err
	}
	b = bytes.TrimLeftFunc(b, unicode.IsSpace)
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	return c.Blob(b, codes...)
}

func (c *XContext) RenderBy(name string, content func(string) ([]byte, error), data interface{}, codes ...int) (b []byte, err error) {
	name, err = c.echo.Template(c, name, data)
	if err != nil {
		return
	}
	c.dataEngine.SetTmplFuncs()
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := bufferpool.Get()
	defer bufferpool.Release(buf)
	data = c.getRenderData(data)
	err = c.renderer.RenderBy(buf, name, content, data, c)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// HTML sends an HTTP response with status code.
func (c *XContext) HTML(html string, codes ...int) error {
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	return c.Blob(engine.Str2bytes(html), codes...)
}

// String sends a string response with status code.
func (c *XContext) String(s string, codes ...int) error {
	c.response.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	return c.Blob(engine.Str2bytes(s), codes...)
}

func (c *XContext) Blob(b []byte, codes ...int) (err error) {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	err = c.preResponse()
	if err != nil {
		return
	}
	c.response.WriteHeader(c.code)
	_, err = c.response.Write(b)
	return
}

// JSON sends a JSON response with status code.
func (c *XContext) JSON(i interface{}, codes ...int) (err error) {
	if m, y := i.(JSONModifer); y {
		i, err = m.JSON(c)
		if err != nil {
			return err
		}
	}
	var b []byte
	if ft, ok := c.Route().Get(metaKeyEncodingConfig).(EncodingConfig); ok {
		b, err = json.MarshalWithOption(
			i,
			json.OptionFilter(ft.filter),
			json.OptionSelector(ft.selector),
		)
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(b, codes...)
}

// JSONBlob sends a JSON blob response with status code.
func (c *XContext) JSONBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	err = c.Blob(b, codes...)
	return
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *XContext) JSONP(callback string, i interface{}, codes ...int) (err error) {
	if m, y := i.(JSONModifer); y {
		i, err = m.JSON(c)
		if err != nil {
			return err
		}
	}
	var b []byte
	if ft, ok := c.Route().Get(metaKeyEncodingConfig).(EncodingConfig); ok {
		b, err = json.MarshalWithOption(
			i,
			json.OptionFilter(ft.filter),
			json.OptionSelector(ft.selector),
		)
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	b = []byte(callback + "(" + string(b) + ");")
	err = c.Blob(b, codes...)
	return
}

// XML sends an XML response with status code.
func (c *XContext) XML(i interface{}, codes ...int) (err error) {
	if m, y := i.(XMLModifer); y {
		i, err = m.XML(c)
		if err != nil {
			return err
		}
	}
	var b []byte
	if ft, ok := c.Route().Get(metaKeyEncodingConfig).(EncodingConfig); ok {
		b, err = xml.MarshalWithOption(
			i,
			xml.OptionFilter(ft.filter),
			xml.OptionSelector(ft.selector),
		)
	} else {
		b, err = xml.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.XMLBlob(b, codes...)
}

// XMLBlob sends a XML blob response with status code.
func (c *XContext) XMLBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	b = []byte(xml.Header + string(b))
	err = c.Blob(b, codes...)
	return
}

func (c *XContext) Stream(step func(w io.Writer) (bool, error)) error {
	return c.response.Stream(step)
}

func (c *XContext) SSEvent(event string, data chan interface{}) error {
	hdr := c.response.Header()
	hdr.Set(HeaderContentType, MIMEEventStream)
	hdr.Set(HeaderCacheControl, `no-cache`)
	hdr.Set(HeaderConnection, `keep-alive`)
	hdr.Set(HeaderTransferEncoding, `chunked`)
	return c.Stream(func(w io.Writer) (bool, error) {
		recv, ok := <-data
		if !ok {
			return ok, nil
		}
		b, err := c.Fetch(event, recv)
		if err != nil {
			return false, err
		}
		//c.Logger().Debugf(`SSEvent: %s`, b)
		_, err = w.Write(b)
		if err != nil {
			return false, err
		}
		return true, err
	})
}

func (c *XContext) Attachment(r io.Reader, name string, modtime time.Time, inline ...bool) error {
	SetAttachmentHeader(c, name, true, inline...)
	return c.ServeContent(r, name, modtime)
}

func (c *XContext) CacheableAttachment(r io.Reader, name string, modtime time.Time, maxAge time.Duration, inline ...bool) error {
	SetAttachmentHeader(c, name, true, inline...)
	return c.ServeContent(r, name, modtime, maxAge)
}

func (c *XContext) File(file string, fs ...http.FileSystem) error {
	return c.CacheableFile(file, 0, fs...)
}

func (c *XContext) CacheableFile(file string, maxAge time.Duration, fs ...http.FileSystem) (err error) {
	var f http.File
	customFS := len(fs) > 0 && fs[0] != nil
	if customFS {
		f, err = fs[0].Open(file)
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return ErrNotFound
	}
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		f.Close()
		file = filepath.Join(file, "index.html")
		if customFS {
			f, err = fs[0].Open(file)
		} else {
			f, err = os.Open(file)
		}
		if err != nil {
			return ErrNotFound
		}
		fi, err = f.Stat()
		if err != nil {
			return err
		}
	}
	if maxAge > time.Second {
		if c.IsValidCache(fi.ModTime()) {
			return c.NotModified()
		}
		c.SetCacheHeader(fi.ModTime(), maxAge)
	}
	c.Response().ServeContent(f, fi.Name(), fi.ModTime())
	return nil
}

func (c *XContext) ServeContent(content io.Reader, name string, modtime time.Time, cacheMaxAge ...time.Duration) error {
	if readSeeker, ok := content.(io.ReadSeeker); ok {
		if c.IsValidCache(modtime) {
			return c.NotModified()
		}
		c.SetCacheHeader(modtime, cacheMaxAge...)
		c.Response().ServeContent(readSeeker, name, modtime)
		return nil
	}
	return c.ServeCallbackContent(func(_ Context) (io.Reader, error) {
		return content, nil
	}, name, modtime)
}

func (c *XContext) IsValidCache(modifiedAt time.Time) bool {
	t, err := time.Parse(http.TimeFormat, c.Request().Header().Get(HeaderIfModifiedSince))
	return err == nil && modifiedAt.Before(t.Add(1*time.Second))
}

func (c *XContext) SetCacheHeader(modifiedAt time.Time, maxAge ...time.Duration) {
	rs := c.Response()
	hdr := rs.Header()
	if len(maxAge) > 0 && maxAge[0] > time.Second {
		now := time.Now().UTC()
		expiredAt := now.Add(maxAge[0])
		hdr.Set(HeaderExpires, expiredAt.Format(http.TimeFormat))
		hdr.Set(HeaderCacheControl, CacheControlPrefix+strconv.Itoa(int(maxAge[0].Seconds())))
	}
	hdr.Set(HeaderLastModified, modifiedAt.UTC().Format(http.TimeFormat))
}

func (c *XContext) NotModified() error {
	rs := c.Response()
	rs.Header().Del(HeaderContentType)
	rs.Header().Del(HeaderContentLength)
	return c.NoContent(http.StatusNotModified)
}

func (c *XContext) ServeCallbackContent(callback func(Context) (io.Reader, error), name string, modtime time.Time, cacheMaxAge ...time.Duration) error {
	if c.IsValidCache(modtime) {
		return c.NotModified()
	}
	content, err := callback(c)
	if err != nil {
		return err
	}
	if readSeeker, ok := content.(io.ReadSeeker); ok {
		c.SetCacheHeader(modtime, cacheMaxAge...)
		c.Response().ServeContent(readSeeker, name, modtime)
		return nil
	}
	rs := c.Response()
	rs.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	c.SetCacheHeader(modtime, cacheMaxAge...)
	rs.WriteHeader(http.StatusOK)
	rs.KeepBody(false)
	_, err = io.Copy(rs, content)
	return err
}

// NoContent sends a response with no body and a status code.
func (c *XContext) NoContent(codes ...int) error {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	c.response.WriteHeader(c.code)
	return nil
}

// Redirect redirects the request with status code.
func (c *XContext) Redirect(url string, codes ...int) error {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	err := c.preResponse()
	if err != nil {
		return err
	}
	format := c.Format()
	if format != `html` && c.auto {
		if render, ok := c.echo.formatRenderers[format]; ok && render != nil {
			if c.dataEngine.GetData() == nil {
				c.dataEngine.SetData(c.Stored(), c.dataEngine.GetCode().Int())
			}
			c.dataEngine.SetURL(url)
			return render(c, c.dataEngine.GetData())
		}
	}
	c.response.Redirect(url, code)
	return nil
}
