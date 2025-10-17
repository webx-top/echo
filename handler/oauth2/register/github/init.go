package github

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/github"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`github`, func(account *oauth2.Account) goth.Provider {
		return github.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
