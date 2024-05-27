package main

import (
	_ "github.com/BTBurke/caddy-jwt"
	_ "github.com/caddy-plugins/loginsrv/caddy"
	"github.com/caddyserver/caddy/caddy/caddymain"
)

func main() {
	caddymain.Run()
}
