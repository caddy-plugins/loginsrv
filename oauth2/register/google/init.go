package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/google"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`google`, func(cfg *oauth2.Config) goth.Provider {
		return google.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI())
	})
}
