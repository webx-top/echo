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

package ipfilter

import (
	"net"

	"github.com/admpub/ipfilter"
	"github.com/webx-top/echo"
)

type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper echo.Skipper `json:"-"`
	Options ipfilter.Options
	filter  *ipfilter.IPFilter
}

func (c *Config) Init() {
	c.filter = ipfilter.New(c.Options)
}

func (c *Config) Filter() *ipfilter.IPFilter {
	if c.filter == nil {
		c.Init()
	}
	return c.filter
}

var (
	// DefaultConfig is the default ipfilter middleware config.
	DefaultConfig = Config{
		Skipper: echo.DefaultSkipper,
		Options: ipfilter.Options{},
	}
)

// IPFilter returns a IPFilter middleware with config.
func IPFilter(config Config) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	config.Init()
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}
			var ip string
			if config.Options.TrustProxy {
				ip = c.RealIP()
			} else {
				ip, _, _ = net.SplitHostPort(c.Request().RemoteAddress())
			}
			//show simple forbidden text
			if !config.Filter().Allowed(ip) {
				return echo.ErrForbidden
			}
			return next.Handle(c)
		}
	}
}
