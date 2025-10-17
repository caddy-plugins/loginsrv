package wechat

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/wechat"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`wechat`, func(cfg *oauth2.Config) goth.Provider {
		return wechat.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), wechat.WECHAT_LANG_CN)
	})
}
