package echo

import "context"

type KVxOption[X any, Y any] func(*KVx[X, Y])

func KVxOptK[X any, Y any](k string) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.K = k
	}
}

func KVxOptV[X any, Y any](v string) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.V = v
	}
}

func KVxOptPriority[X any, Y any](priority int) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.priority = priority
	}
}

func KVxOptH[X any, Y any](h H) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.H = h
	}
}

func KVxOptHKV[X any, Y any](k string, v interface{}) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		if a.H == nil {
			a.H = H{}
		}
		a.H.Set(k, v)
	}
}

func KVxOptX[X any, Y any](x X) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.X = x
	}
}

func KVxOptFn[X any, Y any](fn func(context.Context) Y) KVxOption[X, Y] {
	return func(a *KVx[X, Y]) {
		a.fn = fn
	}
}
