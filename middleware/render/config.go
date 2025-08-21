package render

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type Config struct {
	// TmplDir is the directory where templates are stored.
	TmplDir string
	// Theme is the name of the theme to be used for rendering templates.
	Theme string
	// Engine is the name of the template engine to be used for rendering templates.
	Engine string
	// Style is the name of the style to be used for rendering templates.
	Style string
	// Reload indicates whether the templates should be reloaded on each modify.
	Reload bool
	// ParseStrings is a map of strings to be replaced in the template content.
	ParseStrings map[string]string
	// ParseStringFuncs is a map of functions that return strings to be replaced in the template content.
	ParseStringFuncs map[string]func() string
	StaticOptions    *middleware.StaticOptions
	Debug            bool
	renderer         driver.Driver
	FuncMapGlobal    map[string]interface{}
	RendererDo       []func(driver.Driver)
	CustomParser     func(tmpl string, content []byte) []byte

	// - HTTPErrorHandler -

	// ErrorPages defines the error pages to be used for specific HTTP error codes.
	ErrorPages map[int]string
	// DefaultHTTPErrorCode is the default HTTP error code to use when no specific error code is provided.
	DefaultHTTPErrorCode int
	// ErrorProcessors defines a list of error processors that can be used to handle errors.
	ErrorProcessors []ErrorProcessor
	// ErrorCodeLinks defines a mapping of error codes to links that can be used to provide additional information about the error.
	ErrorCodeLinks map[code.Code]echo.KVList
	// DefaultRenderer is a function that can be used to render the default response for an error.
	DefaultRenderer func(c echo.Context, data echo.H, code int) ([]byte, error)
	// UsingDefaultRenderer is a function that can be used to determine if the default renderer should be used.
	UsingDefaultRenderer func(echo.Context) bool
	errorPageFuncSetter  []echo.HandlerFunc
}

var DefaultFuncMapSkipper = func(c echo.Context) bool {
	return c.Format() != echo.ContentTypeHTML && !c.IsAjax() && !c.IsPjax()
}

func (t *Config) SetRendererDo(rd ...func(driver.Driver)) *Config {
	t.RendererDo = rd
	return t
}

func (t *Config) AddRendererDo(rd ...func(driver.Driver)) *Config {
	if t.RendererDo == nil {
		t.RendererDo = []func(driver.Driver){}
	}
	t.RendererDo = append(t.RendererDo, rd...)
	return t
}

func (t *Config) Parser() func(tmpl string, content []byte) []byte {
	if t.ParseStrings == nil {
		return t.CustomParser
	}
	var replaces []string
	for oldVal, newVal := range t.ParseStrings {
		replaces = append(replaces, oldVal, newVal)
	}
	if t.ParseStringFuncs != nil {
		for oldVal, newVal := range t.ParseStringFuncs {
			replaces = append(replaces, oldVal, newVal())
		}
	}
	if len(replaces) == 0 {
		return t.CustomParser
	}
	repl := strings.NewReplacer(replaces...)
	return func(tmpl string, content []byte) []byte {
		s := engine.Bytes2str(content)
		s = repl.Replace(s)
		content = engine.Str2bytes(s)
		if t.CustomParser != nil {
			content = t.CustomParser(tmpl, content)
		}
		return content
	}
}

// NewRenderer 新建渲染接口
func (t *Config) NewRenderer(manager ...driver.Manager) driver.Driver {
	tmplDir := t.TmplDir
	if len(t.Theme) > 0 {
		tmplDir = filepath.Join(tmplDir, t.Theme)
	}
	renderer := New(t.Engine, tmplDir)
	if len(manager) > 0 && manager[0] != nil {
		renderer.SetManager(manager[0])
	}
	if t.RendererDo != nil {
		for _, rendererDo := range t.RendererDo {
			rendererDo(renderer)
		}
	}
	renderer.Init()
	renderer.SetContentProcessor(t.Parser())
	return renderer
}

func (t *Config) AddFuncSetter(set ...echo.HandlerFunc) *Config {
	if t.errorPageFuncSetter == nil {
		t.errorPageFuncSetter = make([]echo.HandlerFunc, len(DefaultOptions.SetFuncMap))
		copy(t.errorPageFuncSetter, DefaultOptions.SetFuncMap)
	}
	t.errorPageFuncSetter = append(t.errorPageFuncSetter, set...)
	return t
}

func (t *Config) SetFuncSetter(set ...echo.HandlerFunc) *Config {
	t.errorPageFuncSetter = set
	return t
}

func (t *Config) HTTPErrorHandler() echo.HTTPErrorHandler {
	opt := &Options{
		ErrorPages:           t.ErrorPages,
		ErrorProcessors:      t.ErrorProcessors,
		ErrorCodeLinks:       t.ErrorCodeLinks,
		DefaultHTTPErrorCode: t.DefaultHTTPErrorCode,
		DefaultRenderer:      t.DefaultRenderer,
		UsingDefaultRenderer: t.UsingDefaultRenderer,
	}
	opt.SetFuncSetter(t.errorPageFuncSetter...)
	return HTTPErrorHandler(opt)
}

func (t *Config) StaticMiddleware() interface{} {
	if t.StaticOptions != nil {
		return middleware.Static(t.StaticOptions)
	}
	return nil
}

func (t *Config) ApplyTo(e *echo.Echo, manager ...driver.Manager) *Config {
	if t.renderer != nil {
		t.renderer.Close()
	}
	e.SetHTTPErrorHandler(t.HTTPErrorHandler())
	staticMW := t.StaticMiddleware()
	if staticMW != nil {
		e.Use(staticMW)
	}
	renderer := t.MakeRenderer(manager...)
	e.SetRenderer(renderer)
	return t
}

func defaultTplFuncMap() map[string]interface{} {
	return tplfunc.TplFuncMap
}

func (t *Config) MakeRenderer(manager ...driver.Manager) driver.Driver {
	renderer := t.NewRenderer(manager...)
	if t.FuncMapGlobal == nil {
		renderer.SetFuncMap(defaultTplFuncMap)
	} else {
		renderer.SetFuncMap(func() map[string]interface{} { return t.FuncMapGlobal })
	}
	t.renderer = renderer
	return renderer
}

func (t *Config) Renderer() driver.Driver {
	return t.renderer
}

// ThemeDir 主题所在文件夹的路径
func (t *Config) ThemeDir(args ...string) string {
	if len(args) < 1 {
		return filepath.Join(t.TmplDir, t.Theme)
	}
	return filepath.Join(t.TmplDir, args[0])
}
