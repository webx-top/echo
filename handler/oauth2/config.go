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
	"github.com/imdario/mergo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/webx-top/echo"
)

const (
	// DefaultPath /oauth
	DefaultPath = "/oauth"

	// DefaultContextKey oauth_user
	DefaultContextKey = "oauth_user"
)

type Account struct {
	Name        string
	Key         string
	Secret      string
	Extra       echo.H
	Constructor func(*Account) goth.Provider `json:"-" xml:"-"`
}

func (a *Account) SetConstructor(constructor func(*Account) goth.Provider) {
	a.Constructor = constructor
}

func (a *Account) Instance() goth.Provider {
	return a.Constructor(a)
}

// Config the configs for the gothic oauth/oauth2 authentication for third-party websites
// All Key and Secret values are empty by default strings. Non-empty will be registered as Goth Provider automatically, by Iris
// the users can still register their own providers using goth.UseProviders
// contains the providers' keys  (& secrets) and the relative auth callback url path(ex: "/auth" will be registered as /auth/:provider/callback)
//
type Config struct {
	Host, Path string
	Accounts   []*Account

	// defaults to 'oauth_user' used by plugin to give you the goth.User, but you can take this manually also by `context.Get(ContextKey).(goth.User)`
	ContextKey string
}

// DefaultConfig returns OAuth config, the fields of the iteral are zero-values ( empty strings)
func DefaultConfig() *Config {
	return &Config{
		Path:       DefaultPath,
		Accounts:   []*Account{},
		ContextKey: DefaultContextKey,
	}
}

// MergeSingle merges the default with the given config and returns the result
func (c *Config) MergeSingle(cfg *Config) (config *Config) {
	config = cfg
	mergo.Merge(config, c)
	return
}

func (c *Config) CallbackURL(providerName string) string {
	return c.Host + c.Path + "/callback/" + providerName
}

// GenerateProviders returns the valid goth providers and the relative url paths (because the goth.Provider doesn't have a public method to get the Auth path...)
// we do the hard-core/hand checking here at the configs.
//
// receives one parameter which is the host from the server,ex: http://localhost:3000, will be used as prefix for the oauth callback
func (c *Config) GenerateProviders() *Config {
	goth.ClearProviders()
	var providers []goth.Provider
	//we could use a map but that's easier for the users because of code completion of their IDEs/editors
	for _, account := range c.Accounts {
		if len(account.Key) == 0 || len(account.Secret) == 0 {
			continue
		}
		if account.Constructor != nil {
			providers = append(providers, account.Instance())
			continue
		}
		switch account.Name {
		case "twitter":
			providers = append(providers, twitter.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "facebook":
			providers = append(providers, facebook.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "gplus":
			providers = append(providers, gplus.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "github":
			providers = append(providers, github.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "spotify":
			providers = append(providers, spotify.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "linkedin":
			providers = append(providers, linkedin.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "lastfm":
			providers = append(providers, lastfm.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "twitch":
			providers = append(providers, twitch.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "dropbox":
			providers = append(providers, dropbox.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "digitalocean":
			providers = append(providers, digitalocean.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "bitbucket":
			providers = append(providers, bitbucket.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "instagram":
			providers = append(providers, instagram.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "box":
			providers = append(providers, box.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "salesforce":
			providers = append(providers, salesforce.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "amazon":
			providers = append(providers, amazon.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "yammer":
			providers = append(providers, yammer.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "onedrive":
			providers = append(providers, onedrive.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "yahoo":
			providers = append(providers, yahoo.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "slack":
			providers = append(providers, slack.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "stripe":
			providers = append(providers, stripe.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "wepay":
			providers = append(providers, wepay.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "paypal":
			providers = append(providers, paypal.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "steam":
			providers = append(providers, steam.New(account.Key, c.CallbackURL(account.Name)))
		case "heroku":
			providers = append(providers, heroku.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "uber":
			providers = append(providers, uber.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "soundcloud":
			providers = append(providers, soundcloud.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		case "gitlab":
			providers = append(providers, gitlab.New(account.Key, account.Secret, c.CallbackURL(account.Name)))
		}
	}

	goth.UseProviders(providers...)
	return c
}

func (c *Config) AddAccount(accounts ...*Account) *Config {
	c.Accounts = append(c.Accounts, accounts...)
	return c
}
