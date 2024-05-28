package main

import (
	"strings"

	"github.com/admpub/caddy/caddy/caddymain"
	_ "github.com/caddy-plugins/caddy-jwt/v3"
	_ "github.com/caddy-plugins/loginsrv/caddy"

	"github.com/admpub/goth"
	"github.com/caddy-plugins/loginsrv/oauth2"
	"github.com/caddy-plugins/loginsrv/oauth2/provider/nging"
)

func main() {
	caddymain.Run()
}

func init() {
	oauth2.Register(`nging`, func(cfg *oauth2.Config) goth.Provider {
		hostURL := cfg.Extra[`host_url`]
		if len(hostURL) > 0 {
			hostURL = strings.TrimSuffix(hostURL, `/`)
		}
		return nging.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, hostURL, `profile`)
	})
}
