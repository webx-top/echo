package echo

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHandlerFunc(ctx Context) error {
	return nil
}

type testHandler struct {
}

func (t *testHandler) Handle(ctx Context) error {
	return nil
}

func TestHandlerPath(t *testing.T) {
	ppath := HandlerPath(testHandlerFunc)
	assert.Equal(t, "github.com/webx-top/echo.testHandlerFunc", ppath)
	ppath = HandlerPath(HandlerFunc(testHandlerFunc))
	assert.Equal(t, "github.com/webx-top/echo.testHandlerFunc", ppath)
	ppath = HandlerPath(&testHandler{})
	assert.Equal(t, "github.com/webx-top/echo.testHandler", ppath)
	ppath = HandlerTmpl(`github.com/webx-top/echo.(*TestHandler).Index-fm`)
	assert.Equal(t, "/echo/test_handler/index", ppath)
}

func TestLogIf(t *testing.T) {
	LogIf(errors.New(`test`), `debug`)
}

func TestFileSafePath(t *testing.T) {
	full, err := FilePathJoin(`abc`, `../../abc/`)
	assert.NoError(t, err)
	assert.Equal(t, `abc/abc`, full)

	assert.True(t, PathHasDots(`..`))
	assert.True(t, PathHasDots(`../`))
	assert.True(t, PathHasDots(`/../`))
	assert.True(t, PathHasDots(`/..`))
	assert.False(t, PathHasDots(`/..adb`))
	assert.True(t, PathHasDots(`cc/../adb`))
	assert.True(t, PathHasDots(`..\`))
	assert.True(t, PathHasDots(`\..\`))
	assert.True(t, PathHasDots(`\..`))
	assert.False(t, PathHasDots(`\..adb`))
	assert.True(t, PathHasDots(`cc\..\adb`))
}

func TestCreateInRoot(t *testing.T) {
	dir := `./testdata`
	os.Mkdir(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	f, err := CreateInRoot(dir, `abc.txt`)
	assert.NoError(t, err)
	f.WriteString(`test`)
	f.Close()

	err = WriteFileInRoot(dir, `def.txt`, []byte(`test1`), os.ModePerm)
	assert.NoError(t, err)

	b, err := ReadFileInRoot(dir, `def.txt`)
	assert.NoError(t, err)
	assert.Equal(t, `test1`, string(b))
}

func TestURLEncode(t *testing.T) {
	raw := `1 2?a=b`
	encoded := URLEncode(raw)
	assert.Equal(t, "1+2%3Fa%3Db", encoded)
	content, _ := URLDecode(encoded)
	assert.Equal(t, raw, content)
	encoded = URLEncode(raw, true)
	assert.Equal(t, "1%202%3Fa%3Db", encoded)
	content, _ = URLDecode(encoded, true)
	assert.Equal(t, raw, content)
}

func TestInSliceFold(t *testing.T) {
	assert.True(t, InSliceFold(`post`, []string{`POST`}))
}

func TestParseTemplateError(t *testing.T) {
	content := `template: /Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html:7:831: executing "/Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html" at <call>: wrong number of args for call: want at least 1 got 0`
	matches := regErrorTemplateFile.FindAllStringSubmatch(content, -1)
	assert.Equal(t, `template: /Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html:7:831: `, matches[0][0])
	assert.Equal(t, `/Users/hank/go/src/github.com/admpub/nging/template/backend/manager/role_edit_perm_page.html`, matches[0][1])
	assert.Equal(t, `7`, matches[0][2])
	assert.Equal(t, `831`, matches[0][3])
	//panic(Dump(matches, false))

	content = `template: /Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html:6: function "Now2" not defined`
	matches = regErrorTemplateFile.FindAllStringSubmatch(content, -1)
	assert.Equal(t, `template: /Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html:6: `, matches[0][0])
	assert.Equal(t, `/Users/hank/go/src/github.com/webx-top/echo/middleware/render/standard/test/template/test.html`, matches[0][1])
	assert.Equal(t, `6`, matches[0][2])
}
