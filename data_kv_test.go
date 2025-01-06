package echo_test

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestDataKVSort(t *testing.T) {
	m := echo.NewKVData()
	third := echo.NewKV(`third`, `3`).SetPriority(0)
	second := echo.NewKV(`second`, `2`).SetPriority(10)
	first := echo.NewKV(`first`, `1`).SetPriority(30)
	m.AddItem(third)
	m.AddItem(second)
	m.AddItem(first)
	test.Eq(t, third, m.GetItem(`third`))
	test.Eq(t, second, m.GetItem(`second`))
	test.Eq(t, first, m.GetItem(`first`))
	test.Eq(t, []*echo.KV{third, second, first}, m.Slice())
	m.Sort()
	test.Eq(t, []*echo.KV{first, second, third}, m.Slice())
	test.Eq(t, third, m.GetItem(`third`))
	test.Eq(t, second, m.GetItem(`second`))
	test.Eq(t, first, m.GetItem(`first`))
}

func TestDataKVxSort(t *testing.T) {
	m := echo.NewKVxData[any, any]()
	third := echo.NewKVx[any, any](`third`, `3`).SetPriority(0)
	second := echo.NewKVx[any, any](`second`, `2`).SetPriority(10)
	first := echo.NewKVx[any, any](`first`, `1`).SetPriority(30)
	m.AddItem(third)
	m.AddItem(second)
	m.AddItem(first)
	test.Eq(t, third, m.GetItem(`third`))
	test.Eq(t, second, m.GetItem(`second`))
	test.Eq(t, first, m.GetItem(`first`))
	test.Eq(t, []*echo.KVx[any, any]{third, second, first}, m.Slice())
	m.Sort()
	test.Eq(t, []*echo.KVx[any, any]{first, second, third}, m.Slice())
	test.Eq(t, third, m.GetItem(`third`))
	test.Eq(t, second, m.GetItem(`second`))
	test.Eq(t, first, m.GetItem(`first`))
}
