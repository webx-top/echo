package echo_test

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestHostParse(t *testing.T) {
	h := echo.NewHost(`<a:[0-9]+>.<b>.cc`)
	h.Parse()
	test.Eq(t, `%v.%v.cc`, h.Format())
	test.Eq(t, `^([0-9]+)\.([^.]+)\.cc$`, h.RegExp().String())
	r := h.RegExp().FindStringSubmatch(`1.baa.cc`)
	test.Eq(t, []string{
		`1.baa.cc`,
		`1`,
		`baa`,
	}, r)
}
