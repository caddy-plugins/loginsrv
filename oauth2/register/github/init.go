package github

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/github"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`github`, func(cfg *oauth2.Config) goth.Provider {
		return github.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
