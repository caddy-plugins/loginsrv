package stripe

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/stripe"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`stripe`, func(cfg *oauth2.Config) goth.Provider {
		return stripe.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
