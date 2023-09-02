//go:build go1.18

package param

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
