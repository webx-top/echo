package echo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var h = HandlerFunc(func(c Context) error {
	return c.String(`OK`)
})

func TestRouterRegexpKind(t *testing.T) {
	e := New()
	r := NewRouter(e)
	rt := &Route{
		Method:  `GET`,
		Path:    `/1_<id:[\d]+>_2/123`,
		Handler: h,
	}
	r.Add(rt, 0)
	assert.Equal(t, `/1_%v_2/123`, rt.Format)
	assert.Equal(t, []string{`id`}, rt.Params)
	var h2 = HandlerFunc(func(c Context) error {
		return c.String(`OK2`)
	})
	rt = &Route{
		Method:  `GET`,
		Path:    `/1_<id:[\d]+>_2/123/<id2:[\d]+>`,
		Handler: h2,
	}
	r.Add(rt, 0)
	assert.Equal(t, `/1_%v_2/123/%v`, rt.Format)
	rt = &Route{
		Method:  `GET`,
		Path:    `/g/<id:[\d]+>`,
		Handler: h2,
	}
	r.Add(rt, 0)
	assert.Equal(t, `/g/%v`, rt.Format)
	rt = &Route{
		Method:  `GET`,
		Path:    `/g/b<id:[\d]+>`,
		Handler: h2,
	}
	r.Add(rt, 0)
	assert.Equal(t, `/g/b%v`, rt.Format)
	//fmt.Println(r.tree.String())
	ctx := e.NewContext(nil, nil)
	found := r.Find(`GET`, `/1_2000_2/123`, ctx)
	assert.True(t, found)
	assert.Equal(t, fmt.Sprintf(`%p`, h), fmt.Sprintf(`%p`, ctx.(*xContext).handler))

	ctx = e.NewContext(nil, nil)
	found = r.Find(`GET`, `/1_2000_2/123/100`, ctx)
	assert.True(t, found)
	assert.Equal(t, fmt.Sprintf(`%p`, h2), fmt.Sprintf(`%p`, ctx.(*xContext).handler))

	ctx = e.NewContext(nil, nil)
	found = r.Find(`GET`, `/g/100`, ctx)
	assert.True(t, found)
	assert.Equal(t, fmt.Sprintf(`%p`, h2), fmt.Sprintf(`%p`, ctx.(*xContext).handler))

	ctx = e.NewContext(nil, nil)
	found = r.Find(`GET`, `/g/a`, ctx)
	assert.False(t, found)

	ctx = e.NewContext(nil, nil)
	found = r.Find(`GET`, `/g/b100`, ctx)
	assert.True(t, found)
	assert.Equal(t, fmt.Sprintf(`%p`, h2), fmt.Sprintf(`%p`, ctx.(*xContext).handler))
}

func TestRouterParamKind(t *testing.T) {
	e := New()
	r := NewRouter(e)
	r.Add(&Route{
		Method:  `GET`,
		Path:    `/:id`,
		Handler: h,
	}, 0)
	fmt.Println(r.tree.String())
	assert.Equal(t, `/`, string(r.tree.label))
	assert.Equal(t, skind, r.tree.kind)
	assert.Equal(t, ``, r.tree.ppath)
	assert.Equal(t, `/`, r.tree.prefix)

	assert.Equal(t, `:`, string(r.tree.children[0].label))
	assert.Equal(t, pkind, r.tree.children[0].kind)
	assert.Equal(t, `/:id`, r.tree.children[0].ppath)
	assert.Equal(t, `:`, r.tree.children[0].prefix)
}
