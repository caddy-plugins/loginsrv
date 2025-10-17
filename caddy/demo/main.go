package main

import (
	"github.com/admpub/caddy/caddy/caddymain"
	_ "github.com/caddy-plugins/caddy-jwt/v3"
	_ "github.com/caddy-plugins/loginsrv/caddy"

	_ "github.com/caddy-plugins/loginsrv/oauth2/register/github"
	_ "github.com/caddy-plugins/loginsrv/oauth2/register/nging"
)

func main() {
	caddymain.Run()
}
