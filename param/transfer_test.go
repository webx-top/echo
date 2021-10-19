package param

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTo(t *testing.T) {
	raw := Store{
		`Name`:   `Tester`,
		`Age`:    20,
		`Gender`: `male`,
		`Email`:  `test@webx.top`,
	}
	res := raw.Transform(map[string]Transfer{
		`Name`: &Transform{
			Key: "name",
			Func: func(value interface{}, row Store) interface{} {
				return strings.ToUpper(value.(string))
			},
		},
		`Age`:    &Transform{Key: `age`},
		`Gender`: nil,
	})
	expected := Store{
		`name`:   `TESTER`,
		`age`:    20,
		`Gender`: `male`,
	}
	assert.Equal(t, expected, res)
}

func TestTransferFrom(t *testing.T) {
	raw := Store{
		`Name`:   `Tester`,
		`Age`:    20,
		`Gender`: `male`,
		`Email`:  `test@webx.top`,
	}

	// define transfers
	transefers := NewTransfers()
	transefers.AddFunc(`Name`, func(value interface{}, row Store) interface{} {
		return strings.ToUpper(value.(string))
	}, `name`)
	transefers.AddFunc(`Age`, nil, `age`)
	transefers.Add(`Gender`, nil)

	// transform
	res := transefers.Transform(raw)
	expected := Store{
		`name`:   `TESTER`,
		`age`:    20,
		`Gender`: `male`,
	}
	assert.Equal(t, expected, res)
}

func TestTransferMutilLevel(t *testing.T) {
	raw := Store{
		`Name`: `Tester`,
		`Info`: Store{
			`Age`:    20,
			`Gender`: `male`,
			`Email`:  `test@webx.top`,
			`Other`:  `data`,
		},
		`List`: []Store{
			{`Item1`: `value1`},
			{`Item2`: `value2`},
		},
	}

	// define transfers
	transefers := NewTransfers()
	transefers.AddFunc(`Name`, func(value interface{}, row Store) interface{} {
		return strings.ToUpper(value.(string))
	}, `name`)
	transefers.AddFunc(`Info.Age`, nil, `age`)
	transefers.Add(`Info.Gender`, nil)
	transefers.AddFunc(`Info.Email`, func(value interface{}, row Store) interface{} {
		return strings.ToUpper(value.(string))
	}, `info.email`)
	transefers.AddFunc(`List.Item1`, func(value interface{}, row Store) interface{} {
		if value == nil {
			return ``
		}
		return strings.ToUpper(value.(string))
	}, `list.item1`)

	// transform
	res := transefers.Transform(raw)
	expected := Store{
		`name`: `TESTER`,
		`age`:  20,
		`Info`: Store{
			`Gender`: `male`,
		},
		`info`: Store{
			`email`: `TEST@WEBX.TOP`,
		},
		`list`: []Store{
			{`item1`: `VALUE1`},
			{`item1`: ``},
		},
	}
	assert.Equal(t, expected, res)
}
