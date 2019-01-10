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

package host

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	codec "github.com/admpub/securecookie"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/pprof"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

// New 创建Hosts实例，例如：New(nil)
func New(newContext func(*echo.Echo) echo.Context) (s *Hosts) {
	s = &Hosts{
		mapHosts:      make(map[string]*Host),
		mapNames:      make(map[string]*Host),
		URL:           `/`,
		MaxUploadSize: 10 * 1024 * 1024,
		RootName:      `base`,
		URLConvert:    LowerCaseFirst,
		URLRecovery:   UpperCaseFirst,
		SessionOptions: &echo.SessionOptions{
			Engine: `cookie`,
			Name:   `GOSID`,
			CookieOptions: &echo.CookieOptions{
				Prefix:   `e_`,
				HttpOnly: true,
				Path:     `/`,
			},
		},
	}
	s.NewContext = newContext //[1]
	if s.NewContext == nil {
		s.Core = echo.New()
	} else {
		s.Core = echo.NewWithContext(s.NewContext)
	}

	//Core middleware
	s.FuncMap = s.DefaultFuncMap()
	s.DefaultMiddlewares = []interface{}{
		mw.Log(),
		mw.Recover(),
		mw.FuncMap(s.FuncMap, func(ctx echo.Context) bool {
			return ctx.Format() != `html`
		}),
	}
	s.Core.Use(s.DefaultMiddlewares...)

	s.SetHTTPErrorOptions(nil)
	return
}

type (
	//URLConvert 网址转换
	URLConvert func(string) string

	//URLRecovery 网址还原
	URLRecovery func(string) string
)

var (
	//SnakeCase 单词全部小写并用下划线连接
	SnakeCase URLConvert = com.SnakeCase

	//LowerCaseFirst 小写首字母
	LowerCaseFirst URLConvert = com.LowerCaseFirst

	//PascalCase 帕斯卡命名法
	PascalCase URLRecovery = com.PascalCase

	//UpperCaseFirst 大写首字母
	UpperCaseFirst URLRecovery = strings.Title
)

// Hosts 定义应用
type Hosts struct {
	Core                  *echo.Echo
	Name                  string
	URLConvert            URLConvert  `json:"-" xml:"-"`
	URLRecovery           URLRecovery `json:"-" xml:"-"`
	MaxUploadSize         int64
	RootName              string
	URL                   string
	DefaultMiddlewares    []interface{} `json:"-" xml:"-"`
	DefaultPreMiddlewares []interface{} `json:"-" xml:"-"`
	SessionOptions        *echo.SessionOptions
	Renderer              driver.Driver                 `json:"-" xml:"-"`
	FuncMap               map[string]interface{}        `json:"-" xml:"-"`
	NewContext            func(*echo.Echo) echo.Context `json:"-" xml:"-"`
	Codec                 codec.Codec                   `json:"-" xml:"-"`
	mapHosts              map[string]*Host              //域名关联
	mapNames              map[string]*Host              //名称关联
	rootDir               string
	mutex                 sync.RWMutex
	before                []func() error
	after                 []func() error
}

// Before 服务启动之前执行钩子
func (s *Hosts) Before(f ...func() error) *Hosts {
	s.before = append(s.before, f...)
	return s
}

// After 服务结束之后执行钩子
func (s *Hosts) After(f ...func() error) *Hosts {
	s.after = append(s.after, f...)
	return s
}

// Pre 全局前置中间件
func (s *Hosts) Pre(middleware ...interface{}) {
	if len(s.mapNames) > 0 {
		panic(`The global pre-middleware must be set before Module is not created`)
	}
	if len(middleware) == 1 && middleware[0] == nil {
		s.DefaultPreMiddlewares = []interface{}{}
		s.Core.Clear()
		s.Core.Use(s.DefaultMiddlewares...)
		return
	}
	s.Core.Pre(middleware...)
	s.DefaultPreMiddlewares = append(middleware, s.DefaultPreMiddlewares...)
}

// Use 全局中间件
func (s *Hosts) Use(middleware ...interface{}) {
	if len(s.mapNames) > 0 {
		panic(`The global middleware must be set before Module is not created`)
	}
	if len(middleware) == 1 && middleware[0] == nil {
		s.DefaultMiddlewares = []interface{}{}
		s.Core.Clear()
		s.Core.Pre(s.DefaultPreMiddlewares...)
		return
	}
	s.Core.Use(middleware...)
	s.DefaultMiddlewares = append(s.DefaultMiddlewares, middleware...)
}

// FindByDomain 根据域名查询对应模块实例
func (s *Hosts) FindByDomain(host string) (*Host, bool) {
	module, has := s.mapHosts[host]
	if !has {
		if p := strings.LastIndexByte(host, ':'); p > -1 {
			module, has = s.mapHosts[host[0:p]]
		}
	}
	return module, has
}

// ServeHTTP HTTP服务执行入口
func (s *Hosts) ServeHTTP(r engine.Request, w engine.Response) {
	var h *echo.Echo
	domain := r.Host()
	host, has := s.FindByDomain(domain)
	if !has || host.Handler == nil {
		h = s.Core
	} else {
		h = host.Handler
	}

	if h != nil {
		h.ServeHTTP(r, w)
	} else {
		w.NotFound()
	}
}

// SetHTTPErrorOptions 设置错误处理选项
func (s *Hosts) SetHTTPErrorOptions(options *render.Options) *Hosts {
	s.Core.SetHTTPErrorHandler(render.HTTPErrorHandler(options))
	return s
}

// SetDomain 为模块设置域名
func (s *Hosts) SetDomain(name string, domain string) *Hosts {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	a, ok := s.mapNames[name]
	if !ok {
		s.Core.Logger().Warn(`Module does not exist: `, name)
		return s
	}
	if len(domain) == 0 { // 取消域名，加入到Core的Group中
		domain = a.Domain
		if _, ok := s.mapHosts[domain]; !ok {
			return s
		}
		delete(s.mapHosts, domain)
		var prefix string
		if name != s.RootName {
			prefix = `/` + name
			a.Dir = prefix + `/`
		} else {
			a.Dir = `/`
		}
		routes := a.Handler.Routes()
		for _, r := range routes {
			if r.Path == `/` {
				if len(prefix) > 0 {
					r.Path = prefix
					r.Format = prefix
				}
			} else {
				r.Path = prefix + r.Path
				r.Format = prefix + r.Format
			}
			r.Prefix = prefix
		}
		a.URL = a.Dir
		if s.URL != `/` {
			a.URL = strings.TrimSuffix(s.URL, `/`) + a.URL
		}
		a.Domain = ``
		a.Group = s.Core.Group(prefix)
		a.Group.Use(a.Middlewares...)
		s.Core.AppendRouter(routes)
		a.Handler = nil
		return s
	}
	if len(a.Domain) > 0 { // 从一个域名换为另一个域名
		if a.Domain == domain {
			return s
		}
		if _, ok := s.mapHosts[a.Domain]; ok {
			delete(s.mapHosts, a.Domain)
		}
		s.mapHosts[domain] = a
		a.Domain = domain
		return s
	}
	// 从Group移到域名
	s.mapHosts[domain] = a
	routes := []*echo.Route{}
	coreRoutes := []*echo.Route{}
	for _, r := range s.Core.Routes() {
		if r.Prefix == `/`+name {
			if r.Path == `/`+name {
				r.Path = `/`
				r.Format = `/`
			} else {
				r.Path = `/` + strings.TrimPrefix(r.Path, `/`+name+`/`)
				r.Format = `/` + strings.TrimPrefix(r.Format, `/`+name+`/`)
			}
			r.Prefix = ``
			routes = append(routes, r)
		} else {
			coreRoutes = append(coreRoutes, r)
		}
	}
	a.Domain = domain
	a.Group = nil
	e := echo.NewWithContext(s.NewContext)
	e.SetRenderer(a.Renderer)
	e.SetHTTPErrorHandler(s.Core.HTTPErrorHandler())
	e.Pre(s.DefaultPreMiddlewares...)
	e.Use(s.DefaultMiddlewares...)
	e.Use(a.Middlewares...)
	s.Core.RebuildRouter(coreRoutes)
	e.RebuildRouter(routes)
	a.Handler = e
	scheme := `http`
	if s.SessionOptions.Secure {
		scheme = `https`
	}
	a.URL = scheme + `://` + a.Domain + `/`
	a.Dir = `/`
	return s
}

// NewHost 创建新Host
func (s *Hosts) NewHost(name string, middlewares ...interface{}) *Host {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	r := strings.SplitN(name, `@`, 2) //blog@www.blog.com
	var domain string
	if len(r) > 1 {
		name = r[0]
		domain = r[1]
	}
	a := NewHost(name, domain, s, middlewares...)
	if len(domain) > 0 {
		s.mapHosts[domain] = a
	}
	s.mapNames[name] = a
	return a
}

// Get 获取Host实例
func (s *Hosts) Get(args ...string) *Host {
	name := s.RootName
	if len(args) > 0 {
		name = args[0]
		if ap, ok := s.mapNames[name]; ok {
			return ap
		}
	}
	return s.NewHost(name)
}

// SetSessionOptions 设置Session配置
func (s *Hosts) SetSessionOptions(sessionOptions *echo.SessionOptions) *Hosts {
	if sessionOptions.CookieOptions == nil {
		sessionOptions.CookieOptions = &echo.CookieOptions{
			Path:     `/`,
			HttpOnly: true,
		}
	}
	if len(sessionOptions.Name) == 0 {
		sessionOptions.Name = `GOSID`
	}
	if len(sessionOptions.Engine) == 0 {
		sessionOptions.Engine = `cookie`
	}
	s.SessionOptions = sessionOptions
	return s
}

// GetOk 获取Host实例
func (s *Hosts) GetOk(args ...string) (app *Host, ok bool) {
	name := s.RootName
	if len(args) > 0 {
		name = args[0]
	}
	app, ok = s.mapNames[name]
	return
}

// ListAll 获取Host列表，如果传递参数值为true，返回所有域名所对应的模块列表
func (s *Hosts) ListAll(args ...bool) map[string]*Host {
	if len(args) > 0 && args[0] {
		return s.mapHosts
	}
	return s.mapNames
}

// Exists 检查Host是否存在
func (s *Hosts) Exists(name string) bool {
	_, ok := s.mapNames[name]
	return ok
}

// RootDir 网站根目录
func (s *Hosts) RootDir() string {
	if len(s.rootDir) == 0 {
		ppath := os.Getenv(strings.ToUpper(s.Name) + `PATH`)
		if len(ppath) == 0 {
			ppath, _ = filepath.Abs(os.Args[0])
			ppath = filepath.Dir(ppath)
		}
		s.rootDir = ppath
	}
	return s.rootDir
}

// Debug 开关debug模式
func (s *Hosts) Debug(on bool) *Hosts {
	s.Core.SetDebug(on)
	return s
}

// Run 运行服务
func (s *Hosts) Run(args ...interface{}) error {
	for _, f := range s.before {
		if err := f(); err != nil {
			return err
		}
	}
	var eng engine.Engine
	var arg interface{}
	size := len(args)
	if size > 0 {
		arg = args[0]
	}
	if size > 1 {
		if conf, ok := arg.(*engine.Config); ok {
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.NewWithConfig(conf)
				} else {
					eng = standard.NewWithConfig(conf)
				}
			} else {
				eng = fasthttp.NewWithConfig(conf)
			}
		} else {
			addr := `:80`
			if v, ok := arg.(string); ok && len(v) > 0 {
				addr = v
			}
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.New(addr)
				} else {
					eng = standard.New(addr)
				}
			} else {
				eng = fasthttp.New(addr)
			}
		}
	} else {
		switch v := arg.(type) {
		case string:
			eng = fasthttp.New(v)
		case engine.Engine:
			eng = v
		default:
			eng = fasthttp.New(`:80`)
		}
	}
	s.Core.Logger().Infof(`Server "%v" has been launched.`, s.Name)
	s.Core.Run(eng, s)
	s.Core.Logger().Infof(`Server "%v" has been closed.`, s.Name)
	for _, f := range s.after {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// InitCodec 初始化 加密/解密 接口
func (s *Hosts) InitCodec(hashKey []byte, blockKey []byte) {
	s.Codec = codec.New(hashKey, blockKey)
}

// Pprof 启用pprof
func (s *Hosts) Pprof() *Hosts {
	pprof.Wrapper(s.Core)
	return s
}

// DefaultFuncMap 模板的默认函数
func (s *Hosts) DefaultFuncMap() (r map[string]interface{}) {
	r = tplfunc.New()
	r["RootURL"] = func(p ...string) string {
		if len(p) > 0 {
			return s.URL + p[0]
		}
		return s.URL
	}
	return
}

// FuncMapCopyTo 获取模板的默认函数副本
func (s *Hosts) FuncMapCopyTo(m map[string]interface{}) *Hosts {
	for k, v := range s.FuncMap {
		m[k] = v
	}
	return s
}
