package param

import (
	"strconv"
)

type Param string

func (p Param) String() string {
	return string(p)
}

func (p Param) Int() int {
	if len(p) > 0 {
		r, _ := strconv.Atoi(p.String())
		return r
	}
	return 0
}

func (p Param) Int64() int64 {
	if len(p) > 0 {
		r, _ := strconv.ParseInt(p.String(), 10, 64)
		return r
	}
	return 0
}

func (p Param) Int32() int32 {
	if len(p) > 0 {
		r, _ := strconv.ParseInt(p.String(), 10, 32)
		return int32(r)
	}
	return 0
}

func (p Param) Uint() uint {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 64)
		return uint(r)
	}
	return 0
}

func (p Param) Uint64() uint64 {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 64)
		return r
	}
	return 0
}

func (p Param) Uint32() uint32 {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 32)
		return uint32(r)
	}
	return 0
}

func (p Param) Float32() float32 {
	if len(p) > 0 {
		r, _ := strconv.ParseFloat(p.String(), 32)
		return float32(r)
	}
	return 0
}

func (p Param) Float64() float64 {
	if len(p) > 0 {
		r, _ := strconv.ParseFloat(p.String(), 64)
		return r
	}
	return 0
}

func (p Param) Bool() bool {
	if len(p) > 0 {
		r, _ := strconv.ParseBool(p.String())
		return r
	}
	return false
}
