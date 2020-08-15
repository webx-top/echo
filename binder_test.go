package echo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestForm struct {
	Files  []string
	Data   *TestData
	IDs    string `form_delimiter:","`
	secret string
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
	*TestProfile
	Name string
	Age  uint
}

type TestProfile struct {
	Address string
}

type TestAnonymous struct {
	*TestUser
	Title      string
	ListStruct []*TestAnonymous
	ListString []string
	IsOk       *sql.NullBool
	Alias      *string
	Time       time.Time
}

func TestMapToAnonymous(t *testing.T) {
	e := New()
	m := &TestAnonymous{}
	formData := map[string][]string{
		`name`:                   {`lily`},
		`age`:                    {`1`},
		`title`:                  {`test`},
		`listStruct[0][address]`: {`secret`},
		`listStruct[0][name]`:    {`a`},
		`listStruct[0][age]`:     {`2`},
		`listStruct[0][title]`:   {`test2`},
		`listStruct[1][address]`: {`secret3`},
		`listStruct[1][name]`:    {`b`},
		`listStruct[1][age]`:     {`3`},
		`listStruct[1][title]`:   {`test3`},
		`listString[]`:           {`A`, `B`},
		`isOk`:                   {`1`},
		`alias`:                  {`hah`},
		`time`:                   {`2020-08-10 12:00:00`},
	}
	NamedStructMap(e, m, formData, ``)
	s := `hah`
	tm, _ := time.ParseInLocation(`2006-01-02 15:04:05`, `2020-08-10 12:00:00`, time.Local)
	expected := &TestAnonymous{
		TestUser: &TestUser{
			TestProfile: &TestProfile{Address: ``},
			Name:        `lily`,
			Age:         1,
		},
		Title: `test`,
		ListStruct: []*TestAnonymous{
			{
				TestUser: &TestUser{
					TestProfile: &TestProfile{
						Address: `secret`,
					},
					Name: `a`,
					Age:  2,
				},
				Title: `test2`,
			},
			{
				TestUser: &TestUser{
					TestProfile: &TestProfile{
						Address: `secret3`,
					},
					Name: `b`,
					Age:  3,
				},
				Title: `test3`,
			},
		},
		ListString: []string{`A`, `B`},
		IsOk: &sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Alias: &s,
		Time:  tm,
	}
	assert.Equal(t, expected, m)

	for _, v := range expected.ListStruct {
		v.TestUser.Name = ``
	}
	for _, v := range m.ListStruct {
		v.TestUser.Name = ``
	}
	NamedStructMap(e, m, formData, ``, ExcludeFieldName(`*.*.Name`))
	assert.Equal(t, expected, m)
}

func TestMapToStruct(t *testing.T) {
	e := New()
	m := &TestForm{}
	NamedStructMap(e, m, map[string][]string{
		`files[]`:             {`a.txt`, `b.txt`, `c.txt`},
		`data[file][users][]`: {`1`, `2`, `3`},
		`IDs[]`:               {`1`, `2`, `3`},
		`secret`:              {`nothing`},
	}, ``)
	assert.Equal(t, &TestForm{
		Files: []string{`a.txt`, `b.txt`, `c.txt`},
		Data: &TestData{
			File: &TestFile{
				Users: []int64{1, 2, 3},
			},
		},
		IDs:    `1,2,3`,
		secret: ``,
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
