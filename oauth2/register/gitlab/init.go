package gitlab

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitlab"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`gitlab`, func(cfg *oauth2.Config) goth.Provider {
		hostURL := cfg.GetCustomisedHostURL()
		if len(hostURL) > 0 {
			return gitlab.NewCustomisedURL(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), hostURL+`/oauth/authorize`, hostURL+`/oauth/token`, hostURL+`/api/v3/user`, cfg.GetScopes()...)
		}
		return gitlab.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI())
	})
}
