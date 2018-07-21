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
}
