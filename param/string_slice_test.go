package param

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	v := StringSlice([]string{"a", "a", "b", "c", "c", "e", "e"})
	r := v.Unique().Join(`,`)
	if r != `a,b,c,e` {
		panic(fmt.Sprintf("%v != a,b,c,e", r))
	}
	v = StringSlice([]string{"a", "a"})
	r = v.Unique().Join(`,`)
	if r != `a` {
		panic(fmt.Sprintf("%v != a", r))
	}
}
