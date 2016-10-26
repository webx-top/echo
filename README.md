# Echo
A fast and unfancy micro web framework for Go.

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
		return c.String(200, "Hello, World!")
	})
	e.Get("/echo/:name", func(c echo.Context) error {
		return c.String(200, "Echo " + c.Param("name"))
	})
	
	e.Get("/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`standard net/http handleFunc`))
		w.WriteHeader(200)
	})

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
```

[See other examples...](https://github.com/admpub/echo-example/blob/master/_v2/main.go)