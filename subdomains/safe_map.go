package subdomains

import "sync"

func InitSafeMap[T any]() SafeMap[T] {
	return SafeMap[T]{mp: map[string]T{}}
}

func NewSafeMap[T any]() *SafeMap[T] {
	return &SafeMap[T]{mp: map[string]T{}}
}

type SafeMap[T any] struct {
	mu sync.RWMutex
	mp map[string]T
}

func (a *SafeMap[T]) Size() int {
	a.mu.RLock()
	size := len(a.mp)
	a.mu.RUnlock()
	return size
}

func (a *SafeMap[T]) GetOk(key string) (T, bool) {
	a.mu.RLock()
	val, ok := a.mp[key]
	a.mu.RUnlock()
	return val, ok
}

func (a *SafeMap[T]) Get(key string) T {
	a.mu.RLock()
	val := a.mp[key]
	a.mu.RUnlock()
	return val
}

func (a *SafeMap[T]) Gets(keys ...string) []T {
	a.mu.RLock()
	res := make([]T, 0, len(keys))
	for _, key := range keys {
		val, ok := a.mp[key]
		if ok {
			res = append(res, val)
		}
	}
	a.mu.RUnlock()
	return res
}

func (a *SafeMap[T]) GetsFunc(keys ...func() string) []T {
	a.mu.RLock()
	res := make([]T, 0, len(keys))
	for _, key := range keys {
		val, ok := a.mp[key()]
		if ok {
			res = append(res, val)
		}
	}
	a.mu.RUnlock()
	return res
}

func (a *SafeMap[T]) Set(key string, info T) {
	a.mu.Lock()
	a.mp[key] = info
	a.mu.Unlock()
}

func (a *SafeMap[T]) Remove(keys ...string) {
	a.mu.Lock()
	for _, key := range keys {
		delete(a.mp, key)
	}
	a.mu.Unlock()
}

func (a *SafeMap[T]) Range(f func(key string, val T) bool) {
	a.mu.RLock()
	for key, val := range a.mp {
		if !f(key, val) {
			break
		}
	}
	a.mu.RUnlock()
}

func (a *SafeMap[T]) ClearEmpty(f func(key string, val T) bool) {
	a.mu.Lock()
	for key, val := range a.mp {
		if f(key, val) {
			delete(a.mp, key)
		}
	}
	a.mu.Unlock()
}
