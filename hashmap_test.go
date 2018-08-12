package echo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapx(t *testing.T) {
	assert.Equal(t, "[a b c d]", fmt.Sprint(FormNames("a[b][c][d]")))
	assert.Equal(t, "[a b c d]", fmt.Sprint(FormNames("a[[b][c][d]")))
	assert.Equal(t, "[a b c d]", fmt.Sprint(FormNames("a][[b][c][d]")))
	assert.Equal(t, "[a  b c d]", fmt.Sprint(FormNames("a[][b][c][d]")))
	data := map[string][]string{
		"a[d]":   []string{"first"},
		"a[e]":   []string{"second"},
		"a[f]":   []string{"third"},
		"a[g]":   []string{"fourth"},
		"b[]":    []string{"index_0", "index_1"},
		"c[][a]": []string{"index 0.a"},
		"c[][b]": []string{"index 1.b"},
	}
	mx := NewMapx(data)
	//Dump(mx)

	assert.Equal(t, "first", mx.Value("a", "d"))
	assert.Equal(t, "second", mx.Value("a", "e"))
	assert.Equal(t, "third", mx.Value("a", "f"))
	assert.Equal(t, "fourth", mx.Value("a", "g"))
	assert.Equal(t, "[index_0 index_1]", fmt.Sprint(mx.Values("b")))
	assert.Equal(t, "index 0.a", mx.Value("c", "0", "a"))
	assert.Equal(t, "index 1.b", mx.Value("c", "1", "b"))
}

func TestKVList(t *testing.T) {
	list := KVList{}
	list.Add(`k`, `v`)
	assert.Equal(t, "k", list[0].K)
	assert.Equal(t, "v", list[0].V)

	list.Add(`k2`, `v2`)
	assert.Equal(t, "k2", list[1].K)
	assert.Equal(t, "v2", list[1].V)

	list.Del(1)
	assert.Equal(t, 1, len(list))

	list.Reset()
	assert.Equal(t, 0, len(list))
}

func TestMapClone(t *testing.T) {
	list := H{
		`key`: `value`,
		`children`: []H{
			{`childKey`: `childValue`},
		},
	}
	listCopy := list.Clone()
	listCopy[`key`] = `valueChg`
	assert.Equal(t, "value", list[`key`])
	assert.Equal(t, "valueChg", listCopy[`key`])

	listCopy[`children`].([]H)[0][`childKey`] = `childValueChg`
	assert.Equal(t, "childValue", list[`children`].([]H)[0][`childKey`])
	assert.Equal(t, "childValueChg", listCopy[`children`].([]H)[0][`childKey`])
}
