package standard

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	for i := 0; i < 10; i++ {
		go func() {
			buf := bytes.NewBuffer(nil)
			err = tmpl.Execute(buf, nil)
			assert.NoError(t, err)
			//assert.Equal(t, now, buf.String())
			time.Sleep(time.Second * time.Duration(rand.Intn(5)+1))
		}()
	}
	time.Sleep(time.Second * 5)
}

func TestParseError(t *testing.T) {
	content := `template: /Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html:7:831: executing "/Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html" at <call>: wrong number of args for call: want at least 1 got 0`
	matches := regErrorFile.FindAllStringSubmatch(content, -1)
	assert.Equal(t, `template: /Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html:7:831: `, matches[0][0])
	assert.Equal(t, `/Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html`, matches[0][1])
	assert.Equal(t, `7`, matches[0][2])
	assert.Equal(t, `831`, matches[0][3])
	//panic(echo.Dump(matches, false))

	content = `template: /Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html:6: function "Now2" not defined`
	matches = regErrorFile.FindAllStringSubmatch(content, -1)
	assert.Equal(t, `template: /Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html:6: `, matches[0][0])
	assert.Equal(t, `/Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html`, matches[0][1])
	assert.Equal(t, `6`, matches[0][2])
}
