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
package render

import (
	_ "embed"
	"net/http"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var (
	DefaultOptions = &Options{
		Skipper:              echo.DefaultSkipper,
		ErrorPages:           make(map[int]string),
		DefaultHTTPErrorCode: http.StatusInternalServerError,
		SetFuncMap:           []echo.HandlerFunc{},
		DefaultRenderer:      defaultRender,
	}
)

type ErrorProcessor func(ctx echo.Context, err error) (processed bool, newErr error)

type Options struct {
	Skipper              echo.Skipper
	ErrorPages           map[int]string
	ErrorProcessors      []ErrorProcessor
	ErrorCodeLinks       map[code.Code]echo.KVList
	DefaultHTTPErrorCode int
	SetFuncMap           []echo.HandlerFunc
	DefaultRenderer      func(c echo.Context, data echo.H, code int) ([]byte, error)
}

func (opt *Options) AddFuncSetter(set ...echo.HandlerFunc) *Options {
	if opt.SetFuncMap == nil {
		opt.SetFuncMap = make([]echo.HandlerFunc, len(DefaultOptions.SetFuncMap))
		copy(opt.SetFuncMap, DefaultOptions.SetFuncMap)
	}
	opt.SetFuncMap = append(opt.SetFuncMap, set...)
	return opt
}

func (opt *Options) SetFuncSetter(set ...echo.HandlerFunc) *Options {
	opt.SetFuncMap = set
	return opt
}

func (opt *Options) AddErrorProcessor(h ...ErrorProcessor) *Options {
	opt.ErrorProcessors = append(opt.ErrorProcessors, h...)
	return opt
}

func (opt *Options) SetErrorProcessor(h ...ErrorProcessor) *Options {
	opt.ErrorProcessors = h
	return opt
}

func (opt *Options) SetDefaultRender(renderer func(c echo.Context, data echo.H, code int) ([]byte, error)) *Options {
	opt.DefaultRenderer = renderer
	return opt
}

func (opt *Options) GenTmplGetter() func(code int) string {
	tmplNum := len(opt.ErrorPages)
	if tmplNum == 0 {
		return func(code int) string {
			return ""
		}
	}
	return func(code int) string {
		tmpl, ok := opt.ErrorPages[code]
		if ok {
			return tmpl
		}
		if code != 0 {
			tmpl = opt.ErrorPages[0]
		}
		return tmpl
	}
}

// Middleware set renderer
func Middleware(d echo.Renderer) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetRenderer(d)
			return h.Handle(c)
		})
	}
}

func Auto() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetAuto(true)
			return h.Handle(c)
		})
	}
}

func setAndGetErrorMessage(c echo.Context, debugMessage string, prodMessage string) string {
	if c.Echo().Debug() {
		c.Data().SetInfo(debugMessage, 0)
		return debugMessage
	}
	c.Data().SetInfo(prodMessage, 0)
	return prodMessage
}

//go:embed error.tpl.html
var errorHTML []byte

func ErrorHTMLTemplate(_ string) ([]byte, error) {
	return errorHTML, nil
}

func defaultRender(c echo.Context, data echo.H, code int) ([]byte, error) {
	return c.RenderBy(`error.tpl`, ErrorHTMLTemplate, data, code)
}

func HTTPErrorHandler(opt *Options) echo.HTTPErrorHandler {
	if opt == nil {
		opt = DefaultOptions
	} else {
		if opt.Skipper == nil {
			opt.Skipper = DefaultOptions.Skipper
		}
		if opt.ErrorPages == nil {
			opt.ErrorPages = DefaultOptions.ErrorPages
		}
		if opt.DefaultHTTPErrorCode < 1 {
			opt.DefaultHTTPErrorCode = DefaultOptions.DefaultHTTPErrorCode
		}
		if opt.SetFuncMap == nil {
			opt.SetFuncMap = DefaultOptions.SetFuncMap
		}
		if opt.DefaultRenderer == nil {
			opt.DefaultRenderer = DefaultOptions.DefaultRenderer
		}
	}
	getTmpl := opt.GenTmplGetter()
	return func(err error, c echo.Context) {
		if opt.Skipper(c) {
			return
		}
		if err != nil {
			defer c.Logger().Debug(err, `: `, c.Request().URL().String())
		}
		if c.Response().Committed() {
			return
		}
		code := http.StatusInternalServerError
		data := c.Data().Reset()
		title := http.StatusText(code)
		var (
			panicErr  *echo.PanicError
			msg       string
			processed bool
			tmpl      string
		)
		for _, processor := range opt.ErrorProcessors {
			if processor == nil {
				continue
			}
			processed, err = processor(c, err)
			if processed {
				break
			}
		}
		var links echo.KVList
		if v, y := c.Get(`links`).(echo.KVList); y {
			links = v
		}
		switch e := err.(type) {
		case *echo.HTTPError:
			if e.Code > 0 {
				code = e.Code
			}
			msg = e.Message
			title = com.TextLine(msg)
			data.SetError(e)
		case *echo.PanicError:
			panicErr = e
			msg = setAndGetErrorMessage(c, e.Error(), title)
		case *echo.Error:
			code = e.Code.HTTPCode()
			msg = e.Message
			title = com.TextLine(msg)
			if opt.ErrorCodeLinks != nil {
				if v, y := opt.ErrorCodeLinks[e.Code]; y {
					links = append(links, v...)
				}
			}
			data.SetError(e)
		default:
			msg = e.Error()
			title = com.TextLine(msg)
			data.SetError(e)
		}
		c.SetCode(code)
		if c.Request().Method() == echo.HEAD {
			c.NoContent(code)
			return
		}
		dt := echo.H{
			"title":   title,
			"content": msg,
			"debug":   c.Echo().Debug(),
			"code":    code,
			"panic":   panicErr,
			"links":   links,
		}
		data.SetData(dt, data.GetCode().Int())
		var val interface{}
		switch c.Format() {
		case echo.ContentTypeText:
			val = msg
		case echo.ContentTypeJSON, echo.ContentTypeJSONP, echo.ContentTypeXML:
			val = data.GetData()
		case echo.ContentTypeHTML:
			fallthrough
		default:
			tmpl = getTmpl(code)
			for _, setFunc := range opt.SetFuncMap {
				err = setFunc(c)
				if err != nil {
					c.String(err.Error())
					return
				}
			}
			val = data.GetData()
			if len(tmpl) == 0 || echo.IsEmptyRoute(c.Route()) {
				b, renderErr := opt.DefaultRenderer(c, dt, code)
				if renderErr != nil {
					c.String(msg+"\n"+renderErr.Error(), code)
					return
				}
				c.Blob(b, code)
				return
			}
		}
		c.SetAuto(true)
		if renderErr := c.Render(tmpl, val); renderErr != nil {
			c.String(msg+"\n"+renderErr.Error(), code)
		}
	}
}
