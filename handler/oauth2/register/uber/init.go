package uber

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/uber"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`uber`, func(account *oauth2.Account) goth.Provider {
		return uber.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
