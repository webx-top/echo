package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitea"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`gitea`, func(account *oauth2.Account) goth.Provider {
		return gitea.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
