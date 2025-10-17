package salesforce

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/salesforce"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`salesforce`, func(cfg *oauth2.Config) goth.Provider {
		return salesforce.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
