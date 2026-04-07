package echo

import "github.com/webx-top/echo/param"

// Get retrieves data from the context.
func (c *XContext) Get(key string, defaults ...any) any {
	return c.store.Get(key, defaults...)
}

// Set saves data in the context.
func (c *XContext) Set(key string, val any) {
	c.store.Set(key, val)
}

func (c *XContext) SetMore(keyVal ...any) {
	c.store.SetMore(keyVal...)
}

// Incr Increment the value and return the new value
func (c *XContext) Incr(key string, n any, defaults ...any) int64 {
	return c.store.Incr(key, param.AsInt64(n), defaults...)
}

// Decr Decrement the value and return the new value
func (c *XContext) Decr(key string, n any, defaults ...any) int64 {
	return c.store.Decr(key, param.AsInt64(n), defaults...)
}

// Delete saves data in the context.
func (c *XContext) Delete(keys ...string) {
	c.store.Delete(keys...)
}

func (c *XContext) Stored() Store {
	return c.store.CloneStore()
}
