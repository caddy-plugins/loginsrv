package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/yahoo"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`yahoo`, func(cfg *oauth2.Config) goth.Provider {
		return yahoo.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
