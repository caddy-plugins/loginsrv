package caddy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/caddy-plugins/loginsrv/login"
)

// CaddyHandler is the loginsrv handler wrapper for caddy
type CaddyHandler struct {
	next         httpserver.Handler
	config       *login.Config
	loginHandler *login.Handler
}

// NewCaddyHandler create the handler
func NewCaddyHandler(next httpserver.Handler, loginHandler *login.Handler, config *login.Config) *CaddyHandler {
	return &CaddyHandler{
		next:         next,
		config:       config,
		loginHandler: loginHandler,
	}
}

func (h *CaddyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	isLoginPath := strings.HasPrefix(r.URL.Path, h.config.LoginPath)
	// Fetch jwt token. If valid set a Caddy replacer for {user}
	userInfo, valid := h.loginHandler.GetToken(r)
	if !valid && !isLoginPath && strings.HasPrefix(r.UserAgent(), "docker/") { // Auth gate: Docker requests require JWT
		return h.dockerAuthResponse(w, r)
	}
	if valid {
		// let upstream middleware (e.g. fastcgi and cgi) know about authenticated
		// user; this replaces the request with a wrapped instance
		r = r.WithContext(context.WithValue(r.Context(),
			httpserver.RemoteUserCtxKey, userInfo.Subject))

		// Provide username to be used in log by replacer
		repl := httpserver.NewReplacer(r, nil, "-")
		repl.Set("user", userInfo.Subject)
	}

	if isLoginPath {
		h.loginHandler.ServeHTTP(w, r)
		return 0, nil
	}

	return h.next.ServeHTTP(w, r)
}

func (h *CaddyHandler) dockerAuthResponse(w http.ResponseWriter, r *http.Request) (int, error) {
	// Advertise the token endpoint to Docker clients via Www-Authenticate header
	realm := h.config.LoginPath
	service := r.URL.Query().Get("service")
	if len(service) == 0 {
		service = r.Host
	}
	scope := r.URL.Query().Get("scope")
	authHeader := fmt.Sprintf(`Bearer realm=%q,service=%q`, realm, service)
	if len(scope) > 0 {
		authHeader += fmt.Sprintf(`,scope=%q`, scope)
	}
	w.Header().Set("Www-Authenticate", authHeader)
	return http.StatusUnauthorized, nil
}
