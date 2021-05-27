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
