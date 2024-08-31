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
