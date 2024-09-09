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
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/render/driver"
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

	default:
		return nil, fmt.Errorf(`%w: %s`, os.ErrNotExist, name)
	}
}

func TestSlotRender(t *testing.T) {
	a := New(`test`)
	a.SetManager(&testTemplateMgr{})
	a.Init()
	r := a.Fetch(`index`, nil, defaults.NewMockContext())
	assert.Equal(t, `首页 -- powered by webx
这是一个演示`, r)
	r = a.Fetch(`new`, nil, defaults.NewMockContext())
	assert.Equal(t, `首页 -- powered by webx
这是一个插槽演示`, r)
	r = a.Fetch(`new2`, nil, defaults.NewMockContext())
	assert.Equal(t, `首页 -- powered by webx
这是一个[插槽]演示`, r)
}
