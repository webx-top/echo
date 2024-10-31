package oauth2

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webx-top/com"
)

func TestCompressValue(t *testing.T) {
	v, err := CompressValue(`http://github.com/webx-top/echo`)
	require.NoError(t, err)
	require.Equal(t, true, len(v) > 0)
	t.Log(v)

	encrypted := `BKtLthApRvX9Sbeq6+X3MGFT7rN4qBZfW/izn1d6TNjRfI8MHjNNy2bjZLppLvEQkw/hTKw4RxoMkLup54BRrnKN9kQCfyjVyARkRXputHpn4M5LiU7ZgtqfK/nFKTX8xxxx8LFs6916J2oQRoV/DJ3lqTi9R1Tis1EG5netdt5Gn8OgvnP4q4YHWb5Zm44SK5HH6wSOBmIi1+ni0AA8UBfjm/HMt+pJZfd093PRIHMvGD9dPeS1Oxd3U8/AEqW0Q5AWhKgqP9rtUqsWOktT/24CVjAHxwY3TIDSwrXFmil+tEIkttOnG1VjpfEUygO84foGAFVQOnCPzxYlIYYnLw==`

	v, err = com.Base64Decode(encrypted)
	require.NoError(t, err)
	t.Log(v)

	encrypted = url.QueryEscape(encrypted)
	t.Log(encrypted)
	encrypted, err = url.QueryUnescape(encrypted)
	require.NoError(t, err)
	v, err = com.Base64Decode(encrypted)
	require.NoError(t, err)
	t.Log(v)

	encrypted, err = url.QueryUnescape(encrypted)
	require.NoError(t, err)
	v, err = com.Base64Decode(encrypted)
	require.Error(t, err)
	t.Log(v)
}
