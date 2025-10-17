package stripe

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/stripe"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`stripe`, func(account *oauth2.Account) goth.Provider {
		return stripe.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
