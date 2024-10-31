package oauth2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompressValue(t *testing.T) {
	v, err := CompressValue(`http://github.com/webx-top/echo`)
	require.NoError(t, err)
	require.Equal(t, true, len(v) > 0)
	t.Log(v)
}
