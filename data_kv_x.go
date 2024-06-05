//go:build go1.18

package echo

import "context"

func NewKVx[X any, Y any](k, v string) *KVx[X, Y] {
	return &KVx[X, Y]{K: k, V: v}
}

// KV 键值对
type KVx[X any, Y any] struct {
	K  string
	V  string
	H  H `json:",omitempty" xml:",omitempty"`
	X  X `json:",omitempty" xml:",omitempty"`
	fn func(context.Context) Y
}

func (a *KVx[X, Y]) SetK(k string) *KVx[X, Y] {
	a.K = k
	return a
}

func (a *KVx[X, Y]) SetV(v string) *KVx[X, Y] {
	a.V = v
	return a
}

func (a *KVx[X, Y]) SetKV(k, v string) *KVx[X, Y] {
	a.K = k
	a.V = v
	return a
}

func (a *KVx[X, Y]) SetH(h H) *KVx[X, Y] {
	a.H = h
	return a
}

func (a *KVx[X, Y]) SetHKV(k string, v interface{}) *KVx[X, Y] {
	if a.H == nil {
		a.H = H{}
	}
	a.H.Set(k, v)
	return a
}

func (a *KVx[X, Y]) SetX(x X) *KVx[X, Y] {
	a.X = x
	return a
}

func (a *KVx[X, Y]) SetFn(fn func(context.Context) Y) *KVx[X, Y] {
	a.fn = fn
	return a
}

func (a *KVx[X, Y]) Fn() func(context.Context) Y {
	return a.fn
}

type KVxList[X any, Y any] []*KVx[X, Y]

func (list *KVxList[X, Y]) Add(k, v string, options ...KVxOption[X, Y]) {
	a := &KVx[X, Y]{K: k, V: v}
	for _, option := range options {
		option(a)
	}
	*list = append(*list, a)
}

func (list *KVxList[X, Y]) AddItem(item *KVx[X, Y]) {
	*list = append(*list, item)
}

func (list *KVxList[X, Y]) Delete(i int) {
	n := len(*list)
	if i+1 < n {
		*list = append((*list)[0:i], (*list)[i+1:]...)
	} else if i < n {
		*list = (*list)[0:i]
	}
}

func (list *KVxList[X, Y]) Reset() {
	*list = (*list)[0:0]
}

// NewKVxData 键值对数据
func NewKVxData[X any, Y any]() *KVxData[X, Y] {
	return &KVxData[X, Y]{
		slice: []*KVx[X, Y]{},
		index: map[string][]int{},
	}
}

// KVxData 键值对数据（保持顺序）
type KVxData[X any, Y any] struct {
	slice []*KVx[X, Y]
	index map[string][]int
}

// Slice 返回切片
func (a *KVxData[X, Y]) Slice() []*KVx[X, Y] {
	return a.slice
}

// Keys 返回所有K值
func (a *KVxData[X, Y]) Keys() []string {
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
func (a *KVxData[X, Y]) Index(k string) []int {
	v := a.index[k]
	return v
}

// Indexes 返回所有索引值
func (a *KVxData[X, Y]) Indexes() map[string][]int {
	return a.index
}

// Reset 重置
func (a *KVxData[X, Y]) Reset() *KVxData[X, Y] {
	a.index = map[string][]int{}
	a.slice = []*KVx[X, Y]{}
	return a
}

// Add 添加键值
func (a *KVxData[X, Y]) Add(k, v string, options ...KVxOption[X, Y]) *KVxData[X, Y] {
	if _, y := a.index[k]; !y {
		a.index[k] = []int{}
	}
	a.index[k] = append(a.index[k], len(a.slice))
	an := &KVx[X, Y]{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = append(a.slice, an)
	return a
}

func (a *KVxData[X, Y]) AddItem(item *KVx[X, Y]) *KVxData[X, Y] {
	if _, y := a.index[item.K]; !y {
		a.index[item.K] = []int{}
	}
	a.index[item.K] = append(a.index[item.K], len(a.slice))
	a.slice = append(a.slice, item)
	return a
}

// Set 设置首个键值
func (a *KVxData[X, Y]) Set(k, v string, options ...KVxOption[X, Y]) *KVxData[X, Y] {
	a.index[k] = []int{0}
	an := &KVx[X, Y]{K: k, V: v}
	for _, option := range options {
		option(an)
	}
	a.slice = []*KVx[X, Y]{an}
	return a
}

func (a *KVxData[X, Y]) SetItem(item *KVx[X, Y]) *KVxData[X, Y] {
	a.index[item.K] = []int{0}
	a.slice = []*KVx[X, Y]{item}
	return a
}

func (a *KVxData[X, Y]) Get(k string, defaults ...string) string {
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

func (a *KVxData[X, Y]) GetItem(k string, defaults ...func() *KVx[X, Y]) *KVx[X, Y] {
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

func (a *KVxData[X, Y]) GetByIndex(index int, defaults ...string) string {
	if len(a.slice) > index {
		return a.slice[index].V
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

func (a *KVxData[X, Y]) GetItemByIndex(index int, defaults ...func() *KVx[X, Y]) *KVx[X, Y] {
	if len(a.slice) > index {
		return a.slice[index]
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	return nil
}

func (a *KVxData[X, Y]) Size() int {
	return len(a.slice)
}

func (a *KVxData[X, Y]) Has(k string) bool {
	_, ok := a.index[k]
	return ok
}

// Delete 设置某个键的所有值
func (a *KVxData[X, Y]) Delete(ks ...string) *KVxData[X, Y] {
	indexes := []int{}
	for _, k := range ks {
		v, y := a.index[k]
		if !y {
			continue
		}
		indexes = append(indexes, v...)
	}
	newSlice := []*KVx[X, Y]{}
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
	return a
}
