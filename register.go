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

package echo

import "github.com/webx-top/echo/param"

var (
	_ ICore = &Echo{}
	_ ICore = &Group{}
)

type RouteRegister interface {
	MiddlewareRegister
	Group(prefix string, middleware ...any) *Group
	Any(path string, h any, middleware ...any) IRouter
	Route(methods string, path string, h any, middleware ...any) IRouter
	Match(methods []string, path string, h any, middleware ...any) IRouter
	Connect(path string, h any, m ...any) IRouter
	Delete(path string, h any, m ...any) IRouter
	Get(path string, h any, m ...any) IRouter
	Head(path string, h any, m ...any) IRouter
	Options(path string, h any, m ...any) IRouter
	Patch(path string, h any, m ...any) IRouter
	Post(path string, h any, m ...any) IRouter
	Put(path string, h any, m ...any) IRouter
	Trace(path string, h any, m ...any) IRouter
	Static(prefix, root string)
	File(path, file string)
	Prefix() string
}

type ContextRegister interface {
	SetContext(Context)
}

type RendererRegister interface {
	SetRenderer(Renderer)
}

type MiddlewareRegister interface {
	Use(middleware ...any)
	Pre(middleware ...any)
}

type URLBuilder interface {
	URL(any, ...any) string
}

type ICore interface {
	RouteRegister
	URLBuilder
	RendererRegister
	Prefixer
}

type IRouter interface {
	SetName(string) IRouter
	GetName() string
	SetMeta(param.Store) IRouter
	SetMetaKV(string, any) IRouter
	GetMeta() param.Store
	SetEncodingConfig(ef EncodingConfig) IRouter
	SetEncodingOmitFields(names ...string) IRouter
	SetEncodingOnlyFields(names ...string) IRouter
}

type Closer interface {
	Close() error
}

type Prefixer interface {
	Prefix() string
}

type URLGenerator interface {
	RelativeURL(uri string) string
	URLFor(uri string, relative ...bool) string
	URLByName(name string, args ...any) string
	RelativeURLByName(name string, args ...any) string
}

type IRouteDispatchPath interface {
	SetDispatchPath(route string)
	DispatchPath() string
}
