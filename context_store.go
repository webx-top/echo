package echo

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/webx-top/echo/param"
)

var (
	mutex         sync.RWMutex
	emptyHTML     = template.HTML(``)
	emptyJS       = template.JS(``)
	emptyCSS      = template.CSS(``)
	emptyHTMLAttr = template.HTMLAttr(``)
	emptyStore    = Store{}
)

type Store map[string]interface{}

func (s Store) Set(key string, value interface{}) Store {
	mutex.Lock()
	s[key] = value
	mutex.Unlock()
	return s
}

func (s Store) Has(key string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, y := s[key]
	return y
}

func (s Store) Get(key string, defaults ...interface{}) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if v, y := s[key]; y {
		if v == nil && len(defaults) > 0 {
			return defaults[0]
		}
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

func (s Store) String(key string, defaults ...interface{}) string {
	p := s.Get(key, defaults...)
	if v, y := p.(string); y {
		return v
	}
	return fmt.Sprint(p)
}

func (s Store) Split(key string, sep string, limit ...int) param.StringSlice {
	str := s.String(key)
	if len(str) == 0 {
		return param.StringSlice{}
	}
	if len(limit) > 0 {
		return strings.SplitN(str, sep, limit[0])
	}
	return strings.Split(str, sep)
}

func (s Store) Trim(key string, defaults ...interface{}) param.String {
	return param.String(strings.TrimSpace(s.String(key, defaults...)))
}

func (s Store) HTML(key string, defaults ...interface{}) template.HTML {
	val := s.Get(key, defaults...)
	if v, y := val.(template.HTML); y {
		return v
	}
	if v, y := val.(string); y {
		return template.HTML(v)
	}
	return emptyHTML
}

func (s Store) HTMLAttr(key string, defaults ...interface{}) template.HTMLAttr {
	val := s.Get(key, defaults...)
	if v, y := val.(template.HTMLAttr); y {
		return v
	}
	if v, y := val.(string); y {
		return template.HTMLAttr(v)
	}
	return emptyHTMLAttr
}

func (s Store) JS(key string, defaults ...interface{}) template.JS {
	val := s.Get(key, defaults...)
	if v, y := val.(template.JS); y {
		return v
	}
	if v, y := val.(string); y {
		return template.JS(v)
	}
	return emptyJS
}

func (s Store) CSS(key string, defaults ...interface{}) template.CSS {
	val := s.Get(key, defaults...)
	if v, y := val.(template.CSS); y {
		return v
	}
	if v, y := val.(string); y {
		return template.CSS(v)
	}
	return emptyCSS
}

func (s Store) Bool(key string, defaults ...interface{}) bool {
	if v, y := s.Get(key, defaults...).(bool); y {
		return v
	}
	return false
}

func (s Store) Float64(key string, defaults ...interface{}) float64 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case int32:
		return float64(v)
	case uint32:
		return float64(v)
	case int:
		return float64(v)
	case uint:
		return float64(v)
	case string:
		i, _ := strconv.ParseFloat(v, 64)
		return i
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseFloat(s, 64)
		return i
	}
}

func (s Store) Float32(key string, defaults ...interface{}) float32 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case float32:
		return v
	case int32:
		return float32(v)
	case uint32:
		return float32(v)
	case string:
		f, _ := strconv.ParseFloat(v, 32)
		return float32(f)
	default:
		s := fmt.Sprint(val)
		f, _ := strconv.ParseFloat(s, 32)
		return float32(f)
	}
}

func (s Store) Int8(key string, defaults ...interface{}) int8 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case int8:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 8)
		return int8(i)
	default:
		s := fmt.Sprint(val)
		i, _ := strconv.ParseInt(s, 10, 8)
		return int8(i)
	}
}

func (s Store) Int16(key string, defaults ...interface{}) int16 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case int16:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 16)
		return int16(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 16)
		return int16(i)
	}
}

func (s Store) Int(key string, defaults ...interface{}) int {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case int:
		return v
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.Atoi(s)
		return i
	}
}

func (s Store) Int32(key string, defaults ...interface{}) int32 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case int32:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 32)
		return int32(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 32)
		return int32(i)
	}
}

func (s Store) Int64(key string, defaults ...interface{}) int64 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 64)
		return i
	}
}

func (s Store) Decr(key string, n int64, defaults ...interface{}) int64 {
	v, _ := s.Get(key, defaults...).(int64)
	v -= n
	s.Set(key, v)
	return v
}

func (s Store) Incr(key string, n int64, defaults ...interface{}) int64 {
	v, _ := s.Get(key, defaults...).(int64)
	v += n
	s.Set(key, v)
	return v
}

func (s Store) Uint8(key string, defaults ...interface{}) uint8 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case uint8:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 8)
		return uint8(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 8)
		return uint8(i)
	}
}

func (s Store) Uint16(key string, defaults ...interface{}) uint16 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case uint16:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 16)
		return uint16(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 16)
		return uint16(i)
	}
}

func (s Store) Uint(key string, defaults ...interface{}) uint {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case uint:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 32)
		return uint(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 32)
		return uint(i)
	}
}

func (s Store) Uint32(key string, defaults ...interface{}) uint32 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case uint32:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 32)
		return uint32(i)
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 32)
		return uint32(i)
	}
}

func (s Store) Uint64(key string, defaults ...interface{}) uint64 {
	val := s.Get(key, defaults...)
	switch v := val.(type) {
	case uint64:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 64)
		return i
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 64)
		return i
	}
}

func (s Store) Timestamp(key string, defaults ...interface{}) time.Time {
	p := s.String(key, defaults...)
	if len(p) > 0 {
		s := strings.SplitN(p, `.`, 2)
		var sec int64
		var nsec int64
		switch len(s) {
		case 2:
			nsec = param.String(s[1]).Int64()
			fallthrough
		case 1:
			sec = param.String(s[0]).Int64()
		}
		return time.Unix(sec, nsec)
	}
	return time.Time{}
}

func (s Store) DateTime(key string, layouts ...string) time.Time {
	p := s.String(key)
	if len(p) > 0 {
		layout := `2006-01-02 15:04:05`
		if len(layouts) > 0 {
			layout = layouts[0]
		}
		t, _ := time.Parse(layout, p)
		return t
	}
	return time.Time{}
}

func (s Store) Store(key string, defaults ...interface{}) Store {
	val := s.Get(key, defaults...)
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

func (s Store) Delete(keys ...string) {
	mutex.Lock()
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	mutex.Unlock()
}

// MarshalXML allows type Store to be used with xml.Marshal
func (s Store) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local == `Store` {
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

// ToData conversion to *RawData
func (s Store) ToData() *RawData {
	var info, zone, data interface{}
	if v, y := s["Data"]; y {
		data = v
	}
	if v, y := s["Zone"]; y {
		zone = v
	}
	if v, y := s["Info"]; y {
		info = v
	}
	var code State
	if v, y := s["Code"]; y {
		switch c := v.(type) {
		case State:
			code = c
		case int:
			code = State(c)
		case string:
			i, _ := strconv.Atoi(c)
			code = State(i)
		default:
			s := fmt.Sprint(c)
			i, _ := strconv.Atoi(s)
			code = State(i)
		}
	}
	return &RawData{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
}

func (s Store) DeepMerge(source Store) {
	for k, value := range source {
		var (
			destValue interface{}
			ok        bool
		)
		if destValue, ok = s[k]; !ok {
			s[k] = value
			continue
		}
		sourceM, sourceOk := value.(H)
		destM, destOk := destValue.(H)
		if sourceOk && sourceOk == destOk {
			destM.DeepMerge(sourceM)
		} else {
			s[k] = value
		}
	}
}

func (s Store) Clone() Store {
	r := make(Store)
	for k, value := range s {
		switch v := value.(type) {
		case Store:
			r[k] = v.Clone()
		case []Store:
			vCopy := make([]Store, len(v))
			for i, row := range v {
				vCopy[i] = row.Clone()
			}
			r[k] = vCopy
		default:
			r[k] = value
		}
	}
	return r
}
