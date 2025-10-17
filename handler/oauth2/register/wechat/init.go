package wechat

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/wechat"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauth2.Register(`wechat`, func(account *oauth2.Account) goth.Provider {
		return wechat.New(account.Key, account.Secret, account.CallbackURL, wechat.WECHAT_LANG_CN)
	})
}
