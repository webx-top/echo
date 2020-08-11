package echo

import (
	"reflect"
	"testing"

	"github.com/admpub/copier"
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

type TestRole struct {
	Name  string
	Users []*TestUser
}

type TestRoleM struct {
	Name  string
	Users map[string]*TestUser
}

type TestUser struct {
	Name string
	Age  uint
}

type TestAnonymous struct {
	*TestUser
	Title string
}

func TestMapToAnonymous(t *testing.T) {
	e := New()
	m := &TestAnonymous{}
	v := reflect.ValueOf(m)
	copier.InitNilFields(v.Type(), v, ``, map[string]struct{}{
		`TestUser`: struct{}{},
	})
	NamedStructMap(e, m, map[string][]string{
		`name`:  {`lily`},
		`age`:   {`1`},
		`title`: {`test`},
	}, ``)
	assert.Equal(t, &TestAnonymous{
		TestUser: &TestUser{
			Name: `lily`,
			Age:  1,
		},
		Title: `test`,
	}, m)
}

func TestMapToStruct(t *testing.T) {
	e := New()
	m := &TestForm{}
	NamedStructMap(e, m, map[string][]string{
		`files[]`:             {`a.txt`, `b.txt`, `c.txt`},
		`data[file][users][]`: {`1`, `2`, `3`},
		`IDs[]`:               {`1`, `2`, `3`},
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
		`files`:             {`a.txt`, `b.txt`, `c.txt`},
		`data[file][users]`: {`1`, `2`, `3`},
		`IDs`:               {`1`, `2`, `3`},
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

func TestMapToSliceStruct(t *testing.T) {
	e := New()
	m := &TestRole{}
	NamedStructMap(e, m, map[string][]string{
		`name`:           {`manager`},
		`users[0][name]`: {`john`},
		`users[0][age]`:  {`18`},
		`users[1][name]`: {`smith`},
		`users[1][age]`:  {`25`},
		`users[3][name]`: {`hank`},
		`users[3][age]`:  {`28`},
	}, ``)
	assert.Equal(t, &TestRole{
		Name: `manager`,
		Users: []*TestUser{
			{ // 0
				Name: `john`,
				Age:  18,
			},
			{ // 1
				Name: `smith`,
				Age:  25,
			},
			nil, // 2
			{ // 3
				Name: `hank`,
				Age:  28,
			},
		},
	}, m)
	//Dump(m)
}

func TestMapToMapStruct(t *testing.T) {
	e := New()
	m := &TestRoleM{}
	NamedStructMap(e, m, map[string][]string{
		`name`:           {`manager`},
		`users[0][name]`: {`john`},
		`users[0][age]`:  {`18`},
		`users[1][name]`: {`smith`},
		`users[1][age]`:  {`25`},
		`users[3][name]`: {`hank`},
		`users[3][age]`:  {`28`},
	}, ``)
	assert.Equal(t, &TestRoleM{
		Name: `manager`,
		Users: map[string]*TestUser{
			`0`: { // 0
				Name: `john`,
				Age:  18,
			},
			`1`: { // 1
				Name: `smith`,
				Age:  25,
			},
			`3`: { // 3
				Name: `hank`,
				Age:  28,
			},
		},
	}, m)
	//Dump(m)
}
