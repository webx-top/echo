package apple

import (
	"net/http"

	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/apple"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`apple`, func(account *oauth2.Account) goth.Provider {
		return apple.New(account.Key, account.Secret, account.CallbackURL, http.DefaultClient, account.Scopes...)
	})
}
