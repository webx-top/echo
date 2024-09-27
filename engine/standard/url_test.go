package standard

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestURL(t *testing.T) {
	u, err := url.Parse(`https://github.com/webx-top/echo?test=1&name=测试`)
	assert.NoError(t, err)
	assert.Equal(t, `https://github.com/webx-top/echo?test=1`, u.String())
	echo.Dump(u)
}
