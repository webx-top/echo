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

func (s *SafeMap) Get(key interface{}, defaults ...interface{}) interface{} {
	value, ok := s.GetOk(key)
	if (!ok || value == nil) && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() interface{}); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func (s *SafeMap) Has(key interface{}) bool {
	_, ok := s.GetOk(key)
	return ok
}

// GetOrSet get or set
//
//	value: `func() interface{}` or `func() (interface{}, bool)` or other
func (s *SafeMap) GetOrSet(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.GetOk(key)
	if loaded {
		return
	}
	switch f := value.(type) {
	case func() interface{}:
		actual = f()
	case func() (interface{}, bool):
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

func (s *SafeMap) Load(key interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.GetOk(key)
	return
}

// LoadOrStore load or store
//
//	value: `func() interface{}` or `func() (interface{}, bool)` or other
func (s *SafeMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.GetOrSet(key, value)
	return
}

func (s *SafeMap) Store(key, value interface{}) {
	s.Set(key, value)
}

func (s *SafeMap) String(key interface{}, defaults ...interface{}) string {
	return AsString(s.Get(key, defaults...))
}

func (s *SafeMap) Split(key interface{}, sep string, limit ...int) StringSlice {
	return Split(s.Get(key), sep, limit...)
}

func (s *SafeMap) Trim(key interface{}, defaults ...interface{}) String {
	return Trim(s.Get(key, defaults...))
}

func (s *SafeMap) HTML(key interface{}, defaults ...interface{}) template.HTML {
	return AsHTML(s.Get(key, defaults...))
}

func (s *SafeMap) HTMLAttr(key interface{}, defaults ...interface{}) template.HTMLAttr {
	return AsHTMLAttr(s.Get(key, defaults...))
}

func (s *SafeMap) JS(key interface{}, defaults ...interface{}) template.JS {
	return AsJS(s.Get(key, defaults...))
}

func (s *SafeMap) CSS(key interface{}, defaults ...interface{}) template.CSS {
	return AsCSS(s.Get(key, defaults...))
}

func (s *SafeMap) Bool(key interface{}, defaults ...interface{}) bool {
	return AsBool(s.Get(key, defaults...))
}

func (s *SafeMap) Float64(key interface{}, defaults ...interface{}) float64 {
	return AsFloat64(s.Get(key, defaults...))
}

func (s *SafeMap) Float32(key interface{}, defaults ...interface{}) float32 {
	return AsFloat32(s.Get(key, defaults...))
}

func (s *SafeMap) Int8(key interface{}, defaults ...interface{}) int8 {
	return AsInt8(s.Get(key, defaults...))
}

func (s *SafeMap) Int16(key interface{}, defaults ...interface{}) int16 {
	return AsInt16(s.Get(key, defaults...))
}

func (s *SafeMap) Int(key interface{}, defaults ...interface{}) int {
	return AsInt(s.Get(key, defaults...))
}

func (s *SafeMap) Int32(key interface{}, defaults ...interface{}) int32 {
	return AsInt32(s.Get(key, defaults...))
}

func (s *SafeMap) Int64(key interface{}, defaults ...interface{}) int64 {
	return AsInt64(s.Get(key, defaults...))
}

func (s *SafeMap) Decr(key interface{}, n int64, defaults ...interface{}) int64 {
	v := Decr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s *SafeMap) Incr(key interface{}, n int64, defaults ...interface{}) int64 {
	v := Incr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s *SafeMap) Uint8(key interface{}, defaults ...interface{}) uint8 {
	return AsUint8(s.Get(key, defaults...))
}

func (s *SafeMap) Uint16(key interface{}, defaults ...interface{}) uint16 {
	return AsUint16(s.Get(key, defaults...))
}

func (s *SafeMap) Uint(key interface{}, defaults ...interface{}) uint {
	return AsUint(s.Get(key, defaults...))
}

func (s *SafeMap) Uint32(key interface{}, defaults ...interface{}) uint32 {
	return AsUint32(s.Get(key, defaults...))
}

func (s *SafeMap) Uint64(key interface{}, defaults ...interface{}) uint64 {
	return AsUint64(s.Get(key, defaults...))
}

func (s *SafeMap) GetStore(key interface{}, defaults ...interface{}) Store {
	return AsStore(s.Get(key, defaults...))
}

func (s *SafeMap) Timestamp(key interface{}, defaults ...interface{}) time.Time {
	return AsTimestamp(s.Get(key, defaults...))
}

func (s *SafeMap) DateTime(key interface{}, layouts ...string) time.Time {
	return AsDateTime(s.Get(key), layouts...)
}
