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

	encrypted := `BCILe6uQCv9HsJeSql0IHITsI6B08bFU198fSoKC90rijfTkYFp3Ks8jxy4bkj1Lk+pzRCp4BWOurEuLZv9zFtVW+wNUPJUsUgI6iBcXdrk/IwgXjekZO1dSFysXl+OQMctiySC7VuCRO1Z6tI3UvUQCx8HZuYTEJ+GKiUwWNrEgjFICJPS09ztUW4aMXYskpFVEG1BfPBfyBjwUyumjDbAOSUJrTrn31bps4/XFctfc5PzHrbU9CBmg4Ps56UGMqZdw8svc8mRbwuKWb22qHmKg4BARKz/zRJCMOICMOOS7b2pU5L2GrWvOMMunVfrjUJ/XLj0XtqKHl/icfIJ13ksruLfDfIHzhaAPFGNLSYp02TGeT5WhjXBCdoMsoDp46yuYSeIGkoudCTh6cp2JctazUo/Fy//QJMry8lsr49ugu2PNrfQYqcRCBff26MLQvvpJWe3YHc+lQNsUAlNiFDGY665S9cyFB6V7cs7SSi+1k/1F27n6WbgpR9HAq+LwqmjuhyD/xESALZjoLI3iHtc/80C9sf24Tl8vqTPJU3KC`

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
