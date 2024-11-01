/*
Copyright 2016 Wenhui Shen <www.webx.top>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package oauth2

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/goth"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// SessionName is the key used to access the session store.
// we could use the echo's sessions default, but this session should be not confict with the cookie session name defined by the sessions manager
const SessionName = "EchoGoth"
const StateSessionName = "EchoGothState"
const redirectURIQueryString = `redirect_uri=%2F`

var (
	_         goth.Params = (*url.Values)(nil)
	EmptyUser             = goth.User{}
)

func fixedRedirectURIQueryString(ctx echo.Context, authURL string) string {
	pos := strings.Index(authURL, redirectURIQueryString)
	if pos > 0 {
		authURL = authURL[0:pos] + `redirect_uri=` + url.QueryEscape(ctx.Site()) + authURL[pos+len(redirectURIQueryString):]
	}
	return authURL
}

/*
BeginAuthHandler is a convienence handler for starting the authentication process.
It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".
BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.
*/
func BeginAuthHandler(ctx echo.Context) error {
	authURL, err := GetAuthURL(ctx)
	if err != nil {
		return echo.NewHTTPError(400, err.Error()).SetRaw(err)
	}
	next := ctx.Form(echo.DefaultNextURLVarName)
	if len(next) > 0 {
		ctx.Cookie().Set(echo.DefaultNextURLVarName, next)
	}
	return ctx.Redirect(authURL)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx echo.Context) (string, error) {
	state := ctx.Query("state")
	if len(state) > 0 {
		return state, nil
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		err = fmt.Errorf("gothic: source of randomness unavailable: %w", err)
		return state, err
	}
	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx echo.Context) string {
	state := ctx.Query("state")
	if len(state) == 0 && ctx.IsPost() {
		state = ctx.Request().FormValue("state")
	}
	return state
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or url query parameter ":provider".
I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx echo.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	state, err := SetState(ctx)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(state)
	if err != nil {
		return "", err
	}

	if cr, ok := sess.(echo.ContextRegister); ok {
		cr.SetContext(ctx)
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	authURL = fixedRedirectURIQueryString(ctx, authURL)
	authURL = fixedURL(ctx, authURL)
	sessData, err := EncryptValue(ctx, sess.Marshal())
	if err != nil {
		return "", err
	}
	state, err = EncryptValue(ctx, state)
	if err != nil {
		return "", err
	}
	ctx.Cookie().Set(SessionName, sessData).Set(StateSessionName, state)
	return authURL, err
}

func fixedURL(ctx echo.Context, url string) string {
	length := len(url)
	if length == 0 {
		return url
	}
	switch url[0] {
	case '/':
		url = ctx.Site() + strings.TrimPrefix(url, `/`)
	case '.':
		url = ctx.Site() + url
	default:
		if length > 7 {
			switch url[0:7] {
			case `https:/`, `http://`:
				return url
			}
		}
		url = ctx.Site() + url
	}
	return url
}

func fetchUser(ctx echo.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return EmptyUser, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return EmptyUser, err
	}

	sv, ok := ctx.Session().Get(SessionName).(string)
	if !ok || len(sv) == 0 {
		return EmptyUser, ErrSessionDismatched
	}

	defer func() {
		if err != nil {
			ctx.Session().Delete(SessionName).Save()
		}
	}()

	var sess goth.Session
	sess, err = provider.UnmarshalSession(sv)
	if err != nil {
		return EmptyUser, err
	}

	if cr, ok := sess.(echo.ContextRegister); ok {
		cr.SetContext(ctx)
	}

	var user goth.User
	user, err = provider.FetchUser(sess)
	if err != nil {
		return EmptyUser, err
	}
	// user can be found with existing session data
	return user, err
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.
It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".
*/
var CompleteUserAuth = func(ctx echo.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return EmptyUser, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return EmptyUser, err
	}

	//error=invalid_request&error_description=The provided value for the input parameter 'redirect_uri' is not valid. The scope 'openid offline_access user.read' requires that the request must be sent over a secure connection using SSL.&state=state
	errorDescription := ctx.Query(`error_description`)
	if len(errorDescription) > 0 {
		return EmptyUser, errors.New(providerName + `: ` + errorDescription)
	}

	sv, err := DecryptValue(ctx, ctx.Cookie().Get(SessionName))
	if err != nil {
		return EmptyUser, fmt.Errorf(providerName+`: %w`, err)
	}
	if len(sv) == 0 {
		return EmptyUser, ErrSessionDismatched
	}

	defer func() {
		ctx.Cookie().Set(SessionName, ``, -1)
	}()

	var sess goth.Session
	sess, err = provider.UnmarshalSession(sv)
	if err != nil {
		return EmptyUser, err
	}

	if cr, ok := sess.(echo.ContextRegister); ok {
		cr.SetContext(ctx)
	}

	err = validateState(ctx, sess)
	if err != nil {
		return EmptyUser, err
	}

	var user goth.User
	user, err = provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	params := ctx.Queries()
	if len(params) == 0 && ctx.IsPost() {
		params = ctx.Request().PostForm().All()
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, url.Values(params))
	if err != nil {
		return EmptyUser, err
	}

	ctx.Session().Set(SessionName, sess.Marshal()).Save()

	user, err = provider.FetchUser(sess)
	return user, err
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(ctx echo.Context, _ goth.Session) error {
	originalState, err := DecryptValue(ctx, ctx.Cookie().Get(StateSessionName))
	if err != nil || len(originalState) == 0 || originalState != GetState(ctx) {
		return ErrStateTokenMismatch
	}
	ctx.Cookie().Set(StateSessionName, ``, -1)
	return nil
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = getProviderName

func getProviderName(ctx echo.Context) (string, error) {
	provider := ctx.Param("provider")
	if len(provider) > 0 {
		return provider, nil
	}
	provider = ctx.Query("provider")
	if len(provider) == 0 {
		return provider, ErrMustSelectProvider
	}
	return provider, nil
}

func PackValue(ip, ua, value string) string {
	return ip + `|` + com.Md5(ua) + `@` + com.String(time.Now().Unix()) + `|` + value
}

func UnpackValue(ip, ua, value string, maxAge time.Duration) (string, error) {
	parts := strings.SplitN(value, `|`, 3)
	if len(parts) != 3 {
		return "", nil
	}
	if parts[0] != ip {
		return "", fmt.Errorf(`%w: %q != %q`, ErrIPAddressDismatched, parts[0], ip)
	}
	parts2 := strings.SplitN(parts[1], `@`, 2)
	if len(parts2) == 2 {
		if ts := com.Int64(parts2[1]); ts > 0 {
			ti := time.Unix(ts, 0)
			if ti.Before(time.Now().Add(-maxAge)) {
				return "", fmt.Errorf(`%w: %s (maxAge: %s)`, ErrDataExpired, ti.Format(time.DateTime), maxAge.String())
			}
		}
	}
	if uaMd5 := com.Md5(ua); parts2[0] != uaMd5 {
		return "", fmt.Errorf(`%w: %q != %q`, ErrUserAgentDismatched, parts2[0], uaMd5)
	}
	return parts[2], nil
}

func EncryptValue(ctx echo.Context, value string) (string, error) {
	if len(value) == 0 {
		return value, nil
	}
	var err error
	value = PackValue(ctx.RealIP(), ctx.Request().UserAgent(), value)
	value, err = CompressValue(value)
	if err != nil {
		return "", err
	}
	if ctx.CookieOptions().Cryptor != nil {
		value, err = ctx.CookieOptions().Cryptor.EncryptString(value)
		if err != nil {
			return "", err
		}
		value = com.URLSafeBase64(value, true)
	} else {
		value = url.QueryEscape(value)
	}
	return value, err
}

func DecryptValue(ctx echo.Context, value string) (string, error) {
	if len(value) == 0 {
		return value, nil
	}
	var err error
	if ctx.CookieOptions().Cryptor != nil {
		value = com.URLSafeBase64(value, false)
		value, err = ctx.CookieOptions().Cryptor.DecryptString(value)
		if err != nil {
			return "", err
		}
	} else {
		value, err = url.QueryUnescape(value)
		if err != nil {
			return "", err
		}
	}
	value, err = UncompressValue(value)
	if err != nil {
		return "", err
	}
	var unpackErr error
	value, unpackErr = UnpackValue(ctx.RealIP(), ctx.Request().UserAgent(), value, time.Hour)
	if unpackErr != nil {
		log.Warn(unpackErr.Error())
	}
	return value, err
}

func UncompressValue(value string) (string, error) {
	rdata := strings.NewReader(value)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return com.Bytes2str(s), nil
}

func CompressValue(value string) (val string, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err = gz.Write(com.Str2bytes(value)); err != nil {
		return
	}
	if err = gz.Flush(); err != nil {
		return
	}
	if err = gz.Close(); err != nil {
		return
	}
	val = b.String()
	return
}
