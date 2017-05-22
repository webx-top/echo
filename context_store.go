package echo

import (
	"html/template"
	"sync"
)

var (
	mutex         sync.RWMutex
	emptyHTML     = template.HTML(``)
	emptyJS       = template.JS(``)
	emptyCSS      = template.CSS(``)
	emptyHTMLAttr = template.HTMLAttr(``)
)

type store map[string]interface{}

func (s store) Set(key string, value interface{}) store {
	mutex.Lock()
	s[key] = value
	mutex.Unlock()
	return s
}

func (s store) Get(key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if v, y := s[key]; y {
		return v
	}
	return nil
}

func (s store) String(key string) string {
	if v, y := s.Get(key).(string); y {
		return v
	}
	return ``
}

func (s store) HTML(key string) template.HTML {
	val := s.Get(key)
	if v, y := val.(string); y {
		return template.HTML(v)
	}
	if v, y := val.(template.HTML); y {
		return v
	}
	return emptyHTML
}

func (s store) HTMLAttr(key string) template.HTMLAttr {
	val := s.Get(key)
	if v, y := val.(string); y {
		return template.HTMLAttr(v)
	}
	if v, y := val.(template.HTMLAttr); y {
		return v
	}
	return emptyHTMLAttr
}

func (s store) JS(key string) template.JS {
	val := s.Get(key)
	if v, y := val.(string); y {
		return template.JS(v)
	}
	if v, y := val.(template.JS); y {
		return v
	}
	return emptyJS
}

func (s store) CSS(key string) template.CSS {
	val := s.Get(key)
	if v, y := val.(string); y {
		return template.CSS(v)
	}
	if v, y := val.(template.CSS); y {
		return v
	}
	return emptyCSS
}

func (s store) Bool(key string) bool {
	if v, y := s.Get(key).(bool); y {
		return v
	}
	return false
}

func (s store) Float64(key string) float64 {
	if v, y := s.Get(key).(float64); y {
		return v
	}
	return 0
}

func (s store) Float32(key string) float32 {
	if v, y := s.Get(key).(float32); y {
		return v
	}
	return 0
}

func (s store) Int8(key string) int8 {
	if v, y := s.Get(key).(int8); y {
		return v
	}
	return 0
}

func (s store) Int16(key string) int16 {
	if v, y := s.Get(key).(int16); y {
		return v
	}
	return 0
}

func (s store) Int(key string) int {
	if v, y := s.Get(key).(int); y {
		return v
	}
	return 0
}

func (s store) Int32(key string) int32 {
	if v, y := s.Get(key).(int32); y {
		return v
	}
	return 0
}

func (s store) Int64(key string) int64 {
	if v, y := s.Get(key).(int64); y {
		return v
	}
	return 0
}

func (s store) Decr(key string, n int64) int64 {
	v, _ := s.Get(key).(int64)
	v -= n
	s.Set(key, v)
	return v
}

func (s store) Incr(key string, n int64) int64 {
	v, _ := s.Get(key).(int64)
	v += n
	s.Set(key, v)
	return v
}

func (s store) Uint8(key string) uint8 {
	if v, y := s.Get(key).(uint8); y {
		return v
	}
	return 0
}

func (s store) Uint16(key string) uint16 {
	if v, y := s.Get(key).(uint16); y {
		return v
	}
	return 0
}

func (s store) Uint(key string) uint {
	if v, y := s.Get(key).(uint); y {
		return v
	}
	return 0
}

func (s store) Uint32(key string) uint32 {
	if v, y := s.Get(key).(uint32); y {
		return v
	}
	return 0
}

func (s store) Uint64(key string) uint64 {
	if v, y := s.Get(key).(uint64); y {
		return v
	}
	return 0
}

func (s store) Delete(keys ...string) {
	mutex.Lock()
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	mutex.Unlock()
}
