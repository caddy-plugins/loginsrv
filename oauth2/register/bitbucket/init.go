package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/bitbucket"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`bitbucket`, func(cfg *oauth2.Config) goth.Provider {
		return bitbucket.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
