package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/webx-top/codec"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/oauth2"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/session"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)

func main() {
	account := &oauth2.Account{
		On:     true,
		Name:   `github`,
		Key:    os.Getenv(`OAUTH_CLIENT_ID`),
		Secret: os.Getenv(`OAUTH_CLIENT_SECRET`),
	}
	port := flag.String(`p`, "8080", "port")
	flag.StringVar(&account.Name, `provider`, account.Name, "-provider "+account.Name)
	flag.StringVar(&account.Key, `clientID`, account.Key, "-clientID <clientID>")
	flag.StringVar(&account.Secret, `clientSecret`, account.Secret, "-clientSecret <clientSecret>")
	flag.Parse()
	e := echo.New()
	e.Use(mw.Log())
	sessionOptions := &echo.SessionOptions{
		Engine: `cookie`,
		Name:   `SESSIONID`,
		CookieOptions: &echo.CookieOptions{
			Path:     `/`,
			Domain:   ``,
			MaxAge:   0,
			Secure:   false,
			HttpOnly: true,
			Cryptor:  echo.NewCookieCryptor(codec.Default, com.RandomString(16)),
		},
	}

	cookieStore.RegWithOptions(&cookieStore.CookieOptions{
		KeyPairs: [][]byte{
			[]byte(`123456789012345678901234567890ab`),
		},
	})

	e.Use(session.Middleware(sessionOptions))

	e.Get("/", func(c echo.Context) error {
		return c.HTML(`Login: <a href="/oauth/login/github" target="_blank">github</a>`)
	})
	config := oauth2.NewConfig()
	config.AddAccount(account)
	options := oauth2.New(`http://www.coscms.com`, config)
	options.SetSuccessHandler(func(ctx echo.Context) error {
		user := options.User(ctx)
		b, e := json.MarshalIndent(user, "", "  ")
		if e != nil {
			return e
		}
		return ctx.String(string(b))
	})
	options.Wrapper(e)

	switch `` {
	case `fast`:
		// FastHTTP
		e.Run(fasthttp.New("127.0.0.1:" + *port))

	default:
		// Standard
		e.Run(standard.New("127.0.0.1:" + *port))
	}
}
