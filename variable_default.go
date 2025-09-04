/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

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

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

var (
	DefaultAcceptFormats = map[string]string{
		//json
		MIMEApplicationJSON:       ContentTypeJSON,
		`text/javascript`:         ContentTypeJSON,
		MIMEApplicationJavaScript: ContentTypeJSON,

		//xml
		MIMEApplicationXML: ContentTypeXML,
		`text/xml`:         ContentTypeXML,

		//text
		MIMETextPlain: ContentTypeText,

		//html
		`*/*`:               ContentTypeHTML,
		`application/xhtml`: ContentTypeHTML,
		MIMETextHTML:        ContentTypeHTML,

		//default
		`*`: ContentTypeHTML,
	}
	DefaultFormatRenderers = map[string]FormatRender{
		ContentTypeJSON: func(c Context, data interface{}, code ...int) error {
			return c.JSON(c.Data(), code...)
		},
		ContentTypeJSONP: func(c Context, data interface{}, code ...int) error {
			return c.JSONP(c.Query(c.Echo().JSONPVarName), c.Data(), code...)
		},
		ContentTypeXML: func(c Context, data interface{}, code ...int) error {
			return c.XML(c.Data(), code...)
		},
		ContentTypeText: func(c Context, data interface{}, code ...int) error {
			return c.String(fmt.Sprint(data), code...)
		},
	}
	DefaultBinderDecoders = map[string]func(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error{
		MIMEApplicationJSON: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			body := ctx.Request().Body()
			if body == nil {
				return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
			}
			defer body.Close()
			err := json.NewDecoder(body).Decode(i)
			if err != nil {
				switch ev := err.(type) {
				case *json.UnmarshalTypeError:
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ev.Type, ev.Value, ev.Field, ev.Offset)).SetRaw(err)
				case *json.SyntaxError:
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", ev.Offset, ev.Error())).SetRaw(err)
				}
				return NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}
			return err
		},
		MIMEApplicationXML: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			body := ctx.Request().Body()
			if body == nil {
				return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
			}
			defer body.Close()
			err := xml.NewDecoder(body).Decode(i)
			if err != nil {
				switch ev := err.(type) {
				case *xml.UnsupportedTypeError:
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ev.Type, ev.Error())).SetRaw(err)
				case *xml.SyntaxError:
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", ev.Line, ev.Error())).SetRaw(err)
				}
				return NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}
			return err
		},
		MIMEApplicationForm: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().PostForm().All(), ``, valueDecoders, filter...)
		},
		MIMEMultipartForm: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			_, err := ctx.Request().MultipartForm()
			if err != nil && !errors.Is(err, http.ErrMissingBoundary) {
				return err
			}
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().Form().All(), ``, valueDecoders, filter...)
		},
		`*`: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().Form().All(), ``, valueDecoders, filter...)
		},
	}
	// DefaultHTMLFilter html filter (`form_filter:"html"`)
	DefaultHTMLFilter = func(v string) (r string) {
		return v
	}
	DefaultBinderValueEncoders = map[string]BinderValueEncoder{
		`joinKVRows`:     binderValueEncoderJoinKVRows,
		`join`:           binderValueEncoderJoin,
		`unix2time`:      binderValueEncoderUnix2time,      // int64
		`unixmilli2time`: binderValueEncoderUnixmilli2time, // int64
		`unixmicro2time`: binderValueEncoderUnixmicro2time, // int64
		`unixnano2time`:  binderValueEncoderUnixnano2time,  // float64
	}
	DefaultBinderValueDecoders = map[string]BinderValueDecoder{
		`splitKVRows`:    binderValueDecoderSplitKVRows,
		`split`:          binderValueDecoderSplit,
		`time2unix`:      binderValueDecoderTime2unix,      // int64
		`time2unixmilli`: binderValueDecoderTime2unixmilli, // int64
		`time2unixmicro`: binderValueDecoderTime2unixmicro, // int64
		`time2unixnano`:  binderValueDecoderTime2unixnano,  // float64
	}
)

func binderValueDecoderSplitKVRows(field string, values []string, seperator string) (interface{}, error) {
	return com.SplitKVRows(values[0], seperator), nil
}

func binderValueDecoderSplit(field string, values []string, seperator string) (interface{}, error) {
	return strings.Split(values[0], seperator), nil
}

func binderValueEncoderJoin(field string, value interface{}, seperator string) []string {
	return []string{strings.Join(param.AsStdStringSlice(value), seperator)}
}

func binderValueEncoderJoinKVRows(field string, value interface{}, seperator string) []string {
	result := com.JoinKVRows(value, seperator)
	if len(result) == 0 {
		return nil
	}
	return []string{result}
}

var TimeLayouts = map[string]string{
	`datetime`:    time.DateTime,
	`dateonly`:    time.DateOnly,
	`timeonly`:    time.TimeOnly,
	`stampnano`:   time.StampNano,
	`layout`:      time.Layout,
	`ansic`:       time.ANSIC,
	`unixdate`:    time.UnixDate,
	`rubydate`:    time.RubyDate,
	`rfc822`:      time.RFC822,
	`rfc822z`:     time.RFC822Z,
	`rfc850`:      time.RFC850,
	`rfc1123`:     time.RFC1123,
	`rfc1123z`:    time.RFC1123Z,
	`rfc3339`:     time.RFC3339,
	`rfc3339nano`: time.RFC3339Nano,
	`kitchen`:     time.Kitchen,
	`stamp`:       time.Stamp,
	`stampmilli`:  time.StampMilli,
	`stampmicro`:  time.StampMicro,
	// custom
	`datetimeshort`: param.DateTimeShort,
	`timeshort`:     param.TimeShort,
	`dateshort`:     param.DateShort,
	`monthday`:      param.DateMd,
	`hourminute`:    param.TimeShort,
}

func GetTimeLayoutByName(name string) string {
	v, y := TimeLayouts[strings.ToLower(name)]
	if y {
		return v
	}
	return name
}

func binderValueEncoderUnix2time(field string, value interface{}, layout string) []string {
	t := param.AsTimestamp(value)
	if t.IsZero() {
		return []string{}
	}
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return []string{t.Format(layout)}
}

func binderValueDecoderTime2unix(field string, values []string, layout string) (interface{}, error) {
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return param.AsDateTime(values[0], layout).Unix(), nil
}

func binderValueEncoderUnixmilli2time(field string, value interface{}, layout string) []string {
	t := time.UnixMilli(param.AsInt64(value))
	if t.IsZero() {
		return []string{}
	}
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return []string{t.Format(layout)}
}

func binderValueDecoderTime2unixmilli(field string, values []string, layout string) (interface{}, error) {
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return param.AsDateTime(values[0], layout).UnixMilli(), nil
}

func binderValueEncoderUnixmicro2time(field string, value interface{}, layout string) []string {
	t := time.UnixMicro(param.AsInt64(value))
	if t.IsZero() {
		return []string{}
	}
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return []string{t.Format(layout)}
}

func binderValueDecoderTime2unixmicro(field string, values []string, layout string) (interface{}, error) {
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return param.AsDateTime(values[0], layout).UnixMicro(), nil
}

func binderValueEncoderUnixnano2time(field string, value interface{}, layout string) []string {
	t := param.AsTimestamp(value)
	if t.IsZero() {
		return []string{}
	}
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	return []string{t.Format(layout)}
}

func binderValueDecoderTime2unixnano(field string, values []string, layout string) (interface{}, error) {
	if len(layout) == 0 {
		layout = param.DateTimeNormal
	} else {
		layout = GetTimeLayoutByName(layout)
	}
	t := param.AsDateTime(values[0], layout)
	v := fmt.Sprintf(`%d.%d`, t.Unix(), t.Nanosecond())
	return v, nil
}
