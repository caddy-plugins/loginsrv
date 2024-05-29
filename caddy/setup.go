package caddy

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/caddy"
	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/admpub/goth"
	"github.com/admpub/pp"
	"github.com/caddy-plugins/loginsrv/logging"
	"github.com/caddy-plugins/loginsrv/login"
	"github.com/webx-top/com"

	// Import all backends, packaged with the caddy plugin
	_ "github.com/caddy-plugins/loginsrv/htpasswd"
	_ "github.com/caddy-plugins/loginsrv/httpupstream"
	_ "github.com/caddy-plugins/loginsrv/oauth2"
	_ "github.com/caddy-plugins/loginsrv/osiam"
)

func init() {
	caddy.RegisterPlugin("login", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func dump(a ...interface{}) {
	pp.Println(a...)
}

func getSiteURL(c *caddy.Controller) string {
	var siteURL string
	if len(c.ServerBlockKeys) > 0 {
		addrs := make([]string, len(c.ServerBlockKeys))
		for index, addr := range c.ServerBlockKeys {
			parsedValue := com.ParseEnvVar(addr)
			currentHost := `127.0.0.1`
			if len(parsedValue) > 0 {
				switch {
				case parsedValue[0] == ':':
					addr = `http://` + currentHost + parsedValue
				case strings.HasPrefix(parsedValue, `0.0.0.0:`):
					addr = `http://` + currentHost + strings.TrimPrefix(parsedValue, `0.0.0.0`)
				case !strings.Contains(parsedValue, `//`):
					addr = `http://` + parsedValue
				default:
					addr = strings.ReplaceAll(parsedValue, `*`, `test`)
				}
			}
			if strings.HasPrefix(addr, `https://`) {
				siteURL = addr
				break
			}
			addrs[index] = addr
		}
		if len(siteURL) == 0 {
			siteURL = addrs[0]
		}
	}
	return siteURL
}

// setup configures a new loginsrv instance.
func setup(c *caddy.Controller) error {
	logging.Set("info", true)
	goth.ClearProviders()
	//dump(c)
	for c.Next() {
		args := c.RemainingArgs()

		config, err := parseConfig(c)
		if err != nil {
			return err
		}

		if config.Template != "" && !filepath.IsAbs(config.Template) {
			config.Template = filepath.Join(httpserver.GetConfig(c).Root, config.Template)
		}

		if len(args) == 1 {
			logging.Logger.Warnf("DEPRECATED: Please set the login path by parameter login_path and not as directive argument (%v:%v)", c.File(), c.Line())
			config.LoginPath = path.Join(args[0], "/login")
		}

		// if len(config.SiteURL) == 0 {
		// 	config.SiteURL = siteURL
		// }
		loginHandler, err := login.NewHandler(config)
		if err != nil {
			return err
		}

		httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
			return NewCaddyHandler(next, loginHandler, config)
		})
	}

	return nil
}

func parseConfig(c *caddy.Controller) (*login.Config, error) {
	cfg := login.DefaultConfig()
	cfg.Host = ""
	cfg.Port = ""
	cfg.LogLevel = ""

	fs := flag.NewFlagSet("loginsrv-config", flag.ContinueOnError)
	cfg.ConfigureFlagSet(fs)

	secretProvidedByConfig := false
	for c.NextBlock() {
		// caddy prefers '_' in parameter names,
		// so we map them to the '-' from the command line flags
		// the replacement supports both, for backwards compatibility
		name := strings.Replace(c.Val(), "_", "-", -1)
		args := c.RemainingArgs()
		if len(args) != 1 {
			return cfg, fmt.Errorf("Wrong number of arguments for %v: %v (%v:%v)", name, args, c.File(), c.Line())
		}
		value := args[0]

		f := fs.Lookup(name)
		if f == nil {
			return cfg, fmt.Errorf("Unknown parameter for login directive: %v (%v:%v)", name, c.File(), c.Line())
		}
		err := f.Value.Set(value)
		if err != nil {
			return cfg, fmt.Errorf("Invalid value for parameter %v: %v (%v:%v)", name, value, c.File(), c.Line())
		}

		if name == "jwt-secret" {
			secretProvidedByConfig = true
		}
	}

	if err := cfg.ResolveFileReferences(); err != nil {
		return nil, err
	}

	secretFromEnv, secretFromEnvWasSetBefore := os.LookupEnv("JWT_SECRET")
	if !secretProvidedByConfig && secretFromEnvWasSetBefore {
		cfg.JwtSecret = secretFromEnv
	}
	if !secretFromEnvWasSetBefore {
		// populate the secret to caddy.jwt,
		// but do not change a environment variable, which somebody has set it.
		os.Setenv("JWT_SECRET", cfg.JwtSecret)
	}

	return cfg, nil
}
