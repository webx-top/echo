package echo

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
)

func NewKV(k, v string) *KV {
	return &KV{K: k, V: v}
}

// KV 键值对
type KV struct {
	K        string
	V        string
	H        H           `json:",omitempty" xml:",omitempty"`
	X        interface{} `json:",omitempty" xml:",omitempty"`
	fn       func(context.Context) interface{}
	priority int
}

func (a *KV) SetPriority(priority int) *KV {
	a.priority = priority
	return a
}

func (a *KV) Clone() KV {
	return KV{
		K:        a.K,
		V:        a.V,
		H:        a.H.Clone(),
		X:        a.X,
		fn:       a.fn,
		priority: a.priority,
	}
}

func (a *KV) SetK(k string) *KV {
	a.K = k
	return a
}

func (a *KV) SetV(v string) *KV {
	a.V = v
	return a
}

func (a *KV) SetKV(k, v string) *KV {
	a.K = k
	a.V = v
	return a
}

func (a *KV) SetH(h H) *KV {
	a.H = h
	return a
}

func (a *KV) SetHKV(k string, v interface{}) *KV {
	if a.H == nil {
		a.H = H{}
	}
	a.H.Set(k, v)
	return a
}

func (a *KV) SetX(x interface{}) *KV {
	a.X = x
	return a
}

func (a *KV) SetFn(fn func(context.Context) interface{}) *KV {
	a.fn = fn
	return a
}

func (a *KV) Fn() func(context.Context) interface{} {
	return a.fn
}

type KVList []*KV

func (list *KVList) Add(k, v string, options ...KVOption) {
	a := &KV{K: k, V: v}
	for _, option := range options {
		option(a)
	}
	*list = append(*list, a)
}

func (list *KVList) AddItem(item *KV) {
	*list = append(*list, item)
}

func (list *KVList) Delete(i int) {
	n := len(*list)
	if i+1 < n {
		*list = append((*list)[0:i], (*list)[i+1:]...)
	} else if i < n {
		*list = (*list)[0:i]
	}
}

func (list *KVList) Reset() {
	*list = (*list)[0:0]
}

// NewKVData 键值对数据
func NewKVData() *KVData {
	return &KVData{
		slice: []*KV{},
		index: map[string][]int{},
	}
}

// KVData 键值对数据（保持顺序）
type KVData struct {
	slice  []*KV
	index  map[string][]int
	sorted atomic.Bool
	mu     sync.Mutex
}

// Slice 返回切片
func (a *KVData) Slice() []*KV {
	a.Sort()
	return a.slice
}

func (a *KVData) Clone() *KVData {
	b := KVData{
		slice: make([]*KV, len(a.slice)),
		index: map[string][]int{},
		mu:    sync.Mutex{},
	}
	b.sorted.Store(a.sorted.Load())
	for i, v := range a.slice {
		c := v.Clone()
		b.slice[i] = &c
	}
	for name, v := range a.index {
		c := make([]int, len(v))
		copy(c, v)
		b.index[name] = c
	}
	return &b
}

// Keys 返回所有K值
func (a *KVData) Keys() []string {
	a.Sort()
	keys := make([]string, len(a.slice))
	for i, v := range a.slice {
		if v == nil {
			continue
		}
		keys[i] = v.K
	}
	return keys
}

// Index 返回某个key的所有索引值
func (a *KVData) Index(k string) []int {
	v := a.index[k]
	return v
}

// Indexes 返回所有索引值
func (a *KVData) Indexes() map[string][]int {
	return a.index
}

// Reset 重置
func (a *KVData) Reset() *KVData {
	a.index = map[string][]int{}
	a.slice = []*KV{}
	a.sorted.CompareAndSwap(true, false)
	return a
}

// Add 添加键值
func (a *KVData) Add(k, v string, options ...KVOption) *KVData {
	if _, y := a.index[k]; !y {
		a.index[k] = []int{}
	}
	a.index[k] = append(a.index[k], len(a.slice))
	an := &KV{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = append(a.slice, an)
	a.sorted.CompareAndSwap(true, false)
	return a
}

func (a *KVData) AddItem(item *KV) *KVData {
	if _, y := a.index[item.K]; !y {
		a.index[item.K] = []int{}
	}
	a.index[item.K] = append(a.index[item.K], len(a.slice))
	a.slice = append(a.slice, item)
	a.sorted.CompareAndSwap(true, false)
	return a
}

// Set 设置首个键值
func (a *KVData) Set(k, v string, options ...KVOption) *KVData {
	a.index[k] = []int{0}
	an := &KV{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = []*KV{an}
	a.sorted.CompareAndSwap(true, false)
	return a
}

func (a *KVData) SetItem(item *KV) *KVData {
	a.index[item.K] = []int{0}
	a.slice = []*KV{item}
	a.sorted.CompareAndSwap(true, false)
	return a
}

func (a *KVData) Get(k string, defaults ...string) string {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]].V
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

func (a *KVData) GetItem(k string, defaults ...func() *KV) *KV {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]]
		}
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	return nil
}

func (a *KVData) GetByIndex(index int, defaults ...string) string {
	if len(a.slice) > index {
		return a.slice[index].V
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

func (a *KVData) GetItemByIndex(index int, defaults ...func() *KV) *KV {
	if len(a.slice) > index {
		return a.slice[index]
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	return nil
}

func (a *KVData) Size() int {
	return len(a.slice)
}

func (a *KVData) Has(k string) bool {
	_, ok := a.index[k]
	return ok
}

// Delete 设置某个键的所有值
func (a *KVData) Delete(ks ...string) *KVData {
	indexes := []int{}
	for _, k := range ks {
		v, y := a.index[k]
		if !y {
			continue
		}
		indexes = append(indexes, v...)
	}
	newSlice := []*KV{}
	a.index = map[string][]int{}
	for i, v := range a.slice {
		var exists bool
		for _, idx := range indexes {
			if i != idx {
				continue
			}
			exists = true
			break
		}
		if exists {
			continue
		}
		if _, y := a.index[v.K]; !y {
			a.index[v.K] = []int{}
		}
		a.index[v.K] = append(a.index[v.K], len(newSlice))
		newSlice = append(newSlice, v)
	}
	a.slice = newSlice
	a.sorted.CompareAndSwap(true, false)
	return a
}

func (a *KVData) Sort() *KVData {
	if a.sorted.CompareAndSwap(false, true) {
		a.mu.Lock()
		sort.Sort(a)
		a.mu.Unlock()
	}
	return a
}

// sort.Interface

func (a *KVData) Len() int {
	return len(a.slice)
}

func (a *KVData) Less(i, j int) bool {
	return a.slice[i].priority > a.slice[j].priority
}

func (a *KVData) Swap(i, j int) {
	var n int
	if a.slice[i].K == a.slice[j].K {
		for index, sindex := range a.index[a.slice[i].K] {
			switch sindex {
			case i:
				a.index[a.slice[i].K][index] = j
				n++
			case j:
				a.index[a.slice[i].K][index] = i
				n++
			default:
				if n >= 2 {
					goto END
				}
			}
		}
	} else {
		for key, values := range a.index {
			for index, sindex := range values {
				switch sindex {
				case i:
					a.index[key][index] = j
					n++
				case j:
					a.index[key][index] = i
					n++
				default:
					if n >= 2 {
						goto END
					}
				}
			}
		}
	}

END:
	a.slice[i], a.slice[j] = a.slice[j], a.slice[i]
}
