package oauth2

import (
	"strings"

	"github.com/admpub/goth"
	"github.com/webx-top/echo"
)

// Account 账号信息
type Account struct {
	On          bool // on / off
	Name        string
	Key         string
	Secret      string `json:"-" xml:"-"`
	Extra       echo.H
	LoginURL    string
	CallbackURL string
	Scopes      []string
	Constructor func(*Account) goth.Provider `json:"-" xml:"-"`
}

// SetConstructor sets the provider constructor function for the account and returns the account for chaining.
func (a *Account) SetConstructor(constructor func(*Account) goth.Provider) *Account {
	a.Constructor = constructor
	return a
}

// Instance returns a new goth.Provider instance using the Account's Constructor
func (a *Account) Instance() goth.Provider {
	return a.Constructor(a)
}

// GetCustomisedHostURL retrieves the custom host URL from the account's extra data.
// It trims any trailing slash from the URL before returning it.
// Returns empty string if no custom host URL is configured.
func (a *Account) GetCustomisedHostURL() string {
	hostURL := a.GetExtraString(`hostURL`)
	if len(hostURL) > 0 {
		hostURL = strings.TrimSuffix(hostURL, `/`)
	}
	return hostURL
}

// GetExtraString retrieves a string value from the Account's Extra data by key.
// Returns empty string if Extra is nil or key doesn't exist.
func (a *Account) GetExtraString(key string) string {
	if a.Extra == nil {
		return ``
	}
	return a.Extra.String(key)
}
