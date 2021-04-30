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
package captcha

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
)

type IDEncoder func(_ echo.Context, id string, from string) (string, error)
type IDDecoder func(_ echo.Context, id string, from string) (string, error)

var DefaultOptions = &Options{
	EnableImage:    true,
	EnableAudio:    true,
	EnableDownload: true,
	AudioLangs:     []string{`zh`, `ru`, `en`},
	Prefix:         `/captcha`,
	CookieName:     `captchaId`,
	HeaderName:     `X-Captcha-ID`,
	IDEncoder: func(_ echo.Context, id string, _ string) (string, error) {
		return id, nil
	},
	IDDecoder: func(_ echo.Context, id string, _ string) (string, error) {
		return id, nil
	},
}

type Options struct {
	EnableImage    bool
	EnableAudio    bool
	EnableDownload bool
	AudioLangs     []string
	Prefix         string
	CookieName     string
	HeaderName     string
	IDEncoder      IDEncoder
	IDDecoder      IDDecoder
}

func (o Options) Wrapper(e echo.RouteRegister) {
	if o.AudioLangs == nil || len(o.AudioLangs) == 0 {
		o.AudioLangs = DefaultOptions.AudioLangs
	}
	if len(o.Prefix) == 0 {
		o.Prefix = DefaultOptions.Prefix
	}
	o.Prefix = strings.TrimRight(o.Prefix, "/")
	e.Get(o.Prefix+"/*", Captcha(&o))
}

func Captcha(opts ...*Options) func(echo.Context) error {
	var o *Options
	if len(opts) > 0 {
		o = opts[0]
	}
	if o == nil {
		o = DefaultOptions
	}
	if len(o.CookieName) == 0 {
		o.CookieName = DefaultOptions.CookieName
	}
	if len(o.HeaderName) == 0 {
		o.HeaderName = DefaultOptions.HeaderName
	}
	if o.IDDecoder == nil {
		o.IDDecoder = DefaultOptions.IDDecoder
	}
	if o.IDEncoder == nil {
		o.IDEncoder = DefaultOptions.IDEncoder
	}
	return func(ctx echo.Context) (err error) {
		var id, ext string
		param := ctx.P(0)
		if p := strings.LastIndex(param, `.`); p > 0 {
			id = param[0:p]
			ext = param[p:]
		}
		if len(ext) == 0 || len(id) == 0 {
			return echo.ErrNotFound
		}
		id, err = o.IDDecoder(ctx, id, `path`)
		if err != nil {
			return err
		}
		w := ctx.Response()
		header := w.Header()
		ids := []string{id}
		var hasCookieValue, hasHeaderValue bool
		if len(o.CookieName) > 0 {
			idByCookie := ctx.GetCookie(o.CookieName)
			if len(idByCookie) > 0 {
				idByCookie, err = o.IDDecoder(ctx, idByCookie, `cookie`)
				if err != nil {
					header.Add(`X-Captcha-ID-Error`, `cookie: `+err.Error())
					ctx.SetCookie(o.CookieName, ``, -1)
				} else {
					ids = append(ids, idByCookie)
				}
				hasCookieValue = true
			}
		}
		if len(o.HeaderName) > 0 {
			idByHeader := ctx.Header(o.HeaderName)
			if len(idByHeader) > 0 {
				idByHeader, err = o.IDDecoder(ctx, idByHeader, `header`)
				if err != nil {
					header.Add(`X-Captcha-ID-Error`, `header: `+err.Error())
				} else {
					ids = append(ids, idByHeader)
				}
				hasHeaderValue = true
			}
		}
		if len(ctx.Query("reload")) > 0 {
			var ok bool
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if captcha.Reload(id) {
					ok = true
					ids = []string{id}
					break
				}
			}
			if !ok && (hasCookieValue || hasHeaderValue) { // 旧的已经全部失效了，自动申请新ID
				id = captcha.New()
				ids = []string{id}
				if hasCookieValue {
					ctx.SetCookie(o.CookieName, id)
				}
				if hasHeaderValue {
					header.Set(o.HeaderName, id)
				}
			}
		}
		download := o.EnableDownload && ctx.Queryx("download").Bool()
		b := bytes.NewBuffer(nil)
		switch ext {
		case ".png":
			if !o.EnableImage {
				return echo.ErrNotFound
			}
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				err = captcha.WriteImage(b, id, captcha.StdWidth, captcha.StdHeight)
				if err == nil || err != captcha.ErrNotFound {
					break
				}
			}
			if err != nil {
				if err == captcha.ErrNotFound {
					return echo.ErrNotFound
				}
				return
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "image/png")
			}
		case ".wav":
			if !o.EnableAudio {
				return echo.ErrNotFound
			}
			lang := strings.ToLower(ctx.Query("lang"))
			supported := false
			for _, supportedLang := range o.AudioLangs {
				if supportedLang == lang {
					supported = true
					break
				}
			}
			if !supported && len(o.AudioLangs) > 0 {
				lang = o.AudioLangs[0]
			}
			var au *captcha.Audio
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				au, err = captcha.GetAudio(id, lang)
				if err == nil || err != captcha.ErrNotFound {
					break
				}
			}
			if err != nil {
				if err == captcha.ErrNotFound {
					return echo.ErrNotFound
				}
				return
			}
			length := strconv.Itoa(au.EncodedLen())
			_, err = au.WriteTo(b)
			if err != nil {
				return err
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "audio/x-wav")
			}
			header.Set("Content-Length", length)
		default:
			return nil
		}
		return ctx.Blob(b.Bytes())
	}
}
