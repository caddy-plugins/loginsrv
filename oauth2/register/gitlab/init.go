package gitlab

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitlab"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`gitlab`, func(cfg *oauth2.Config) goth.Provider {
		return gitlab.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI())
	})
}
