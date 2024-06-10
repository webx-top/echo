package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParamToRegexpRule(t *testing.T) {
	r, rv, ps, kv := QueryParamToRegexpRule(`/a/b/:id`, false)
	assert.Equal(t, `^/a/b/(?P<id>[^/]+)$`, r)
	assert.Equal(t, `/a/b/$1`, rv)
	assert.Equal(t, []string{`(?P<id>[^/]+)`}, ps)
	assert.Equal(t, map[string]string{`id`: `(?P<id>[^/]+)`}, kv)
	r, rv, ps, kv = QueryParamToRegexpRule(`/a/b/:id/:city`, false)
	assert.Equal(t, `^/a/b/(?P<id>[^/]+)/(?P<city>[^/]+)$`, r)
	assert.Equal(t, `/a/b/$1/$2`, rv)
	assert.Equal(t, []string{`(?P<id>[^/]+)`, `(?P<city>[^/]+)`}, ps)
	assert.Equal(t, map[string]string{`id`: `(?P<id>[^/]+)`, `city`: `(?P<city>[^/]+)`}, kv)
	r, rv, ps, kv = QueryParamToRegexpRule(`/a/b/<name:[a-z]+>`, false)
	assert.Equal(t, `^/a/b/(?P<name>[a-z]+)$`, r)
	assert.Equal(t, `/a/b/$1`, rv)
	assert.Equal(t, []string{`(?P<name>[a-z]+)`}, ps)
	assert.Equal(t, map[string]string{`name`: `(?P<name>[a-z]+)`}, kv)
	r, rv, ps, kv = QueryParamToRegexpRule(`/a/b/<name:[a-z]+>/<city:[a-z]+>`, false)
	assert.Equal(t, `^/a/b/(?P<name>[a-z]+)/(?P<city>[a-z]+)$`, r)
	assert.Equal(t, `/a/b/$1/$2`, rv)
	assert.Equal(t, []string{`(?P<name>[a-z]+)`, `(?P<city>[a-z]+)`}, ps)
	assert.Equal(t, map[string]string{`name`: `(?P<name>[a-z]+)`, `city`: `(?P<city>[a-z]+)`}, kv)
	r, rv, ps, kv = QueryParamToRegexpRule(`http://a.b.c/b/<name:[a-z]+>/<city:[a-z]+>`, true)
	assert.Equal(t, `^http://a\.b\.c/b/(?P<name>[a-z]+)/(?P<city>[a-z]+)$`, r)
	assert.Equal(t, `http://a.b.c/b/$1/$2`, rv)
	assert.Equal(t, []string{`(?P<name>[a-z]+)`, `(?P<city>[a-z]+)`}, ps)
	assert.Equal(t, map[string]string{`name`: `(?P<name>[a-z]+)`, `city`: `(?P<city>[a-z]+)`}, kv)
}

func TestRewriteConfig(t *testing.T) {
	c := &RewriteConfig{
		Rules: map[string]string{
			`/a/b/:id`:                      `/n_$1`,
			`/c/d/<name:[a-z]+>`:            `/m_$1`,
			`/c/d/<name:[a-z]+>/add`:        `/m_$1_add`,
			`/c/d/<name:[a-z]+>/edit/:id`:   `/m_$1_edit_$2`,
			`/c/d/<name1:[a-z]+>/named/:id`: `/m_{name1}_named_{id}`,
		},
	}
	c.Init()
	r := c.Rewrite(`/a/b/hah`)
	assert.Equal(t, `/n_hah`, r)
	assert.Equal(t, `/a/b/hah`, c.Reverse(`/n_hah`))

	r = c.Rewrite(`/c/d/hah`)
	assert.Equal(t, `/m_hah`, r)
	assert.Equal(t, `/c/d/hah`, c.Reverse(`/m_hah`))

	r = c.Rewrite(`/test/c/d/hah`)
	assert.Equal(t, `/test/c/d/hah`, r)
	assert.Equal(t, `/test/c/d/hah`, c.Reverse(`/test/c/d/hah`))

	r = c.Rewrite(`/c/d/hah/add`)
	assert.Equal(t, `/m_hah_add`, r)
	assert.Equal(t, `/c/d/hah/add`, c.Reverse(`/m_hah_add`))

	r = c.Rewrite(`/c/d/hah/edit/1`)
	assert.Equal(t, `/m_hah_edit_1`, r)
	assert.Equal(t, `/c/d/hah/edit/1`, c.Reverse(`/m_hah_edit_1`))

	r = c.Rewrite(`/c/d/hah/named/1`)
	assert.Equal(t, `/m_hah_named_1`, r)
	assert.Equal(t, `/c/d/hah/named/1`, c.Reverse(`/m_hah_named_1`))
}

func TestRewriteConfig2(t *testing.T) {
	c := &RewriteConfig{
		Rules: map[string]string{
			`http://a.b.c/`: `http://b.b.c/`,
			`http://a.*.n/`: `http://b.$1.n/`,
		},
		DisableColonParam: true,
	}
	c.Init()
	r := c.Rewrite(`http://a.b.c/`)
	assert.Equal(t, `http://b.b.c/`, r)
	assert.Equal(t, `http://a.b.c/`, c.Reverse(`http://b.b.c/`))

	r = c.Rewrite(`http://a.e.n/`)
	assert.Equal(t, `http://b.e.n/`, r)
	assert.Equal(t, `http://a.e.n/`, c.Reverse(`http://b.e.n/`))
}
