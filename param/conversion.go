package param

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/dateparse"
	"github.com/webx-top/com"
)

const (
	EmptyHTML      = template.HTML(``)
	EmptyJS        = template.JS(``)
	EmptyCSS       = template.CSS(``)
	EmptyHTMLAttr  = template.HTMLAttr(``)
	DateTimeNormal = `2006-01-02 15:04:05`
	DateTimeShort  = `2006-01-02 15:04`
	DateNormal     = `2006-01-02`
	TimeNormal     = `15:04:05`
	DateMd         = `01-02`
	DateShort      = `06-01-02`
	TimeShort      = `15:04`
)

func AsType(typ string, val interface{}) interface{} {
	return com.AsType(typ, val)
}

func AsString(val interface{}) string {
	return com.ToStr(val)
}

func AsBytes(val interface{}) []byte {
	return com.Bytes(val)
}

func Split(val interface{}, sep string, limit ...int) StringSlice {
	str := AsString(val)
	if len(str) == 0 {
		return StringSlice{}
	}
	if len(limit) > 0 {
		return strings.SplitN(str, sep, limit[0])
	}
	return strings.Split(str, sep)
}

func Trim(val interface{}) String {
	return String(strings.TrimSpace(AsString(val)))
}

func AsHTML(val interface{}) template.HTML {
	switch v := val.(type) {
	case template.HTML:
		return v
	case string:
		return template.HTML(v)
	case nil:
		return EmptyHTML
	default:
		return template.HTML(fmt.Sprint(v))
	}
}

func AsHTMLAttr(val interface{}) template.HTMLAttr {
	switch v := val.(type) {
	case template.HTMLAttr:
		return v
	case string:
		return template.HTMLAttr(v)
	case nil:
		return EmptyHTMLAttr
	default:
		return template.HTMLAttr(fmt.Sprint(v))
	}
}

func AsJS(val interface{}) template.JS {
	switch v := val.(type) {
	case template.JS:
		return v
	case string:
		return template.JS(v)
	case nil:
		return EmptyJS
	default:
		return template.JS(fmt.Sprint(v))
	}
}

func AsCSS(val interface{}) template.CSS {
	switch v := val.(type) {
	case template.CSS:
		return v
	case string:
		return template.CSS(v)
	case nil:
		return EmptyCSS
	default:
		return template.CSS(fmt.Sprint(v))
	}
}

func AsBool(val interface{}) bool {
	return com.Bool(val)
}

func AsFloat64(val interface{}) float64 {
	return com.Float64(val)
}

func AsFloat32(val interface{}) float32 {
	return com.Float32(val)
}

func AsInt8(val interface{}) int8 {
	return com.Int8(val)
}

func AsInt16(val interface{}) int16 {
	return com.Int16(val)
}

func AsInt(val interface{}) int {
	return com.Int(val)
}

func AsInt32(val interface{}) int32 {
	return com.Int32(val)
}

func AsInt64(val interface{}) int64 {
	return com.Int64(val)
}

func Decr(val interface{}, n int64) int64 {
	v := AsInt64(val)
	v -= n
	return v
}

func Incr(val interface{}, n int64) int64 {
	v := AsInt64(val)
	v += n
	return v
}

func AsUint8(val interface{}) uint8 {
	return com.Uint8(val)
}

func AsUint16(val interface{}) uint16 {
	return com.Uint16(val)
}

func AsUint(val interface{}) uint {
	return com.Uint(val)
}

func AsUint32(val interface{}) uint32 {
	return com.Uint32(val)
}

func AsUint64(val interface{}) uint64 {
	return com.Uint64(val)
}

func AsTimestamp(val interface{}) time.Time {
	p := AsString(val)
	if len(p) > 0 {
		s := strings.SplitN(p, `.`, 2)
		var sec int64
		var nsec int64
		switch len(s) {
		case 2:
			nsec = String(s[1]).Int64()
			fallthrough
		case 1:
			sec = String(s[0]).Int64()
		}
		return time.Unix(sec, nsec)
	}
	return EmptyTime
}

func AsDateTime(val interface{}, layouts ...string) time.Time {
	p := AsString(val)
	if len(p) > 0 {
		layout := DateTimeNormal
		if len(layouts) > 0 && len(layouts[0]) > 0 {
			layout = layouts[0]
		}
		t, err := time.ParseInLocation(layout, p, time.Local)
		if err != nil {
			t, _ = dateparse.ParseLocal(p)
		}
		return t
	}
	return EmptyTime
}

func AsDuration(val interface{}, defaults ...time.Duration) time.Duration {
	p := AsString(val)
	if len(p) > 0 {
		t, err := time.ParseDuration(p)
		if err == nil {
			return t
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

func AsStore(val interface{}) Store {
	switch v := val.(type) {
	case Store:
		return v
	case map[string]interface{}:
		return Store(v)
	case map[string]uint64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]uint:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]uint32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]float32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]float64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]string:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	default:
		return emptyStore
	}
}

// AsStdStringSlice p must be slice
func AsStdStringSlice(p interface{}) []string {
	var r []string
	switch v := p.(type) {
	case []uint64:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatUint(a, 10)
		}
		return r
	case []int64:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatInt(a, 10)
		}
		return r
	case []uint:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatUint(uint64(a), 10)
		}
		return r
	case []int:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.Itoa(a)
		}
		return r
	case []uint32:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatUint(uint64(a), 10)
		}
		return r
	case []int32:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatInt(int64(a), 10)
		}
		return r
	case []float32:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatFloat(float64(a), 'f', -1, 32)
		}
		return r
	case []float64:
		r = make([]string, len(v))
		for k, a := range v {
			r[k] = strconv.FormatFloat(a, 'f', -1, 64)
		}
		return r
	case []string:
		return v
	default:
		panic(fmt.Sprintf(`[AsStdStringSlice] unsupported: %T`, p))
	}
}

// Converts 转换 Slice (不支持指针类型元素)
func Converts[V Scalar, T Scalar](p []T, converter ...func(s T) V) []V {
	result := make([]V, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s T) V
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		rv := reflect.ValueOf(result[0])
		typeName := rv.Kind().String()
		convert = func(s T) V {
			return AsType(typeName, s).(V)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}

func InterfacesTo[T Scalar](p []any, converter ...func(s any) T) []T {
	result := make([]T, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s any) T
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		rv := reflect.ValueOf(result[0])
		typeName := rv.Kind().String()
		convert = func(s any) T {
			return AsType(typeName, s).(T)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}

func AsInterfaces[T any](p []T, converter ...func(s T) any) []any {
	result := make([]any, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s T) any
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		convert = func(s T) any {
			return any(s)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}
