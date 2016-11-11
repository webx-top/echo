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
package websocket

import (
	"net/http"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

type Options struct {
	Path     string                      `json:"path"` //UrlPath
	Upgrader *websocket.Upgrader         `json:"upgrader"`
	Executer func(*websocket.Conn) error `json:"-"`
}

var DefaultOptions = &Options{
	Path:     "/websocket/",
	Upgrader: &websocket.Upgrader{},
	Executer: func(c *websocket.Conn) error {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return err
		}
		log.Info("Websocket recv: %s", message)
		return c.WriteMessage(mt, message)
	},
}

func Websocket(opts ...*Options) echo.Middleware {
	var opt *Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = DefaultOptions
	}
	h := func(w http.ResponseWriter, r *http.Request) error {
		c, err := opt.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
		defer c.Close()
		for {
			if err = opt.Executer(c); err != nil {
				break
			}
		}
		return err
	}
	return echo.MiddlewareFunc(func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			return h(c.Response().StdResponseWriter(), c.Request().StdRequest())
		})
	})
}
