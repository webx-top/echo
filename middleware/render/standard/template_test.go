package standard

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func TestTemplate(t *testing.T) {
	tmpl := template.New(`test`)
	var now string
	funcMap := map[string]interface{}{
		`now`: func() template.HTML {
			now = time.Now().Format(time.RFC3339Nano)
			fmt.Println(now)
			return template.HTML(now)
		},
	}
	tmpl.Funcs(funcMap)
	_, err := tmpl.Parse(`{{now}}`)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := bytes.NewBuffer(nil)
			err = tmpl.Execute(buf, nil)
			assert.NoError(t, err)
			//assert.Equal(t, now, buf.String())
			time.Sleep(time.Second * time.Duration(rand.Intn(5)+1))
		}()
	}
	wg.Wait()
}

type testTemplateMgr struct {
	driver.BaseManager
}

func (b *testTemplateMgr) GetTemplate(name string) ([]byte, error) {
	switch filepath.Base(name) {
	case `layout.html`:
		return []byte(`{{Block "title"}}-- powered by webx{{/Block}}
{{Block "body"}}内容{{/Block}}`), nil

	case `index.html`:
		return []byte(`{{Extend "layout"}}
{{Block "title"}}首页 {{Super}}{{/Block}}
{{Block "body"}}这是一个{{Block "demoName"/}}演示{{/Block}}`), nil

	case `new.html`:
		return []byte(`{{Extend "index"}}
{{Block "demoName"}}插槽{{/Block}}`), nil

	case `new2.html`:
		return []byte(`{{Extend "index"}}
{{Block "demoName"}}[插槽]{{/Block}}`), nil

	case `snippet.html`:
		return []byte(`{{Extend "layout"}}
{{Block "body"}}
{{Snippet "testSnippet" A}}
{{/Block}}`), nil

	default:
		return nil, fmt.Errorf(`%w: %s`, os.ErrNotExist, name)
	}
}

func TestSlotRender(t *testing.T) {
	a := New(`test`)
	a.SetManager(&testTemplateMgr{})
	a.Init()
	a.SetFuncMap(func() map[string]interface{} {
		return tplfunc.New()
	})
	ctx := defaults.NewMockContext()
	r := a.Fetch(`index`, nil, ctx)
	assert.Equal(t, `首页 -- powered by webx
这是一个演示`, r)
	r = a.Fetch(`new`, nil, ctx)
	assert.Equal(t, `首页 -- powered by webx
这是一个插槽演示`, r)
	r = a.Fetch(`new2`, nil, ctx)
	assert.Equal(t, `首页 -- powered by webx
这是一个[插槽]演示`, r)

	ctx.SetFunc(`testSnippet`, func(ctx echo.Context, tmpl string, arg string) string {
		return `{{"Z"|ToLower}}@` + arg
	})
	r = a.Fetch(`snippet`, nil, ctx)
	assert.Equal(t, `-- powered by webx
z@A`, r)
}

// go test -bench=BenchmarkXxx
func BenchmarkXxx(b *testing.B) {
	a := New(`test`)
	a.SetManager(&testTemplateMgr{})
	a.Init()
	a.SetFuncMap(func() map[string]interface{} {
		return tplfunc.New()
	})
	b.Run("fetch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx := defaults.NewMockContext()
			a.Fetch(`index`, nil, ctx)
			a.Fetch(`new`, nil, ctx)
			a.Fetch(`new2`, nil, ctx)
			a.Fetch(`snippet`, nil, ctx)
		}
	})

	b.ReportAllocs()
}
