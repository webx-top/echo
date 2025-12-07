package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitea"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`gitea`, func(account *oauth2.Account) goth.Provider {
		hostURL := account.GetCustomisedHostURL()
		if len(hostURL) > 0 {
			return gitea.NewCustomisedURL(account.Key, account.Secret, account.CallbackURL, hostURL+`/login/oauth/authorize`, hostURL+`/login/oauth/access_token`, hostURL+`/api/v1/user`, account.Scopes...)
		}
		return gitea.New(account.Key, account.Secret, account.CallbackURL, account.Scopes...)
	})
}
