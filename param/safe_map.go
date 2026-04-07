package param

import (
	"html/template"
	"time"

	"github.com/webx-top/com"
)

func NewMap() *SafeMap {
	return &SafeMap{
		SafeMap: com.NewSafeMap[any, any](),
	}
}

type SafeMap struct {
	*com.SafeMap[any, any]
}

func (s *SafeMap) Get(key any, defaults ...any) any {
	value, ok := s.GetOk(key)
	if (!ok || value == nil) && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() any); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func (s *SafeMap) Has(key any) bool {
	_, ok := s.GetOk(key)
	return ok
}

// GetOrSet get or set
//
//	value: `func() any` or `func() (any, bool)` or other
func (s *SafeMap) GetOrSet(key, value any) (actual any, loaded bool) {
	actual, loaded = s.GetOk(key)
	if loaded {
		return
	}
	switch f := value.(type) {
	case func() any:
		actual = f()
	case func() (any, bool):
		var store bool
		actual, store = f()
		if !store {
			return
		}
	default:
		actual = value
	}
	s.Set(key, actual)
	return
}

func (s *SafeMap) Load(key any) (actual any, loaded bool) {
	actual, loaded = s.GetOk(key)
	return
}

// LoadOrStore load or store
//
//	value: `func() any` or `func() (any, bool)` or other
func (s *SafeMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	actual, loaded = s.GetOrSet(key, value)
	return
}

func (s *SafeMap) Store(key, value any) {
	s.Set(key, value)
}

func (s *SafeMap) String(key any, defaults ...any) string {
	return AsString(s.Get(key, defaults...))
}

func (s *SafeMap) Split(key any, sep string, limit ...int) StringSlice {
	return Split(s.Get(key), sep, limit...)
}

func (s *SafeMap) Trim(key any, defaults ...any) String {
	return Trim(s.Get(key, defaults...))
}

func (s *SafeMap) HTML(key any, defaults ...any) template.HTML {
	return AsHTML(s.Get(key, defaults...))
}

func (s *SafeMap) HTMLAttr(key any, defaults ...any) template.HTMLAttr {
	return AsHTMLAttr(s.Get(key, defaults...))
}

func (s *SafeMap) JS(key any, defaults ...any) template.JS {
	return AsJS(s.Get(key, defaults...))
}

func (s *SafeMap) CSS(key any, defaults ...any) template.CSS {
	return AsCSS(s.Get(key, defaults...))
}

func (s *SafeMap) Bool(key any, defaults ...any) bool {
	return AsBool(s.Get(key, defaults...))
}

func (s *SafeMap) Float64(key any, defaults ...any) float64 {
	return AsFloat64(s.Get(key, defaults...))
}

func (s *SafeMap) Float32(key any, defaults ...any) float32 {
	return AsFloat32(s.Get(key, defaults...))
}

func (s *SafeMap) Int8(key any, defaults ...any) int8 {
	return AsInt8(s.Get(key, defaults...))
}

func (s *SafeMap) Int16(key any, defaults ...any) int16 {
	return AsInt16(s.Get(key, defaults...))
}

func (s *SafeMap) Int(key any, defaults ...any) int {
	return AsInt(s.Get(key, defaults...))
}

func (s *SafeMap) Int32(key any, defaults ...any) int32 {
	return AsInt32(s.Get(key, defaults...))
}

func (s *SafeMap) Int64(key any, defaults ...any) int64 {
	return AsInt64(s.Get(key, defaults...))
}

func (s *SafeMap) Decr(key any, n int64, defaults ...any) int64 {
	v := Decr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s *SafeMap) Incr(key any, n int64, defaults ...any) int64 {
	v := Incr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s *SafeMap) Uint8(key any, defaults ...any) uint8 {
	return AsUint8(s.Get(key, defaults...))
}

func (s *SafeMap) Uint16(key any, defaults ...any) uint16 {
	return AsUint16(s.Get(key, defaults...))
}

func (s *SafeMap) Uint(key any, defaults ...any) uint {
	return AsUint(s.Get(key, defaults...))
}

func (s *SafeMap) Uint32(key any, defaults ...any) uint32 {
	return AsUint32(s.Get(key, defaults...))
}

func (s *SafeMap) Uint64(key any, defaults ...any) uint64 {
	return AsUint64(s.Get(key, defaults...))
}

func (s *SafeMap) GetStore(key any, defaults ...any) Store {
	return AsStore(s.Get(key, defaults...))
}

func (s *SafeMap) Timestamp(key any, defaults ...any) time.Time {
	return AsTimestamp(s.Get(key, defaults...))
}

func (s *SafeMap) DateTime(key any, layouts ...string) time.Time {
	return AsDateTime(s.Get(key), layouts...)
}
