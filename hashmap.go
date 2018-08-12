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

package echo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Dump 输出对象和数组的结构信息
func Dump(m interface{}, args ...bool) (r string) {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	r = string(v)
	l := len(args)
	if l < 1 || args[0] {
		fmt.Println(r)
	}
	return
}

type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local == `H` {
		start.Name.Local = `Map`
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
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
func (h H) ToData() *RawData {
	var info, zone, data interface{}
	if v, y := h["Data"]; y {
		data = v
	}
	if v, y := h["Zone"]; y {
		zone = v
	}
	if v, y := h["Info"]; y {
		info = v
	}
	var code State
	if v, y := h["Code"]; y {
		if c, y := v.(int); y {
			code = State(c)
		} else if c, y := v.(State); y {
			code = c
		}
	}
	return &RawData{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
}

func (h H) DeepMerge(source H) {
	for k, value := range source {
		var (
			destValue interface{}
			ok        bool
		)
		if destValue, ok = h[k]; !ok {
			h[k] = value
			continue
		}
		sourceM, sourceOk := value.(H)
		destM, destOk := destValue.(H)
		if sourceOk && sourceOk == destOk {
			destM.DeepMerge(sourceM)
		} else {
			h[k] = value
		}
	}
}

func (h H) Clone() H {
	r := make(H)
	for k, value := range h {
		switch v := value.(type) {
		case H:
			r[k] = v.Clone()
		case []H:
			vCopy := make([]H, len(v))
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

type Mapx struct {
	Map   map[string]*Mapx `json:",omitempty"`
	Slice []*Mapx          `json:",omitempty"`
	Val   []string         `json:",omitempty"`
	lock  *sync.RWMutex
}

func NewMapx(data map[string][]string, mutex ...*sync.RWMutex) *Mapx {
	m := &Mapx{
		Map:   map[string]*Mapx{},
		Slice: []*Mapx{},
		Val:   []string{},
	}
	if len(mutex) > 0 {
		m.lock = mutex[0]
	} else {
		m.lock = &sync.RWMutex{}
	}
	if data == nil {
		return m
	}
	return m.Parse(data)
}

func (m *Mapx) Clone() *Mapx {
	mCopy := &Mapx{
		Map:   map[string]*Mapx{},
		Slice: make([]*Mapx, len(m.Slice)),
		Val:   make([]string, len(m.Val)),
		lock:  m.lock,
	}
	for key, mapx := range m.Map {
		mCopy.Map[key] = mapx.Clone()
	}
	for idx, mapx := range m.Slice {
		mCopy.Slice[idx] = mapx.Clone()
	}
	for idx, val := range m.Val {
		mCopy.Val[idx] = val
	}
	return mCopy
}

func (m *Mapx) Parse(data map[string][]string, keySkipper ...func(string) bool) *Mapx {
	m.lock.Lock()
	defer m.lock.Unlock()
	keys := make([]string, 0, len(data))
	var skip func(string) bool
	if len(keySkipper) > 0 {
		skip = keySkipper[0]
	}
	for k := range data {
		if skip != nil && skip(k) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, name := range keys {
		m.Add(name, data[name])
	}
	return m
}

func (m *Mapx) Add(name string, values []string) *Mapx {
	names := FormNames(name)
	end := len(names) - 1
	v := m
	for idx, key := range names {
		if len(key) == 0 {

			if v.Slice == nil {
				v.Slice = []*Mapx{}
			}

			if idx == end {
				v.Slice = append(v.Slice, &Mapx{lock: m.lock, Val: values})
				continue
			}
			mapx := &Mapx{
				lock: m.lock,
				Map:  map[string]*Mapx{},
			}
			v.Slice = append(v.Slice, mapx)
			v = mapx
			continue
		}
		if _, ok := v.Map[key]; !ok {
			if idx == end {
				v.Map[key] = &Mapx{lock: m.lock, Val: values}
				continue
			}
			v.Map[key] = &Mapx{
				lock: m.lock,
				Map:  map[string]*Mapx{},
			}
			v = v.Map[key]
			continue
		}

		if idx == end {
			v.Map[key] = &Mapx{lock: m.lock, Val: values}
		} else {
			v = v.Map[key]
		}
	}
	return m
}

func (m *Mapx) Value(names ...string) string {
	v := m.Values(names...)
	if v != nil {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ``
}

func (m *Mapx) ValueOk(names ...string) (string, bool) {
	v, y := m.ValuesOk(names...)
	if y && v != nil {
		if len(v) > 0 {
			return v[0], true
		}
	}
	return ``, false
}

func (m *Mapx) ValuesOk(names ...string) ([]string, bool) {
	if len(names) == 0 {
		if m.Val == nil {
			return []string{}, false
		}
		return m.Val, true
	}
	v := m.Get(names...)
	if v != nil {
		return v.Val, true
	}
	return []string{}, false
}

func (m *Mapx) Values(names ...string) []string {
	if len(names) == 0 {
		if m.Val == nil {
			return []string{}
		}
		return m.Val
	}
	v := m.Get(names...)
	if v != nil {
		return v.Val
	}
	return []string{}
}

func (m *Mapx) get(k string) (*Mapx, bool) {
	r, y := m.Map[k]
	return r, y
}

func (m *Mapx) Get(names ...string) *Mapx {
	m.lock.Lock()
	defer m.lock.Unlock()
	v := m
	end := len(names) - 1
	for idx, key := range names {
		_, ok := v.get(key)
		if !ok {
			if v.Slice == nil {
				return nil
			}
			i, err := strconv.Atoi(key)
			if err != nil {
				return nil
			}
			if i < 0 {
				return nil
			}
			if i < len(v.Slice) {
				v = v.Slice[i]
				continue
			}
			return nil
		}
		v, _ = v.get(key)

		if idx == end {
			return v
		}
	}
	return nil
}
