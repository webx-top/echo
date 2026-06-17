# Echo
[![Go](https://github.com/webx-top/echo/actions/workflows/go.yml/badge.svg)](https://github.com/webx-top/echo/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/webx-top/echo)](https://goreportcard.com/report/github.com/webx-top/echo) [English](README.md)
#### Echo 是一个快速且简约的 Go 语言 Web 框架。速度远超同类框架。
本包需要 >= **go 1.25**

## 特性

- 优化的 HTTP 路由器，智能路由优先级
- 构建健壮且可扩展的 RESTful API
- 支持标准 HTTP 服务器或 FastHTTP 服务器
- API 分组
- 可扩展的中间件框架
- 可在根、分组或路由级别定义中间件
- 便捷的多种 HTTP 响应发送函数
- 集中的 HTTP 错误处理
- 支持多引擎模板渲染（standard、Jet、SSE 等）
- 可自定义日志格式
- 高度可定制
- 业务错误码系统
- 表单数据过滤与绑定
- 事务支持
- 反向代理与负载均衡（随机、轮询）
- URL 重写与逆向还原
- 速率限制（支持内存和 Redis 后端）
- 请求队列并发控制
- IP 过滤
- 分布式追踪（OpenTracing）
- 服务端推送事件（SSE）

## 快速开始

### 安装

```sh
$ go get github.com/webx-top/echo
```

### Hello, World!

创建 `server.go`

```go
package main

import (
	"net/http"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
)

func main() {
	e := echo.New()
	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!", http.StatusOK)
	})
	e.Run(standard.New(":1323"))
}
```

启动服务器

```sh
$ go run server.go
```

访问 [http://localhost:1323](http://localhost:1323) 即可看到 Hello, World!

### 路由

```go
e.Post("/users", saveUser)
e.Get("/users/:id", getUser)
e.Put("/users/:id", updateUser)
e.Delete("/users/:id", deleteUser)
e.Get("/user/<id:[\\d]+>", getUser)
```

支持两种参数语法：`:param` 和 `<param:regexp>`。

### 路径参数

```go
func getUser(c echo.Context) error {
	// 从路径 `users/:id` 获取用户 ID
	id := c.Param("id")
	// 或 id := c.Paramx("id").Uint64()
}
```

### 查询参数

`/show?team=x-men&member=wolverine`

```go
func show(c echo.Context) error {
	// 从查询字符串获取 team 和 member
	team := c.Query("team")
	member := c.Query("member")
	age := c.Queryx("age").Uint()
}
```

### 表单 `application/x-www-form-urlencoded`

`POST` `/save`

名称 | 值
:--- | :---
name | Joe Smith
email | joe@labstack.com


```go
func save(c echo.Context) error {
	// 获取 name 和 email
	name := c.Form("name")
	email := c.Form("email")
	age := c.Formx("age").Uint()
}
```

### 表单 `multipart/form-data`

`POST` `/save`

名称 | 值
:--- | :---
name | Joe Smith
email | joe@labstack.com
avatar | avatar

```go
func save(c echo.Context) error {
	// 获取 name 和 email
	name := c.Form("name")
	email := c.Form("email")

	//------------
	// 获取头像
	//------------
	_, err := c.SaveUploadedFile("avatar","./")
	return err
}
```

### 请求处理

- 根据 `Content-Type` 请求头将 `JSON` 或 `XML` 绑定到 Go 结构体
- 以 `JSON` 或 `XML` 格式返回响应

```go
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

e.Post("/users", func(c echo.Context) error {
	u := new(User)
	if err := c.MustBind(u); err != nil {
		return err
	}
	return c.JSON(u, http.StatusCreated)
	// 或
	// return c.XML(u, http.StatusCreated)
})
```

### 表单数据过滤 (FormFilter)

在绑定前过滤和转换表单数据。

```go
import "github.com/webx-top/echo/formfilter"

ff := formfilter.New()
ff.Add(
    formfilter.StringToSlice("ids", ","),
    formfilter.DateRange("created", "2006-01-02"),
)

// 将过滤器直接传给 Bind/MustBind
u := new(User)
c.MustBind(u, ff.Build())
```

### 静态文件

为 `/static/*` 路径提供静态文件服务。

```go
e.Use(mw.Static(&mw.StaticOptions{
	Root:"static", // 静态文件物理路径
	Path:"/static/", // 网址访问静态文件的路径
	Browse:true, // 是否显示文件列表
}))
```

### 嵌入文件处理器

通过 `embed.FileSystems` 提供嵌入文件服务，支持 index.html。

```go
import "github.com/webx-top/echo/handler/embed"

//go:embed www/*
var wwwFS embed.FS

fs := embed.NewFileSystems()
fs.Register(wwwFS)
e.Get("/*", embed.File(fs))
```

### 嵌入静态文件 (Bindata)

以中间件方式提供 go-bindata 生成的嵌入静态文件。

```go
import (
	"github.com/admpub/go-bindata-assetfs"
	"github.com/webx-top/echo/middleware/bindata"
)

//go:generate go-bindata -o=bindata_assetfs.go -pkg=main static/...
//go:generate go-bindata-assetfs -o=bindata_assetfs.go -pkg=main static/...

func NewAssetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,     // 由 go-bindata 生成
		AssetDir:  AssetDir,  // 由 go-bindata 生成
		AssetInfo: AssetInfo, // 由 go-bindata 生成
		Prefix:    "",
	}
}

e.Use(bindata.Static("/static/", NewAssetFS()))
```

### 验证码

支持图片和音频验证码。

```go
import "github.com/webx-top/echo/handler/captcha"

captcha.DefaultOptions.Wrapper(e)
// 访问 /captcha/<id>.png
```

### 中间件

```go
// 根级别中间件
e.Use(middleware.Log())
e.Use(middleware.Recover())

// 分组级别中间件
g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string) bool {
	if username == "joe" && password == "secret" {
		return true
	}
	return false
}))

// 路由级别中间件
track := func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("request to /users")
		return next.Handle(c)
	}
}
e.Get("/users", func(c echo.Context) error {
	return c.String("/users", http.StatusOK)
}, track)
```

### Cookie
```go
e.Get("/setcookie", func(c echo.Context) error {
	c.SetCookie("uid","1")
	return c.String("/setcookie: uid="+c.GetCookie("uid"), http.StatusOK)
})
```

### Session

Session 引擎支持基于 cookie 和文件的存储。

```go
...
import (
	...
	"github.com/webx-top/echo/middleware/session"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)
...
sessionOptions := &echo.SessionOptions{
	Engine: `cookie`,
	Name:   `SESSIONID`,
	CookieOptions: &echo.CookieOptions{
		Path:     `/`,
		Domain:   ``,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	},
}

cookieStore.RegWithOptions(&cookieStore.CookieOptions{
	KeyPairs: [][]byte{
		[]byte(`123456789012345678901234567890ab`),
	},
})

e.Use(session.Middleware(sessionOptions))

e.Get("/session", func(c echo.Context) error {
	c.Session().Set("uid",1).Save()
	return c.String(fmt.Sprintf("/session: uid=%v",c.Session().Get("uid")))
})
```

### 事务

通过上下文支持数据库事务。

```go
// Transaction 接口
type Transaction interface {
	Begin(ctx context.Context) error
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
	End(ctx context.Context, succeed bool) error
}

// 在处理器中使用
c.Begin()
// ... 业务操作
c.Commit()
```

### Websocket
```go
...
import (
	...
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"
)
...

e.AddHandlerWrapper(ws.HanderWrapper)

e.Get("/websocket", func(c *websocket.Conn, ctx echo.Context) error {
	// 推送消息
	go func() {
		var counter int
		for {
			if counter >= 10 {
				return
			}
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			ctx.Logger().Info(`推送消息: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				ctx.Logger().Error(`推送错误: `, err.Error())
				return
			}
			counter++
		}
	}()

	// 回显
	ws.DefaultExecuter(c, ctx)
	return nil
})
```
[更多...](https://github.com/webx-top/echo/blob/master/handler/websocket/example/main.go)

### Sockjs
```go
...
import (
	...
	"github.com/webx-top/echo"
	"github.com/admpub/sockjs-go/v3/sockjs"
	ws "github.com/webx-top/echo/handler/sockjs"
)
...

options := ws.Options{
	Handle: func(c sockjs.Session) error {
		// 推送消息
		go func() {
			var counter int
			for {
				if counter >= 10 {
					return
				}
				time.Sleep(5 * time.Second)
				message := time.Now().String()
				log.Info(`推送消息: `, message)
				if err := c.Send(message); err != nil {
					log.Error(`推送错误: `, err.Error())
					return
				}
				counter++
			}
		}()

		// 回显
		ws.DefaultExecuter(c)
		return nil
	},
	Options: &sockjs.DefaultOptions,
	Prefix:  "/websocket",
}
options.Wrapper(e)
```
[更多...](https://github.com/webx-top/echo/blob/master/handler/sockjs/example/main.go)

### 服务端推送事件 (SSE)

```go
import (
	_ "github.com/webx-top/echo/middleware/render/sse"
	"github.com/webx-top/echo/middleware/render"
)

e.Use(render.Middleware(render.New(`sse`, ``)))

e.Get("/events", func(c echo.Context) error {
	return c.SSEvent("message", listener)
})
```

### 反向代理

内置反向代理，支持负载均衡（随机、轮询）。

```go
import "github.com/webx-top/echo/middleware"

// 轮询均衡器
targets := []middleware.ProxyTargeter{
	&middleware.ProxyTarget{
		Name: "api-1",
		URL:  mustParseURL("http://localhost:8081"),
	},
	&middleware.ProxyTarget{
		Name: "api-2",
		URL:  mustParseURL("http://localhost:8082"),
	},
}
balancer := middleware.NewRoundRobinBalancer(targets)
e.Use(middleware.Proxy(balancer))
```

### URL 重写

```go
e.Use(middleware.Rewrite(map[string]string{
	"/old":              "/new",
	"/api/*":            "/$1",
	"/users/:id":        "/user/$1",
	"/users/*/orders/*": "/user/$1/order/$2",
}))
```

### 请求队列

通过队列超时限制并发请求处理。

```go
config := middleware.QueueConfig{
	QueueSize:     100,
	Workers:       10,
	QueueTimeout:  30 * time.Second,
	WorkerTimeout: 10 * time.Second,
}
e.Use(middleware.QueueWithConfig(config))
```

### 请求 ID

```go
e.Use(middleware.RequestID())
// 每个响应都会包含 X-Request-ID 头
```

### IP 过滤

```go
import "github.com/webx-top/echo/middleware/ipfilter"

e.Use(ipfilter.IPFilter(ipfilter.Config{
	Options: ipfilter.Options{
		AllowedIPs: []string{"192.168.0.0/24"},
		BlockedIPs: []string{"0.0.0.0/0"},
	},
}))
```

### OpenTracing

```go
import "github.com/webx-top/echo/middleware/opentracing"

e.Use(opentracing.Trace(tracer))
```

### 其他示例

```go
package main

import (
	"net/http"

	"github.com/webx-top/echo"
	// "github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(mw.Log())

	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!")
	})
	e.Get("/echo/:name", func(c echo.Context) error {
		return c.String("Echo " + c.Param("name"))
	})
	
	e.Get("/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`标准 net/http 处理函数`))
		w.WriteHeader(200)
	})

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// 标准
	e.Run(standard.New(":4444"))
}
```

[查看更多示例...](https://github.com/admpub/echo-example/blob/master/_v2/main.go)

## 中间件列表

### 根中间件 (github.com/webx-top/echo/middleware)

中间件 | 导入路径 | 说明
:----------|:------------|:-----------
[BasicAuth](https://github.com/webx-top/echo/blob/master/middleware/auth_basic.go) | github.com/webx-top/echo/middleware | HTTP 基本认证
[KeyAuth](https://github.com/webx-top/echo/blob/master/middleware/auth_key.go) | github.com/webx-top/echo/middleware | API 密钥认证（header/query/form）
[BodyLimit](https://github.com/webx-top/echo/blob/master/middleware/bodylimit.go) | github.com/webx-top/echo/middleware | 限制请求体大小
[Gzip](https://github.com/webx-top/echo/blob/master/middleware/compress.go) | github.com/webx-top/echo/middleware | gzip 压缩响应
[Secure](https://github.com/webx-top/echo/blob/master/middleware/secure.go) | github.com/webx-top/echo/middleware | 防 XSS、点击劫持等攻击
[CORS](https://github.com/webx-top/echo/blob/master/middleware/cors.go) | github.com/webx-top/echo/middleware | 跨域资源共享
[CSRF](https://github.com/webx-top/echo/blob/master/middleware/csrf.go) | github.com/webx-top/echo/middleware | 跨站请求伪造防护
[Log](https://github.com/webx-top/echo/blob/master/middleware/log.go) | github.com/webx-top/echo/middleware | 记录 HTTP 请求日志
[MethodOverride](https://github.com/webx-top/echo/blob/master/middleware/method_override.go) | github.com/webx-top/echo/middleware | 通过请求头覆盖请求方法
[Recover](https://github.com/webx-top/echo/blob/master/middleware/recover.go) | github.com/webx-top/echo/middleware | 从 panic 中恢复
[HTTPSRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go) | github.com/webx-top/echo/middleware | HTTP 重定向到 HTTPS
[HTTPSWWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go) | github.com/webx-top/echo/middleware | HTTP 重定向到 WWW HTTPS
[WWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go) | github.com/webx-top/echo/middleware | 非 WWW 重定向到 WWW
[NonWWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go) | github.com/webx-top/echo/middleware | WWW 重定向到非 WWW
[AddTrailingSlash](https://github.com/webx-top/echo/blob/master/middleware/slash.go) | github.com/webx-top/echo/middleware | 添加尾部斜杠
[RemoveTrailingSlash](https://github.com/webx-top/echo/blob/master/middleware/slash.go) | github.com/webx-top/echo/middleware | 移除尾部斜杠
[Static](https://github.com/webx-top/echo/blob/master/middleware/static.go) | github.com/webx-top/echo/middleware | 提供静态文件服务
[MaxAllowed](https://github.com/webx-top/echo/blob/master/middleware/limit.go) | github.com/webx-top/echo/middleware | 限制同时请求数
[NoCache](https://github.com/webx-top/echo/blob/master/middleware/nocache.go) | github.com/webx-top/echo/middleware | 设置无缓存头
[Rewrite](https://github.com/webx-top/echo/blob/master/middleware/rewrite.go) | github.com/webx-top/echo/middleware | URL 路径重写
[Proxy](https://github.com/webx-top/echo/blob/master/middleware/proxy.go) | github.com/webx-top/echo/middleware | 反向代理与负载均衡
[Queue](https://github.com/webx-top/echo/blob/master/middleware/queue.go) | github.com/webx-top/echo/middleware | 请求并发队列
[RequestID](https://github.com/webx-top/echo/blob/master/middleware/request_id.go) | github.com/webx-top/echo/middleware | X-Request-ID 头
[Validate](https://github.com/webx-top/echo/blob/master/middleware/validate.go) | github.com/webx-top/echo/middleware | 请求验证器注入
[FuncMap](https://github.com/webx-top/echo/blob/master/middleware/funcmap.go) | github.com/webx-top/echo/middleware | 模板函数映射
[AJAX](https://github.com/webx-top/echo/blob/master/middleware/ajax.go) | github.com/webx-top/echo/middleware | AJAX 操作处理器
[Language](https://github.com/webx-top/echo/tree/master/middleware/language) | github.com/webx-top/echo/middleware/language | 多语言 (i18n) 支持
[RateLimit](https://github.com/webx-top/echo/tree/master/middleware/ratelimit) | github.com/webx-top/echo/middleware/ratelimit | 速率限制
[RateLimiter](https://github.com/webx-top/echo/tree/master/middleware/ratelimiter) | github.com/webx-top/echo/middleware/ratelimiter | 速率限制器（支持内存和 Redis）
[Session](https://github.com/webx-top/echo/blob/master/middleware/session/middleware.go) | github.com/webx-top/echo/middleware/session | Session 管理器
[JWT](https://github.com/webx-top/echo/blob/master/middleware/jwt/jwt.go) | github.com/webx-top/echo/middleware/jwt | JWT 认证
[Markdown](https://github.com/webx-top/echo/blob/master/middleware/markdown/markdown.go) | github.com/webx-top/echo/middleware/markdown | Markdown 渲染
[Render](https://github.com/webx-top/echo/blob/master/middleware/render/middleware.go) | github.com/webx-top/echo/middleware/render | HTML 模板渲染（standard、Jet、SSE）
[IPFilter](https://github.com/webx-top/echo/tree/master/middleware/ipfilter) | github.com/webx-top/echo/middleware/ipfilter | IP 地址过滤
[OpenTracing](https://github.com/webx-top/echo/tree/master/middleware/opentracing) | github.com/webx-top/echo/middleware/opentracing | 分布式追踪
[Bindata Static](https://github.com/webx-top/echo/tree/master/middleware/bindata) | github.com/webx-top/echo/middleware/bindata | 提供嵌入静态文件服务
[ReverseProxy](https://github.com/webx-top/reverseproxy/blob/master/middleware.go) | github.com/webx-top/reverseproxy | 外部反向代理

## 处理器包装器列表
包装器 | 导入路径 | 说明
:-------|:------------|:-----------
Websocket | github.com/webx-top/echo/handler/websocket | [示例](https://github.com/webx-top/echo/blob/master/handler/websocket/example/main.go)
Sockjs | github.com/webx-top/echo/handler/sockjs | [示例](https://github.com/webx-top/echo/blob/master/handler/sockjs/example/main.go)
Oauth2 | github.com/webx-top/echo/handler/oauth2 | [示例](https://github.com/webx-top/echo/blob/master/handler/oauth2/example/main.go)
Pprof | github.com/webx-top/echo/handler/pprof | Go pprof 性能分析
Captcha | github.com/webx-top/echo/handler/captcha | 图片和音频验证码
Embed | github.com/webx-top/echo/handler/embed | 嵌入文件系统文件服务
SSE | github.com/webx-top/echo/middleware/render/sse | 服务端推送事件渲染驱动

## 模板函数文档
[模板函数文档](middleware/tplfunc/TplFuncMap_Documentation.md)

## 模板数据

使用 Render 中间件时，传给 `c.Render(data)` 的 data 会被包装为 `*echo.RenderData`。在 Go 模板中可通过 `$.XXX` 调用其所有公开方法：

```gotmpl
{{$.Now}}           {{/* 当前时间 */}}
{{$.UnixTime}}      {{/* 当前时间戳 */}}
{{$.Site}}          {{/* 网站URL */}}
{{$.SiteRoot}}      {{/* 网站根路径 */}}
{{$.URL}}           {{/* 当前请求URL */}}
{{$.URI}}           {{/* 当前请求URI */}}
{{$.Path}}          {{/* 当前请求路径 */}}
{{$.Domain}}        {{/* 当前域名 */}}
{{$.Port}}          {{/* 当前端口 */}}
{{$.Scheme}}        {{/* http/https */}}
{{$.Lang}}          {{/* 当前语言 */}}
{{$.Referer}}       {{/* 来源地址 */}}
{{$.Query "key"}}   {{/* 获取查询参数 */}}
{{$.Form "key"}}    {{/* 获取表单参数 */}}
{{$.Param "key"}}   {{/* 获取路径参数 */}}
{{$.Get "key"}}     {{/* 获取Context中存储的值 */}}
{{$.Cookie}}        {{/* Cookie操作 */}}
{{$.Session}}       {{/* Session操作 */}}
{{$.Flash "key"}}   {{/* Flash消息 */}}
{{$.T "你好%v" "世界"}}  {{/* 多语言翻译 */}}
{{$.LangURI "zh"}}  {{/* 生成语言链接 */}}
{{$.URLByName "routeName" "param1"}}  {{/* 根据路由名生成URL */}}
{{$.CaptchaForm}}   {{/* 验证码表单 */}}
{{$.TimeAgo $time}} {{/* 时间友好显示 */}}
{{$.TsHumanize $time}}{{/* 时间区间友好显示 */}}
{{$.DurationFormat $t}}{{/* 持续时间格式化 */}}
{{$.Fetch "subtmpl" .}}{{/* 渲染子模板并嵌入 */}}
{{$.MakeURL "handler" "arg"}}{{/* 生成URL */}}
{{$.Ext}}           {{/* 默认扩展名 */}}
{{$.ThemeColor}}    {{/* 主题色 */}}
{{$.Prefix}}        {{/* 当前路由前缀 */}}
{{$.RootPrefix}}    {{/* 根路由前缀 */}}
{{$.UploadURL "subdir"}}{{/* 上传文件URL */}}
{{$.FullURL "/path"}}{{/* 生成完整URL */}}
{{$.HasAnyRequest}} {{/* 是否有任何请求数据 */}}
{{$.GetNextURL}}    {{/* 获取跳转URL */}}
{{$.ReturnToCurrentURL}}{{/* 返回当前URL */}}
```

其中 `$.Data` 为原始传入的数据对象，`$.Stored` 为通过 `c.Set("key",any)` 存储的只读数据。

## 附加包

- `formfilter` - 表单数据过滤工具
- `subdomains` - 子域名路由工具（含 SafeMap）
- `code` - 业务错误码系统
- `code/register` - 错误码注册
- `encoding/dbconfig` - 数据库配置编码
- `param` - 参数类型工具（StringSlice、StringMap、Store）
- `testing` - HTTP 测试工具
- `mockcontext` - 单元测试用的 Mock 上下文
- `defaults` - 默认配置工具
- `logger` - 日志集成

## 案例

- [Nging](https://github.com/admpub/nging)

## 贡献者

- [Vishal Rana](https://github.com/vishr) - 原作者
- [Hank Shen](https://github.com/admpub) - 作者
- [Nitin Rana](https://github.com/nr17) - 顾问
- [贡献者列表](https://github.com/webx-top/echo/graphs/contributors)

## 许可证

[Apache 2](https://github.com/webx-top/echo/blob/master/LICENSE)
