package echo_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func testHandlerFunc(ctx echo.Context) error {
	return nil
}

type testHandler struct {
}

func (t *testHandler) Handle(ctx echo.Context) error {
	return nil
}

func TestHandlerPath(t *testing.T) {
	ppath := echo.HandlerPath(testHandlerFunc)
	assert.Equal(t, "github.com/webx-top/echo_test.testHandlerFunc", ppath)
	ppath = echo.HandlerPath(echo.HandlerFunc(testHandlerFunc))
	assert.Equal(t, "github.com/webx-top/echo_test.testHandlerFunc", ppath)
	ppath = echo.HandlerPath(&testHandler{})
	assert.Equal(t, "github.com/webx-top/echo_test.testHandler", ppath)
	ppath = echo.HandlerTmpl(`github.com/webx-top/echo_test.(*TestHandler).Index-fm`)
	assert.Equal(t, "/echo_test/test_handler/index", ppath)
}

func TestLogIf(t *testing.T) {
	echo.LogIf(errors.New(`test`), `debug`)
}

func TestURLEncode(t *testing.T) {
	raw := `1 2?a=b`
	encoded := echo.URLEncode(raw)
	assert.Equal(t, "1+2%3Fa%3Db", encoded)
	content, _ := echo.URLDecode(encoded)
	assert.Equal(t, raw, content)
	encoded = echo.URLEncode(raw, true)
	assert.Equal(t, "1%202%3Fa%3Db", encoded)
	content, _ = echo.URLDecode(encoded, true)
	assert.Equal(t, raw, content)
}
