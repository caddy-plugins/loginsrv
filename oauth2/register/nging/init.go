package nging

import (
	"strings"

	"github.com/admpub/goth"
	"github.com/caddy-plugins/loginsrv/oauth2"
	"github.com/caddy-plugins/loginsrv/oauth2/provider/nging"
)

func init() {
	oauth2.Register(`nging`, func(cfg *oauth2.Config) goth.Provider {
		hostURL := cfg.Extra[`host_url`]
		if len(hostURL) > 0 {
			hostURL = strings.TrimSuffix(hostURL, `/`)
		}
		return nging.New(cfg.ClientID, cfg.ClientSecret, cfg.GetRedirectURI(), hostURL, `profile`)
	})
}
