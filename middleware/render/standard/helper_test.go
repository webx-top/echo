package standard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceByMatchedIndex(t *testing.T) {
	a := New(`test`).(*Standard)
	content := `<body>
    {{Include "header"}}
	<div id="wrapper">
		{{Include "sidebar"}}
		<div class="body">
		</div>
	</div>
</body>`
	matches := a.incTagRegex.FindAllStringSubmatchIndex(content, -1)
	var replaced string
	fn := replaceByMatchedIndex(content, matches, &replaced)
	for k, v := range matches {
		var tmplFile, passObject string
		getMatchedByIndex(content, v, nil, &tmplFile, &passObject)
		if k == 0 {
			assert.Equal(t, `header`, tmplFile)
			assert.Equal(t, ``, passObject)
		} else {
			assert.Equal(t, `sidebar`, tmplFile)
			assert.Equal(t, ``, passObject)
		}
		fn(k, v, `{P}`)
	}
	expected := `<body>
    {P}
	<div id="wrapper">
		{P}
		<div class="body">
		</div>
	</div>
</body>`
	assert.Equal(t, expected, replaced)
	var replaced2 string
	fn2 := replaceByMatchedIndex(content, matches, &replaced2)
	for k, v := range matches {
		fn2(k, v)
	}
	assert.Equal(t, content, replaced2)
}

func TestReplaceByMatchedIndex2(t *testing.T) {
	a := New(`test`).(*Standard)
	content := `{{Include "sub"}}`
	matches := a.incTagRegex.FindAllStringSubmatchIndex(content, -1)
	var replaced string
	fn := replaceByMatchedIndex(content, matches, &replaced)
	for k, v := range matches {
		var tmplFile, passObject string
		getMatchedByIndex(content, v, nil, &tmplFile, &passObject)
		assert.Equal(t, `sub`, tmplFile)
		assert.Equal(t, ``, passObject)
		fn(k, v, `{P}`)
	}
	expected := `{P}`
	assert.Equal(t, expected, replaced)
	var replaced2 string
	fn2 := replaceByMatchedIndex(content, matches, &replaced2)
	for k, v := range matches {
		fn2(k, v)
	}
	assert.Equal(t, content, replaced2)
}

func TestSliceClip(t *testing.T) {
	v := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, []int{}, v[0:0])
	assert.Equal(t, []int{6}, v[len(v)-1:len(v)]) // OK
	assert.Equal(t, []int{}, v[len(v):])          // OK
	// v[len(v)] =>  panic: runtime error: index out of range [6] with length 6
	// 从切片中裁切子切片时可以使用 len(v) 的值作为下标，而通过下标取单个元素值时则下标的值不能大于 len(v)-1
}
