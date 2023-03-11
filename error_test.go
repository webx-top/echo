package echo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/code"
)

func TestInErrorCode(t *testing.T) {
	err := NewError(`test`, code.Unauthenticated)
	assert.True(t, IsErrorCode(err, code.Unauthenticated))

	err2 := fmt.Errorf(`err2: %w`, err)
	assert.True(t, IsErrorCode(err2, code.Unauthenticated))
	assert.False(t, IsErrorCode(err2, code.Unsupported))
}
