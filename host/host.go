/*

   Copyright 2019 Wenhui Shen <www.webx.top>

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
	"strings"
	"sync"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
)

// NewHost 创建Host实例
func NewHost(name string, domain string, s *Hosts, middlewares ...interface{}) (a *Host) {
	a = &Host{
		Hosts:       s,
		Name:        name,
		Domain:      domain,
		Middlewares: middlewares,
		lock:        &sync.RWMutex{},
	}
	if s.Renderer != nil {
		a.Renderer = s.Renderer
	}
	if len(a.Domain) == 0 {
		var prefix string
		if name != s.RootName {
			prefix = `/` + name
			a.Dir = prefix + `/`
		} else {
			a.Dir = `/`
		}
		a.URL = a.Dir
		if s.URL != `/` {
			a.URL = strings.TrimSuffix(s.URL, `/`) + a.URL
		}
		a.Group = s.Core.Group(prefix)
		a.Group.Use(middlewares...)
	} else {
		var e *echo.Echo
		if s.NewContext == nil {
			e = echo.New()
		} else {
			e = echo.NewWithContext(s.NewContext)
		}
		e.SetRenderer(s.Core.Renderer())
		e.SetHTTPErrorHandler(s.Core.HTTPErrorHandler())
		e.Pre(s.DefaultPreMiddlewares...)
		e.Use(s.DefaultMiddlewares...)
		e.Use(middlewares...)
		a.Handler = e
		scheme := `http`
		if s.SessionOptions.Secure {
			scheme = `https`
		}
		a.URL = scheme + `://` + a.Domain + `/`
		a.Dir = `/`
	}
	if s.RootName == name {
		a.Installed = time.Now().Unix()
	}
	return
}

// Host 定义
type Host struct {
	*Hosts         `json:"-" xml:"-"`
	Group          *echo.Group   `json:"-" xml:"-"`
	Handler        *echo.Echo    `json:"-" xml:"-"` //指定域名时有效
	PreMiddlewares []interface{} `json:"-" xml:"-"`
	Middlewares    []interface{} `json:"-" xml:"-"`
	Renderer       driver.Driver `json:"-" xml:"-"`
	Name           string
	Domain         string
	URL            string
	Dir            string
	Config         Configer
	Installed      int64 //安装时间戳
	Expired        int64 //过期时间戳
	Disabled       bool  //是否禁用
	// 安装和卸载逻辑
	Install   func() error `json:"-" xml:"-"`
	Uninstall func() error `json:"-" xml:"-"`

	lock *sync.RWMutex
}

// Valid 验证服务是否有效
func (a *Host) Valid() error {
	if a.Installed == 0 {
		return ErrHasNotBeenInstalled
	}
	if a.Disabled {
		return ErrHasBeenDisabled
	}
	if a.Expired > 0 && a.Expired < time.Now().Unix() {
		return ErrHasExpired
	}
	return nil
}

// Pre 前置中间件
func (a *Host) Pre(middleware ...interface{}) {
	if a.Handler != nil {
		a.Handler.Pre(middleware...)
	} else {
		a.Core.Pre(middleware...)
	}
	a.PreMiddlewares = append(middleware, a.PreMiddlewares...)

}

// Use 中间件
func (a *Host) Use(middleware ...interface{}) {
	a.Router().Use(middleware...)
	a.Middlewares = append(a.Middlewares, middleware...)
}

// Register 注册路由：module.Register(`/index`,Index.Index,"GET","POST")
func (a *Host) Register(ppath string, handler interface{}, methods ...string) *Host {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	a.Router().Match(methods, ppath, handler)
	return a
}

// Router 路由
func (a *Host) Router() echo.ICore {
	if a.Group != nil {
		return a.Group
	}
	return a.Handler
}
