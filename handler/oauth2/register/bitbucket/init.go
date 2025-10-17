package bitbucket

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/bitbucket"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`bitbucket`, func(account *oauth2.Account) goth.Provider {
		return bitbucket.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
