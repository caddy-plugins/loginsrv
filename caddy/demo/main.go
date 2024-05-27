package main

import (
	"github.com/admpub/caddy/caddy/caddymain"
	_ "github.com/caddy-plugins/caddy-jwt/v3"
	_ "github.com/caddy-plugins/loginsrv/caddy"
)

func main() {
	caddymain.Run()
}
