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
	//"net/http"
	"sync"

	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/amazon"
	//"github.com/admpub/goth/providers/apple"
	"github.com/admpub/goth/providers/bitbucket"
	"github.com/admpub/goth/providers/gitea"
	"github.com/admpub/goth/providers/github"
	"github.com/admpub/goth/providers/paypal"
	"github.com/admpub/goth/providers/salesforce"
	"github.com/admpub/goth/providers/stripe"
	"github.com/admpub/goth/providers/uber"
	"github.com/admpub/goth/providers/wechat"
	"github.com/admpub/goth/providers/yahoo"
	"github.com/webx-top/echo"
)

const (
	// DefaultPath /oauth
	DefaultPath = "/oauth"

	// DefaultContextKey oauth_user
	DefaultContextKey = "oauth_user"
)

// Config the configs for the gothic oauth/oauth2 authentication for third-party websites
// All Key and Secret values are empty by default strings. Non-empty will be registered as Goth Provider automatically, by Iris
// the users can still register their own providers using goth.UseProviders
// contains the providers' keys  (& secrets) and the relative auth callback url path(ex: "/auth" will be registered as /auth/:provider/callback)
type Config struct {
	Host, Path string
	accounts   []*Account
	accountM   map[string]int
	mu         *sync.RWMutex

	// defaults to 'oauth_user' used by plugin to give you the goth.User, but you can take this manually also by `context.Get(ContextKey).(goth.User)`
	ContextKey string
}

func NewConfig() *Config {
	return &Config{
		Path:       DefaultPath,
		accounts:   []*Account{},
		accountM:   map[string]int{},
		mu:         &sync.RWMutex{},
		ContextKey: DefaultContextKey,
	}
}

// DefaultConfig returns OAuth config, the fields of the iteral are zero-values ( empty strings)
func DefaultConfig() *Config {
	return &Config{
		Path:       DefaultPath,
		accounts:   []*Account{},
		accountM:   map[string]int{},
		mu:         &sync.RWMutex{},
		ContextKey: DefaultContextKey,
	}
}

// MergeSingle merges the default with the given config and returns the result
func (c *Config) MergeSingle(cfg *Config) (config Config) {
	config = *cfg
	if len(config.Host) == 0 {
		config.Host = c.Host
	}
	if len(config.Path) == 0 {
		config.Path = c.Path
	}
	if config.mu == nil {
		config.mu = c.mu
	}
	if len(config.accounts) == 0 {
		c.mu.RLock()
		config.accounts = make([]*Account, len(c.accounts))
		for k, v := range c.accounts {
			copyV := *v
			if len(v.Extra) > 0 {
				copyV.Extra = v.Extra.Clone()
			} else {
				copyV.Extra = echo.H{}
			}
			config.accounts[k] = &copyV
		}
		c.mu.RUnlock()
	}
	if len(config.ContextKey) == 0 {
		config.ContextKey = c.ContextKey
	}
	return
}

func (c *Config) CallbackURL(providerName string) string {
	return c.Host + c.Path + "/callback/" + providerName
}

func (c *Config) LoginURL(providerName string) string {
	return c.Host + c.Path + "/login/" + providerName
}

// GenerateProviders returns the valid goth providers and the relative url paths (because the goth.Provider doesn't have a public method to get the Auth path...)
// we do the hard-core/hand checking here at the configs.
//
// receives one parameter which is the host from the server,ex: http://localhost:3000, will be used as prefix for the oauth callback
func (c *Config) GenerateProviders() *Config {
	var providers []goth.Provider
	//we could use a map but that's easier for the users because of code completion of their IDEs/editors
	c.RangeAccounts(func(account *Account) bool {
		if !account.On {
			return true
		}
		if provider := c.NewProvider(account); provider != nil {
			providers = append(providers, provider)
		}
		return true
	})
	goth.UseProviders(providers...)
	return c
}

func (c *Config) ClearAccounts() {
	c.mu.Lock()
	c.accounts = []*Account{}
	c.accountM = map[string]int{}
	c.mu.Unlock()
}

func (c *Config) RangeAccounts(cb func(*Account) bool) (ok bool) {
	c.mu.RLock()
	for _, account := range c.accounts {
		ok = cb(account)
		if !ok {
			break
		}
	}
	c.mu.RUnlock()
	return
}

func (c *Config) ClearProviders() *Config {
	goth.ClearProviders()
	return c
}

func (c *Config) DeleteProvider(names ...string) *Config {
	goth.DeleteProvider(names...)
	return c
}

func (c *Config) NewProvider(account *Account) goth.Provider {
	if len(account.LoginURL) == 0 {
		account.LoginURL = c.LoginURL(account.Name)
	}
	if len(account.CallbackURL) == 0 {
		account.CallbackURL = c.CallbackURL(account.Name)
	}
	if account.Constructor != nil {
		return account.Instance()
	}
	switch account.Name {
	case "gitea":
		return gitea.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "github":
		return github.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "bitbucket":
		return bitbucket.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "salesforce":
		return salesforce.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "amazon":
		return amazon.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "yahoo":
		return yahoo.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "stripe":
		return stripe.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	case "paypal":
		return paypal.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	//case "apple":
	//	return apple.New(account.Key, account.Secret, account.CallbackURL, http.DefaultClient, account.Scopes...)
	case "wechat":
		return wechat.New(account.Key, account.Secret, account.CallbackURL, wechat.WECHAT_LANG_CN)
	case "uber":
		return uber.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	}
	return nil
}

func (c *Config) AddAccount(accounts ...*Account) *Config {
	c.mu.Lock()
	for _, account := range accounts {
		if idx, ok := c.accountM[account.Name]; ok {
			c.accounts[idx] = account
			continue
		}
		c.accountM[account.Name] = len(c.accounts)
		c.accounts = append(c.accounts, account)
	}
	c.mu.Unlock()
	return c
}

func (c *Config) GetAccount(name string) (account *Account) {
	c.mu.RLock()
	idx, ok := c.accountM[name]
	if ok {
		account = c.accounts[idx]
	}
	c.mu.RUnlock()
	return
}

func (c *Config) DeleteAccount(name string) {
	c.mu.Lock()
	idx, ok := c.accountM[name]
	if ok {
		accounts := make([]*Account, 0, len(c.accounts))
		if idx > 0 {
			accounts = append(accounts, c.accounts[0:idx]...)
		}
		if len(c.accounts) > idx+1 {
			for _, account := range c.accounts[idx+1:] {
				c.accountM[account.Name] = len(c.accounts)
				accounts = append(accounts, account)
			}
		}
		c.accounts = accounts
		delete(c.accountM, name)
	}
	c.mu.Unlock()
}

func (c *Config) SetAccount(newAccount *Account) *Config {
	c.mu.Lock()
	idx, ok := c.accountM[newAccount.Name]
	if ok {
		account := c.accounts[idx]
		isOff := account.On && !newAccount.On
		account.On = newAccount.On
		account.Key = newAccount.Key
		account.Secret = newAccount.Secret
		account.Extra = newAccount.Extra
		account.Constructor = newAccount.Constructor
		account.LoginURL = newAccount.LoginURL
		account.CallbackURL = newAccount.CallbackURL
		account.Scopes = make([]string, len(newAccount.Scopes))
		copy(account.Scopes, newAccount.Scopes)
		c.accounts[idx] = account
		if isOff {
			c.DeleteProvider(account.Name)
		} else if account.On {
			if provider := c.NewProvider(account); provider != nil {
				goth.UseProviders(provider)
			}
		}
	} else {
		c.accountM[newAccount.Name] = len(c.accounts)
		c.accounts = append(c.accounts, newAccount)
		if newAccount.On {
			if provider := c.NewProvider(newAccount); provider != nil {
				goth.UseProviders(provider)
			}
		}
	}
	c.mu.Unlock()
	return c
}
