package dbconfig

import (
	"net/url"
	"strings"
)

type Encoder func(*Config) string

var encoders = map[string]Encoder{
	`mysql`: func(c *Config) string {
		var host string
		if strings.HasPrefix(c.Host, `unix:`) {
			host = "unix(" + strings.TrimPrefix(c.Host, `unix:`) + ")"
		} else {
			host = "tcp(" + c.Host
			if len(c.Port) > 0 {
				host += ":" + c.Port
			}
			host += ")"
		}
		dsn := url.QueryEscape(c.User) + ":" + url.QueryEscape(c.Pass) + "@" + host + "/" + c.Name
		if len(c.Charset) > 0 {
			dsn += "?charset=" + c.Charset
		}
		return dsn
	},
	`mymysql`: func(c *Config) string {
		//tcp:localhost:3306*gotest/root/root
		var host string
		if strings.HasPrefix(c.Host, `unix:`) {
			host = c.Host
		} else {
			host = "tcp:" + c.Host
			if len(c.Port) > 0 {
				host += ":" + c.Port
			}
		}
		dsn := host + "*" + c.Name + "/" + url.QueryEscape(c.User) + "/" + url.QueryEscape(c.Pass)
		return dsn
	},
}

func EncoderRegister(engine string, encoder Encoder) {
	encoders[engine] = encoder
}

func EncoderGet(engine string) Encoder {
	encoder, _ := encoders[engine]
	return encoder
}

func EncoderUnregister(engine string) bool {
	_, ok := encoders[engine]
	if ok {
		delete(encoders, engine)
	}
	return ok
}
