package login

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/caddy-plugins/loginsrv/model"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/stretchr/testify/assert"
)

func clearCSS(r string) string {
	parts := strings.SplitN(r, `</style>`, 2)
	if len(parts) != 2 {
		return r
	}
	return parts[1]
}

func Test_form(t *testing.T) {
	// show error
	recorder := httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Error: true,
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
			Oauth:     Options{
				/*
					`github`: map[string]string{
						`client_id`:     `client_id`,
						`client_secret`: `client_secret`,
						`redirect_uri`:  `/login/github`,
					},
				*/
			},
		},
	})
	result := clearCSS(recorder.Body.String())
	Contains(t, result, `<form`)
	NotContains(t, result, `github`)
	NotContains(t, result, `Welcome`)
	Contains(t, result, `Error`)

	// only form
	recorder = httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
		},
	})
	result = clearCSS(recorder.Body.String())
	Contains(t, result, `<form`)
	NotContains(t, result, `github`)
	NotContains(t, result, `Welcome`)
	NotContains(t, result, `Error`)

	// only links
	recorder = httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Config: &Config{
			LoginPath: "/login",
			Oauth:     Options{"github": {}},
		},
	})
	result = clearCSS(recorder.Body.String())
	NotContains(t, result, `<form`)
	Contains(t, result, `href="/login/github"`)
	NotContains(t, result, `Welcome`)
	NotContains(t, result, `Error`)

	// with form and links
	recorder = httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
			Oauth:     Options{"github": {}},
		},
	})
	result = clearCSS(recorder.Body.String())
	Contains(t, result, `<form`)
	Contains(t, result, `href="/login/github"`)
	NotContains(t, result, `Welcome`)
	NotContains(t, result, `Error`)

	// show only the user info
	recorder = httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Authenticated: true,
		UserInfo:      model.UserInfo{RegisteredClaims: jwt.RegisteredClaims{Subject: "smancke"}, Name: "Sebastian Mancke"},
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
			Oauth:     Options{"github": {}},
		},
	})
	result = clearCSS(recorder.Body.String())
	NotContains(t, result, `<form`)
	NotContains(t, result, `href="/login/github"`)
	Contains(t, result, `Welcome smancke`)
	NotContains(t, result, `Error`)
}

func Test_form_executeError(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{})
	Equal(t, 500, recorder.Code)
}

func Test_form_customTemplate(t *testing.T) {
	f, err := os.CreateTemp("", "")
	NoError(t, err)
	f.WriteString(`<html><body>My custom template {{template "login" .}}</body></html>`)
	f.Close()
	defer os.Remove(f.Name())

	recorder := httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Error: true,
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
			Template:  f.Name(),
		},
	})
	result := clearCSS(recorder.Body.String())
	Contains(t, result, `My custom template`)
	Contains(t, result, `<form`)
	NotContains(t, result, `github`)
	NotContains(t, result, `Welcome`)
	NotContains(t, result, `Error`)
	NotContains(t, result, `style`)
}

func Test_form_customTemplate_ParseError(t *testing.T) {
	f, err := os.CreateTemp("", "")
	NoError(t, err)
	f.WriteString(`<html><body>My custom template {{template "login" `)
	f.Close()
	defer os.Remove(f.Name())

	recorder := httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Config: &Config{
			LoginPath: "/login",
			Backends:  Options{"simple": {}},
			Template:  f.Name(),
		},
	})
	Equal(t, 500, recorder.Code)
}

func Test_form_customTemplate_MissingFile(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeLoginForm(recorder, loginFormData{
		Config: &Config{
			Template: "/this/file/does/not/exist",
		},
	})
	Equal(t, 500, recorder.Code)
}

func Test_ucfirst(t *testing.T) {
	Equal(t, "", ucfirst(""))
	Equal(t, "A", ucfirst("a"))
	Equal(t, "Abc def", ucfirst("abc def"))
}
