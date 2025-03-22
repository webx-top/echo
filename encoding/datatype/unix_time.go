package datatype

import (
	"strconv"
	"time"
)

type UnixTime time.Time

func (u UnixTime) MarshalText() ([]byte, error) {
	seconds := time.Time(u).Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

type UnixMilli time.Time

func (u UnixMilli) MarshalText() ([]byte, error) {
	seconds := time.Time(u).UnixMilli()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}
