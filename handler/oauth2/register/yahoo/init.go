package yahoo

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/yahoo"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`yahoo`, func(account *oauth2.Account) goth.Provider {
		return yahoo.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
