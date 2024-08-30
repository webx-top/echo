package param

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestConversion(t *testing.T) {
	var f64 float64 = math.MaxFloat64
	f32 := AsFloat32(f64)
	assert.Equal(t, float32(0), f32)
	var ui64 uint64 = math.MaxUint64
	f32 = AsFloat32(ui64)
	assert.Equal(t, com.NumberFormat(ui64, 0), com.NumberFormat(f32, 0))

	ui32a := AsUint32(ui64)
	assert.Equal(t, uint32(math.MaxUint32), ui32a)
	assert.Greater(t, ui64, uint64(ui32a))

	ui16a := AsUint16(ui64)
	assert.Equal(t, uint16(math.MaxUint16), ui16a)
	assert.Greater(t, ui64, uint64(ui16a))

	var i int64 = math.MinInt
	ui64a := AsUint64(i)
	assert.Equal(t, uint64(0), ui64a)

	ui32a = AsUint32(ui64a)
	assert.Equal(t, uint32(0), ui32a)
	ui32a = AsUint32(i)
	assert.Equal(t, uint32(0), ui32a)

	ui16a = AsUint16(ui32a)
	assert.Equal(t, uint16(0), ui16a)
	ui16a = AsUint16(i)
	assert.Equal(t, uint16(0), ui16a)

	ui64a = AsUint64(f64)
	com.Dump(Store{
		`ui64a`:      com.NumberFormat(ui64a, 3),
		`MaxUint64`:  com.NumberFormat(ui64, 3),
		`MaxFloat64`: com.NumberFormat(f64, 3),
	})
	assert.Equal(t, uint64(0), ui64a)

	ui32a = AsUint32(i)
	assert.Equal(t, uint32(0), ui32a)
}

func TestConverts(t *testing.T) {
	src := []string{`1`, `2`}
	v := Converts[uint](src)
	excepted := []uint{1, 2}
	assert.Equal(t, excepted, v)
	v2 := Converts[string](excepted)
	excepted2 := src
	assert.Equal(t, excepted2, v2)

	type myInt int

	v3 := Converts[myInt](excepted, func(s uint) myInt {
		return myInt(s)
	})
	excepted3 := []myInt{myInt(1), myInt(2)}
	assert.Equal(t, excepted3, v3)

	v4 := Converts[myInt](src, func(s string) myInt {
		return myInt(AsInt(s))
	})
	excepted4 := excepted3
	assert.Equal(t, excepted4, v4)
}
