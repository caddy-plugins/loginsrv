package gitea

import (
	"github.com/admpub/goth"
	"github.com/admpub/goth/providers/gitea"
	"github.com/caddy-plugins/loginsrv/oauth2"
)

func init() {
	oauth2.Register(`gitea`, func(cfg *oauth2.Config) goth.Provider {
		hostURL := cfg.GetCustomisedHostURL()
		if len(hostURL) > 0 {
			return gitea.NewCustomisedURL(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), hostURL+`/login/oauth/authorize`, hostURL+`/login/oauth/access_token`, hostURL+`/api/v1/user`, cfg.GetScopes()...)
		}
		return gitea.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), cfg.GetScopes()...)
	})
}
