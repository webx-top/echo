package echo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestForm struct {
	Files []string
	Data  *TestData
	IDs   string `form_delimiter:","`
}

type TestData struct {
	File *TestFile
}

type TestFile struct {
	Users []int64
}

func TestMapToStruct(t *testing.T) {
	e := New()
	m := &TestForm{}
	NamedStructMap(e, m, map[string][]string{
		`files[]`:             []string{`a.txt`, `b.txt`, `c.txt`},
		`data[file][users][]`: []string{`1`, `2`, `3`},
		`IDs[]`:               []string{`1`, `2`, `3`},
	}, ``)
	assert.Equal(t, &TestForm{
		Files: []string{`a.txt`, `b.txt`, `c.txt`},
		Data: &TestData{
			File: &TestFile{
				Users: []int64{1, 2, 3},
			},
		},
		IDs: `1,2,3`,
	}, m)
}

func TestMapToStruct2(t *testing.T) {
	e := New()
	m := &TestForm{}
	NamedStructMap(e, m, map[string][]string{
		`files`:             []string{`a.txt`, `b.txt`, `c.txt`},
		`data[file][users]`: []string{`1`, `2`, `3`},
		`IDs`:               []string{`1`, `2`, `3`},
	}, ``)
	assert.Equal(t, &TestForm{
		Files: []string{`a.txt`, `b.txt`, `c.txt`},
		Data: &TestData{
			File: &TestFile{
				Users: []int64{1, 2, 3},
			},
		},
		IDs: `1,2,3`,
	}, m)
}
