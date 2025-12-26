package login

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/caddy-plugins/loginsrv/logging"
	"github.com/caddy-plugins/loginsrv/model"
)

type loginFormData struct {
	Error         bool
	Reason        string
	Failure       bool
	Config        *Config
	Authenticated bool
	UserInfo      model.UserInfo
	CSSFiles      Styles
}

type Styles struct {
	Bootstrap   string
	BSSocial    string
	FontAwesome string
	Custom      string
}

func getEnvDefault(k string, dft string) string {
	v := os.Getenv(k)
	if len(v) == 0 {
		v = dft
	}
	return v
}

const queryVarName = `--login-assets--`

var cssFiles = Styles{
	Bootstrap:   getEnvDefault(`CSS_BOOTSTRAP_URL`, `?`+queryVarName+`=bootstrap.min.css`),
	BSSocial:    getEnvDefault(`CSS_BOOTSTRAP_SOCIAL__URL`, `?`+queryVarName+`=bootstrap-social.min.css`),
	FontAwesome: getEnvDefault(`CSS_FONT_AWESOME_URL`, `?`+queryVarName+`=font-awesome.min.css`),
	Custom:      getEnvDefault(`CSS_CUSTOM_URL`, `?`+queryVarName+`=custom.min.css`),
}

func writeLoginForm(w http.ResponseWriter, params loginFormData) {
	params.CSSFiles = cssFiles
	funcMap := template.FuncMap{
		"ucfirst":   ucfirst,
		"trimRight": strings.TrimRight,
		"toCSS":     toCSS,
	}
	templateName := "loginForm"
	if params.Config != nil && params.Config.Template != "" {
		templateName = params.Config.Template
	}
	t := template.New(templateName).Funcs(funcMap)
	t = template.Must(t.Parse(partials))
	if params.Config != nil && params.Config.Template != "" {
		customTemplate, err := os.ReadFile(params.Config.Template)
		if err != nil {
			logging.Logger.WithError(err).Error()
			w.WriteHeader(500)
			w.Write([]byte(`Internal Server Error`))
			return
		}

		t, err = t.Parse(string(customTemplate))
		if err != nil {
			logging.Logger.WithError(err).Error()
			w.WriteHeader(500)
			w.Write([]byte(`Internal Server Error`))
			return
		}
	} else {
		t = template.Must(t.Parse(layout))
	}

	b := bytes.NewBuffer(nil)
	err := t.Execute(b, params)
	if err != nil {
		logging.Logger.WithError(err).Error()
		w.WriteHeader(500)
		w.Write([]byte(`Internal Server Error`))
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", contentTypeHTML)
	if params.Error {
		w.WriteHeader(500)
	}

	w.Write(b.Bytes())
}

func ucfirst(in string) string {
	if in == "" {
		return ""
	}

	return strings.ToUpper(in[0:1]) + in[1:]
}

func toCSS(v string) template.CSS {
	return template.CSS(v)
}
