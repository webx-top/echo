package param

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	var v StringMap
	fmt.Println(`a:`, v.String(`a`))
	fmt.Println(`b:`, v.Int(`b`))
	fmt.Println(`c:`, v.Float64(`c`))
	fmt.Println(`d:`, v.Bool(`d`))
	d := Store{
		`number`:   1.234567e+06,
		`float64`:  float64(1.234),
		`float32`:  float32(1.234),
		`-float64`: -float64(1.234),
		`-float32`: -float32(1.234),
	}
	if d.Int(`float64`) != 1 {
		panic(fmt.Sprintf("%v != 1", d.Int(`float64`)))
	}
	if d.Int(`float32`) != 1 {
		panic(fmt.Sprintf("%v != 1", d.Int(`float32`)))
	}
	if d.Int(`-float64`) != -1 {
		panic(fmt.Sprintf("%v != -1", d.Int(`-float64`)))
	}
	if d.Int(`-float32`) != -1 {
		panic(fmt.Sprintf("%v != -1", d.Int(`-float32`)))
	}
	if d.Uint(`-float64`) != 0 {
		panic(fmt.Sprintf("%v != 0", d.Uint(`-float64`)))
	}
	if d.Int(`number`) != 1234567 {
		panic(fmt.Sprintf("%v != 1234567", d.Int(`number`)))
	}
}
