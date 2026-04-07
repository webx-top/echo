package param

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"time"
)

type MapReadonly map[string]any

func (s MapReadonly) Has(key string) bool {
	_, y := s[key]
	return y
}

func (s MapReadonly) Get(key string, defaults ...any) any {
	value, ok := s[key]
	if (!ok || value == nil) && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() any); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func (s MapReadonly) String(key string, defaults ...any) string {
	return AsString(s.Get(key, defaults...))
}

func (s MapReadonly) Split(key string, sep string, limit ...int) StringSlice {
	return Split(s.Get(key), sep, limit...)
}

func (s MapReadonly) Trim(key string, defaults ...any) String {
	return Trim(s.Get(key, defaults...))
}

func (s MapReadonly) HTML(key string, defaults ...any) template.HTML {
	return AsHTML(s.Get(key, defaults...))
}

func (s MapReadonly) HTMLAttr(key string, defaults ...any) template.HTMLAttr {
	return AsHTMLAttr(s.Get(key, defaults...))
}

func (s MapReadonly) JS(key string, defaults ...any) template.JS {
	return AsJS(s.Get(key, defaults...))
}

func (s MapReadonly) CSS(key string, defaults ...any) template.CSS {
	return AsCSS(s.Get(key, defaults...))
}

func (s MapReadonly) Bool(key string, defaults ...any) bool {
	return AsBool(s.Get(key, defaults...))
}

func (s MapReadonly) Float64(key string, defaults ...any) float64 {
	return AsFloat64(s.Get(key, defaults...))
}

func (s MapReadonly) Float32(key string, defaults ...any) float32 {
	return AsFloat32(s.Get(key, defaults...))
}

func (s MapReadonly) Int8(key string, defaults ...any) int8 {
	return AsInt8(s.Get(key, defaults...))
}

func (s MapReadonly) Int16(key string, defaults ...any) int16 {
	return AsInt16(s.Get(key, defaults...))
}

func (s MapReadonly) Int(key string, defaults ...any) int {
	return AsInt(s.Get(key, defaults...))
}

func (s MapReadonly) Int32(key string, defaults ...any) int32 {
	return AsInt32(s.Get(key, defaults...))
}

func (s MapReadonly) Int64(key string, defaults ...any) int64 {
	return AsInt64(s.Get(key, defaults...))
}

func (s MapReadonly) Uint8(key string, defaults ...any) uint8 {
	return AsUint8(s.Get(key, defaults...))
}

func (s MapReadonly) Uint16(key string, defaults ...any) uint16 {
	return AsUint16(s.Get(key, defaults...))
}

func (s MapReadonly) Uint(key string, defaults ...any) uint {
	return AsUint(s.Get(key, defaults...))
}

func (s MapReadonly) Uint32(key string, defaults ...any) uint32 {
	return AsUint32(s.Get(key, defaults...))
}

func (s MapReadonly) Uint64(key string, defaults ...any) uint64 {
	return AsUint64(s.Get(key, defaults...))
}

func (s MapReadonly) Timestamp(key string, defaults ...any) time.Time {
	return AsTimestamp(s.Get(key, defaults...))
}

func (s MapReadonly) Duration(key string, defaults ...time.Duration) time.Duration {
	return AsDuration(s.Get(key), defaults...)
}

func (s MapReadonly) DateTime(key string, layouts ...string) time.Time {
	return AsDateTime(s.Get(key), layouts...)
}

func (s MapReadonly) Children(keys ...any) MapReadonly {
	r := s
	for _, key := range keys {
		r = r.GetStore(fmt.Sprint(key))
	}
	return r
}

func (s MapReadonly) GetStore(key string, defaults ...any) MapReadonly {
	return MapReadonly(AsStore(s.Get(key, defaults...)))
}

func (s MapReadonly) GetStoreByKeys(keys ...string) MapReadonly {
	sz := len(keys)
	if sz == 0 {
		return s
	}
	r := s.GetStore(keys[0])
	if sz == 1 {
		return r
	}
	for _, key := range keys[1:] {
		r = r.GetStore(key)
	}
	return r
}

func (s MapReadonly) Select(selectKeys ...string) MapReadonly {
	r := MapReadonly{}
	for _, key := range selectKeys {
		r[key] = s.Get(key)
	}
	return r
}

// MarshalXML allows type MapReadonly to be used with xml.Marshal
func (s MapReadonly) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local == `MapReadonly` {
		start.Name.Local = `Map`
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range s {
		elem := xml.StartElement{
			Name: xml.Name{Space: ``, Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
