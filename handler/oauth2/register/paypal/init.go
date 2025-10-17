package paypal

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/paypal"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`paypal`, func(account *oauth2.Account) goth.Provider {
		return paypal.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
