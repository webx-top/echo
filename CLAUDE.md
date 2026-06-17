# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

- `go test ./...` — run all tests
- `go test -v ./...` — verbose test output
- `go test -run TestName ./...` — specific test
- `go vet ./...` — static analysis
- `go build ./...` — build all packages

## Architecture

### Core Layers (top to bottom)

**`echo.Echo`** — main application struct. Holds router, middleware chains, binder, renderer, logger, engine, and all configuration. Created via `echo.New()`.

**`echo.Context`** — per-request interface. Wraps `engine.Request` + `engine.Response`, provides param/query/form access, session, transaction, validation, i18n, binding, JSON/XML/HTML responses.

**`echo.Router`** — radix tree (`node`) + static route map (`map[string]*methodHandler`). Route priority: static > named param(`:name`) > regex param(`<name:regexp>`). Registers via `echo.Group` (prefix-based route grouping).

### Engine Abstraction

Two interchangeable backend servers:

- **`engine/standard/`** — wraps `net/http` (stdlib)
- **`engine/fasthttp/`** — wraps `valyala/fasthttp`
- **`engine/mock/`** — mock request/response for testing

Both implement `engine.Request`, `engine.Response`, `engine.Engine` interfaces.

### Package Map

```
echo/
├── echo.go              — Echo struct: app lifecycle, routing registration, engine run
├── context.go           — Context interface (request-scoped data)
├── context_x*.go        — Context method implementations (request, response, session, store, transaction)
├── router.go            — Radix tree router
├── group.go             — Route grouping with prefix
├── binder*.go           — Request body binding (JSON/XML/Form/Multipart)
├── middleware.go        — Middleware types (Skipper, MiddlewareFunc, MiddlewareFuncd)
├── middleware/          — Built-in middleware (30+)
│   ├── *.go             — Root middleware: auth, cache, cors, csrf, log, proxy, rewrite, queue, etc.
│   ├── render/          — Template rendering (standard/jet/sse engines)
│   ├── session/         — Session (cookie/file engines)
│   ├── jwt/             — JWT auth
│   ├── ratelimit/       — Rate limiting
│   ├── ratelimiter/     — Rate limiter (memory + Redis)
│   ├── language/        — i18n
│   ├── ipfilter/        — IP allow/block
│   ├── opentracing/     — Distributed tracing
│   └── bindata/         — Embedded static files
├── handler/             — Handler wrappers (websocket, sockjs, oauth2, pprof, captcha, embed)
├── engine/              — Server engine implementations (standard, fasthttp, mock)
├── formfilter/          — Build FormDataFilter functions for pre-bind data transformation
├── param/               — Type-safe param access (String, StringSlice, Store)
├── code/                — Business error code system
├── logger/              — Logger interface + logzero integration
├── testing/             — HTTP test helpers
├── mockcontext/         — Mock Context for unit tests
└── subdomains/          — Subdomain routing utilities (SafeMap)
```

### Middleware Patterns

Two middleware function signatures exist:

- **`MiddlewareFunc`** = `func(Handler) Handler` — wraps any Handler
- **`MiddlewareFuncd`** = `func(Handler) HandlerFunc` — wraps Handler returning HandlerFunc

Both accept optional `Skipper` function via config structs. Middleware applies at root (`e.Use()`), group (`g.Use()`), or route level.

### Context Response Helpers

`c.String()`, `c.JSON()`, `c.XML()`, `c.HTML()`, `c.Blob()`, `c.SSEvent()`, `c.Stream()`, `c.ServeContent()`, `c.Redirect()`.

### Binder

`c.MustBind(i)` decodes request body into struct based on Content-Type. Accepts optional `FormDataFilter` variadic args for pre-processing form values. Struct tags: `json`, `xml`, `form`, `form_filter`.

### Key Conventions

- Module path: `github.com/webx-top/echo`
- Go 1.25 minimum
- Tests: standard `_test.go` files alongside packages
- Most middleware configs follow pattern: `DefaultXConfig` + `XWithConfig(config)` constructor
- Logger interface defined in `logger/`, default log adapter is `github.com/admpub/log`
