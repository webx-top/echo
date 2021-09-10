package echo_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestHostParse(t *testing.T) {
	h := NewHost(`<a:[0-9]+>.<b>.cc`)
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

func TestHostGroupHandler(t *testing.T) {
	e := New()
	e.OnHostFound(func(c Context) (bool, error) {
		if c.HostParam(`name`) != `coscms` {
			return true, nil
		}
		fmt.Println(`host param (uid):`, c.HostParam(`uid`))
		switch c.HostParam(`uid`) {
		case `123`:
			return true, nil
		case `0`:
			return false, errors.New(`err`)
		default:
			return false, nil
		}
	})

	// suffix
	g := e.Host(".admpub.com")
	g.Get("/host", func(c Context) error {
		if c.Queryx(`route`).Bool() {
			return c.JSON(c.Route())
		}
		return c.String(c.Host())
	}).SetName(`host`)

	// rule
	g = e.Host("<uid:[0-9]+>.<name>.com").SetAlias(`user`)
	g.Get("/host2", func(c Context) error {
		if c.Queryx(`route`).Bool() {
			return c.JSON(c.Route())
		}
		return c.String(c.HostParam(`uid`) + `.` + c.HostParam(`name`) + `.com`)
	}).SetName(`host2`)

	e.RebuildRouter()

	c, b := request(GET, "/host", e, func(req *http.Request) {
		req.Host = "test.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "test.admpub.com", b)

	c, b = request(GET, "/host", e, func(req *http.Request) {
		req.Host = "test-b.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "test-b.admpub.com", b)

	c, b = request(GET, "/host?route=1", e, func(req *http.Request) {
		req.Host = "test-b.admpub.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `{"Host":".admpub.com","Method":"GET","Path":"/host","Name":"host","Format":"/host","Params":[],"Prefix":"","Meta":{}}`, b)

	c, b = request(GET, "/host2", e, func(req *http.Request) {
		req.Host = "123.coscms.com"
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "123.coscms.com", b)
	assert.Equal(t, "10000.admpub.com/host2", e.TypeHost(`user`, 10000, `admpub`).URI(`host2`))
	assert.Equal(t, "10001.admpub.com/host2", e.TypeHost(`user`, H{`uid`: 10001, `name`: `admpub`}).URI(`host2`))

	c, b = request(GET, "/host2", e, func(req *http.Request) {
		req.Host = "0.coscms.com"
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), b)

	c, b = request(GET, "/host2", e, func(req *http.Request) {
		req.Host = "3.coscms.com"
	})
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, http.StatusText(http.StatusNotFound), b)
}
