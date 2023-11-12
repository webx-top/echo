package echo_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/param"
)

type TestForm struct {
	Files    []string
	Data     *TestData
	IDs      string `form_delimiter:","`
	Interval time.Duration
	secret   string
}

type TestData struct {
	File     *TestFile
	Extra    H
	List     []interface{}
	Strings  param.StringSlice
	SStrings []param.String
	String   param.String
}

type TestFile struct {
	Users []int64
}

type TestRole struct {
	Name  string
	Users []*TestUser
}

type TestRoleM struct {
	Name     string
	Users    map[string]*TestUser
	Profiled map[string]map[string]*TestProfile
	Profiles []map[string]*TestProfile
	Profilea []map[string]string
	Slices   map[string][]string
	Data     H
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

type TestTypeInt int
type TestTypeString string

type TestMapIntKey struct {
	Map  map[int][]string
	Map1 map[TestTypeInt][]TestTypeString
	Map2 map[TestTypeString]TestTypeInt
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
			TestProfile: nil,
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

	formData[`address`] = []string{``}
	expected.TestUser.TestProfile = &TestProfile{Address: ``}

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
		`data[strings][]`:     {`a`, `b`, `c`},
		`data[sStrings][]`:    {`a`, `b`, `c`},
		`data[string]`:        {`a`},
		`secret`:              {`nothing`},
	}, ``)
	assert.Equal(t, &TestForm{
		Files: []string{`a.txt`, `b.txt`, `c.txt`},
		Data: &TestData{
			File: &TestFile{
				Users: []int64{1, 2, 3},
			},
			Strings:  param.StringSlice{`a`, `b`, `c`},
			SStrings: []param.String{`a`, `b`, `c`},
			String:   `a`,
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
		`data[extra][key]`:  {`v1`},
		`data[extra][key2]`: {`v1`, `v2`},
		`data[list]`:        {`v1`, `v2`},
		`IDs`:               {`1`, `2`, `3`},
		`interval`:          {`5m`},
	}, ``)
	assert.Equal(t, &TestForm{
		Files: []string{`a.txt`, `b.txt`, `c.txt`},
		Data: &TestData{
			File: &TestFile{
				Users: []int64{1, 2, 3},
			},
			Extra: H{
				`key`:  `v1`,
				`key2`: []string{`v1`, `v2`},
			},
			List: []interface{}{`v1`, `v2`},
		},
		IDs:      `1,2,3`,
		Interval: 5 * time.Minute,
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
	err := NamedStructMap(e, m, map[string][]string{
		`name`:                       {`manager`},
		`users[0][name]`:             {`john`},
		`users[0][age]`:              {`18`},
		`users[1][name]`:             {`smith`},
		`users[1][age]`:              {`25`},
		`users[3][name]`:             {`hank`},
		`users[3][age]`:              {`28`},
		`profiled[3][1000][address]`: {`address`},
		`profiles[0][2000][address]`: {`address2`},
		`profilea[0][address]`:       {`address2`},
		`slices[3000][]`:             {`3000v`},
		`slices[3001][]`:             {`3001v`},
	}, ``)
	assert.NoError(t, err)
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
		Profiled: map[string]map[string]*TestProfile{
			`3`: {
				`1000`: {Address: `address`},
			},
		},
		Profiles: []map[string]*TestProfile{
			{
				`2000`: {Address: `address2`},
			},
		},
		Profilea: []map[string]string{
			{`address`: `address2`},
		},
		Slices: map[string][]string{
			`3000`: {`3000v`},
			`3001`: {`3001v`},
		},
	}, m)
	Dump(m)
}

func TestStructToForm(t *testing.T) {
	e := New()
	m := &TestRoleM{
		Name: `test`,
		Users: map[string]*TestUser{
			`user`: &TestUser{
				Name: `user`,
				Age:  10,
			},
		},
	}
	ctx := e.NewContext(mock.NewRequest(), mock.NewResponse())
	StructToForm(ctx, m, ``, DefaultFieldNameFormatter)
	forms := ctx.Forms()

	expected := map[string][]string{
		`Name`:            []string{`test`},
		`Users.user.Age`:  []string{`10`},
		`Users.user.Name`: []string{`user`},
	}
	assert.Equal(t, expected, forms)

	m.Users[`hasProfile`] = &TestUser{
		Name: `userWithProfile`,
		Age:  20,
		TestProfile: &TestProfile{
			Address: `China`,
		},
	}
	m.Data = H{
		`data`:   nil,
		`number`: 1,
	}
	StructToForm(ctx, m, ``, DefaultFieldNameFormatter)
	forms = ctx.Forms()
	expected[`Users.hasProfile.Age`] = []string{`20`}
	expected[`Users.hasProfile.Name`] = []string{`userWithProfile`}
	expected[`Users.hasProfile.TestProfile.Address`] = []string{`China`}
	expected[`Data.number`] = []string{`1`}
	assert.Equal(t, expected, forms)
	//Dump(forms)

	m2 := &TestRoleM{
		Data: H{`data`: nil, `number`: 0},
	}
	NamedStructMap(e, m2, expected, ``)
	assert.Equal(t, m, m2)
	//Dump(m2)
}

func TestStructMapIntKey(t *testing.T) {
	e := New()
	m := &TestMapIntKey{}
	err := NamedStructMap(e, m, map[string][]string{
		`map[1]`:    {`manager`},
		`map1[1][]`: {`manager1`},
		`map2[a]`:   {`2`},
	}, ``)
	assert.NoError(t, err)
	assert.Equal(t, &TestMapIntKey{
		Map: map[int][]string{
			1: {`manager`},
		},
		Map1: map[TestTypeInt][]TestTypeString{
			1: {`manager1`},
		},
		Map2: map[TestTypeString]TestTypeInt{
			`a`: 2,
		},
	}, m)
}

type TestBinderWithConvertor struct {
	Options map[string]string `form_decoder:"splitKVRows:="`
	Env     []string          `form_decoder:"split:\n" form_encoder:"join:\n"`
}

func TestBinderConvertor(t *testing.T) {
	e := New()
	m := &TestBinderWithConvertor{}
	err := NamedStructMap(e, m, map[string][]string{
		`options`: {"a=1\nb=2"},
		`env`:     {"A=ONE\nB=TWO"},
	}, ``)
	assert.NoError(t, err)
	expected := &TestBinderWithConvertor{
		Options: map[string]string{
			`a`: `1`,
			`b`: `2`,
		},
		Env: []string{
			`A=ONE`,
			`B=TWO`,
		},
	}
	assert.Equal(t, expected, m)
	ctx := e.NewContext(mock.NewRequest(), mock.NewResponse())
	StructToForm(ctx, expected, ``, LowerCaseFirstLetter)
	assert.Equal(t, map[string][]string{
		`options.a`: {"1"},
		`options.b`: {"2"},
		`env`:       {"A=ONE\nB=TWO"},
	}, ctx.Forms())

}
