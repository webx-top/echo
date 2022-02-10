package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParamToRegexpRule(t *testing.T) {
	r, rv, ps := QueryParamToRegexpRule(`/a/b/:id`)
	assert.Equal(t, `^/a/b/([^/]+)$`, r)
	assert.Equal(t, `/a/b/$1`, rv)
	assert.Equal(t, []string{`([^/]+)`}, ps)
	r, rv, ps = QueryParamToRegexpRule(`/a/b/:id/:city`)
	assert.Equal(t, `^/a/b/([^/]+)/([^/]+)$`, r)
	assert.Equal(t, `/a/b/$1/$2`, rv)
	assert.Equal(t, []string{`([^/]+)`, `([^/]+)`}, ps)
	r, rv, ps = QueryParamToRegexpRule(`/a/b/<name:[a-z]+>`)
	assert.Equal(t, `^/a/b/([a-z]+)$`, r)
	assert.Equal(t, `/a/b/$1`, rv)
	assert.Equal(t, []string{`([a-z]+)`}, ps)
	r, rv, ps = QueryParamToRegexpRule(`/a/b/<name:[a-z]+>/<city:[a-z]+>`)
	assert.Equal(t, `^/a/b/([a-z]+)/([a-z]+)$`, r)
	assert.Equal(t, `/a/b/$1/$2`, rv)
	assert.Equal(t, []string{`([a-z]+)`, `([a-z]+)`}, ps)
}

func TestRewriteConfig(t *testing.T) {
	c := &RewriteConfig{
		Rules: map[string]string{
			`/a/b/:id`:                    `/n_$1`,
			`/c/d/<name:[a-z]+>`:          `/m_$1`,
			`/c/d/<name:[a-z]+>/add`:      `/m_$1_add`,
			`/c/d/<name:[a-z]+>/edit/:id`: `/m_$1_edit_$2`,
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
}
