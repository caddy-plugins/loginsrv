package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitea"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`gitea`, func(cfg *oauth2.Config) goth.Provider {
		return gitea.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
