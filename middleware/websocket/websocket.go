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
	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
)

var DefaultUpgrader = &websocket.Upgrader{}

func HanderWrapper(v interface{}) echo.Handler {
	if h, ok := v.(func(*websocket.Conn, echo.Context) error); ok {
		return Websocket(h)
	}
	return nil
}

func Websocket(executer func(*websocket.Conn, echo.Context) error, opts ...*websocket.Upgrader) echo.HandlerFunc {
	var opt *websocket.Upgrader
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = DefaultUpgrader
	}
	if executer == nil {
		//Test mode
		executer = func(c *websocket.Conn, ctx echo.Context) error {
			mt, message, err := c.ReadMessage()
			if err != nil {
				return err
			}
			log.Infof("Websocket recv: %s", message)
			return c.WriteMessage(mt, message)
		}
	}
	h := func(ctx echo.Context) error {
		w := ctx.Response().StdResponseWriter()
		r := ctx.Request().StdRequest()
		c, err := opt.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
		defer c.Close()
		for {
			if err = executer(c, ctx); err != nil {
				break
			}
		}
		return err
	}
	return echo.HandlerFunc(h)
}
