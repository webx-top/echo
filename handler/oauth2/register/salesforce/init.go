package salesforce

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/salesforce"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`salesforce`, func(account *oauth2.Account) goth.Provider {
		return salesforce.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
