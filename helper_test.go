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
