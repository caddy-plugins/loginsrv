package login

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/caddy-plugins/loginsrv/logging"
	"github.com/caddy-plugins/loginsrv/model"
)

//go:embed bootstrap.min.css
var bootstrapCSS string

//go:embed bootstrap-social.min.css
var bsSocialCSS string

//go:embed font-awesome.css
var fontAwesomeCSS string

const partials = `

{{- define "styles" -}}
<!-- <link uic-remove rel="stylesheet" href="https://cdn.bootcdn.net/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"> -->
<link uic-remove rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
<style>
{{- .Styles.CSS -}}
.vertical-offset-100{padding-top:100px;}
.login-or-container {text-align: center;margin: 0;margin-bottom: 10px;clear: both;color: #6a737c;font-variant: small-caps;}
.login-or-hr {margin-bottom: 0;position: relative;top: 28px;height: 0;border: 0;border-top: 1px solid #e4e6e8;}
.login-or {display: inline-block;position: relative;padding: 10px;background-color: #FFF;}
.login-picture {height: 120px;border-radius: 3px;margin-bottom: 10px;box-shadow: 0 0 5px #ccc;}
@media (prefers-color-scheme: dark) {
body,.login-or{background-color:#111;}.login-or-hr{border-color:#444}body{color:#aaa}
.panel-default{border-color:#333}.panel{background-color:#333}..login-or-container{color:#888}
.panel-default>.panel-heading{color: #fff;background-color: #555;border-color: #333;}
.form-control{background-color: #000;border-color:#333;color:#f1f1f1}
}</style>
{{- end -}}

{{- define "userInfo" -}}
              {{- with .UserInfo -}}
                <h1>Welcome {{.Subject}}!</h1>
                <br/>
                {{if .Picture}}<img class="login-picture" src="{{.Picture}}?s=120">{{end -}}
                {{if .Name}}<h3>{{.Name}}</h3>{{end -}}
              {{- end -}}
              <br/>
              <a class="btn btn-md btn-primary" href="{{ .Config.LoginPath }}?logout=true">Logout</a>
{{- end -}}

{{- define "login" -}}
              {{- range $providerName, $opts := .Config.Oauth  -}}
                <a class="btn btn-block btn-lg btn-social btn-primary btn-{{ $providerName }}" href="{{ trimRight $.Config.LoginPath "/" }}/{{ $providerName }}">
                  <span class="fa fa-user fa-{{ $providerName }}"></span> Sign in with {{ $providerName | ucfirst }}
                </a>
              {{- end -}}

              {{- if and (not (eq (len .Config.Backends) 0)) (not (eq (len .Config.Oauth) 0)) -}}
                <div class="login-or-container">
                  <hr class="login-or-hr">
                  <div class="login-or lead">or</div>
                </div>
              {{- end -}}

              {{- if not (eq (len .Config.Backends) 0)  -}}
                <div class="panel panel-default">
  	          <div class="panel-heading">  
  		    <div class="panel-title">
  		      <h4>Sign in</h4>
		    </div>
	          </div>
	          <div class="panel-body">
			  {{- if .Failure}}<div class="alert alert-warning" role="alert">Invalid credentials</div>{{- end -}} 
		    <form accept-charset="UTF-8" role="form" method="POST" action="{{.Config.LoginPath}}">
                <fieldset>
		        <div class="form-group">
		          <input class="form-control" placeholder="Username" name="username" value="{{.UserInfo.Subject}}" type="text">
		        </div>
		        <div class="form-group">
		          <input class="form-control" placeholder="Password" name="password" type="password" value="">
		        </div>
		        <input class="btn btn-lg btn-success btn-block" type="submit" value="Login">
		      </fieldset>
		    </form>
	          </div>
	        </div>
              {{- end -}}
{{- end}}`

var layout = `<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    {{- template "styles" . -}}
  </head>
  <body>
    <uic-fragment name="content">
      <div class="container">
        <div class="row vertical-offset-100">
    	  <div class="col-md-4 col-md-offset-4">

            {{- if .Error -}}
              <div class="alert alert-danger" role="alert">
			  	{{- if .Reason -}}
				<strong>Error:</strong> {{.Reason}}
				{{- else -}}
                <strong>Internal Error. </strong> Please try again later.
				{{- end -}}
              </div>
            {{end -}}

            {{- if .Authenticated -}}

              {{template "userInfo" .  -}}

            {{- else -}}

              {{- template "login" .  -}}

            {{- end -}}
	  </div>
	</div>
      </div>
    </uic-fragment>
  </body>
</html>`

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
