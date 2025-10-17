package paypal

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/paypal"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`paypal`, func(cfg *oauth2.Config) goth.Provider {
		return paypal.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
