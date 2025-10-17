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
	Styles        Styles
}

type Styles struct {
	Bootstrap   string
	BSSocial    string
	FontAwesome string
}

func (s Styles) String() string {
	return s.Bootstrap + s.BSSocial + s.FontAwesome
}

func (s Styles) CSS() template.CSS {
	return template.CSS(s.String())
}

func writeLoginForm(w http.ResponseWriter, params loginFormData) {
	params.Styles = Styles{
		Bootstrap:   bootstrapCSS,
		BSSocial:    bsSocialCSS,
		FontAwesome: fontAwesomeCSS,
	}
	funcMap := template.FuncMap{
		"ucfirst":   ucfirst,
		"trimRight": strings.TrimRight,
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
