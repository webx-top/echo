package amazon

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/amazon"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`amazon`, func(account *oauth2.Account) goth.Provider {
		return amazon.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
