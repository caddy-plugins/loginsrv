package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/admpub/goth"
	"github.com/admpub/goth/gothic"
	"github.com/admpub/goth/providers/faux"
	"github.com/gorilla/sessions"
	. "github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var testConfig = &Config{
	ClientID:     "client42",
	ClientSecret: "secret",
	RedirectURI:  "http://localhost/callback",
	Scope:        "email other",
	ProviderName: `faux`,
}

type testProvider struct {
	*faux.Provider
}

// BeginAuth is used only for testing.
func (p *testProvider) BeginAuth(state string) (goth.Session, error) {
	c := &oauth2.Config{
		ClientID:     testConfig.ClientID,
		ClientSecret: testConfig.ClientSecret,
		RedirectURL:  testConfig.RedirectURI,
		Scopes:       []string{`email`, `other`},
		Endpoint: oauth2.Endpoint{
			AuthURL: "http://example.com/auth",
		},
	}
	url := c.AuthCodeURL(state)
	return &faux.Session{
		ID:          "id",
		Name:        `test`,
		Email:       ``,
		AuthURL:     url,
		AccessToken: `e72e16c7e42f292c6912e7710c838347ae178b4a`,
	}, nil
}

func setProvider() {
	p := &testProvider{}
	goth.UseProviders(p)
	testConfig.SetProvider(p)
}

var cookieStore *sessions.CookieStore

func init() {
	sessionSecret := `1234567890123456`
	os.Setenv(`SESSION_SECRET`, sessionSecret)

	cookieStore = sessions.NewCookieStore([]byte(sessionSecret))
	cookieStore.Options.HttpOnly = true
	gothic.Store = cookieStore
}

func Test_StartFlow(t *testing.T) {
	setProvider()
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(`GET`, `http://localhost/`, nil)
	StartFlow(testConfig, resp, req)

	Contains(t, resp.Body.String(), `<a href=`)
	Equal(t, http.StatusTemporaryRedirect, resp.Code)

	// assert that we received a state cookie
	cHeader := strings.Split(resp.Header().Get("Set-Cookie"), ";")[0]
	Equal(t, gothic.SessionName, strings.Split(cHeader, "=")[0])
	//state := strings.Split(cHeader, "=")[1]
	// state, err := gothic.GetFromSession(testConfig.ProviderName, req)
	// NoError(t, err)
	// t.Logf(`==============> %+v`, cHeader)

	expectedLocation := fmt.Sprintf("%v?client_id=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v",
		`http://example.com/auth`,
		testConfig.ClientID,
		url.QueryEscape(testConfig.RedirectURI),
		"email+other",
		``,
	)

	Contains(t, resp.Header().Get("Location"), expectedLocation)
}

/*
func Test_Authenticate(t *testing.T) {
	setProvider()
	// mock a server for token exchange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Equal(t, "POST", r.Method)
		Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		Equal(t, "application/json", r.Header.Get("Accept"))

		body, _ := io.ReadAll(r.Body)
		Equal(t, "client_id=client42&client_secret=secret&code=theCode&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"e72e16c7e42f292c6912e7710c838347ae178b4a", "name":"test", "scope":"repo gist", "token_type":"bearer"}`))
	}))
	defer server.Close()

	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	user, err := Authenticate(testConfig, resp, request)
	t.Logf(`%+v`, user)
	NoError(t, err)
	Equal(t, "test", user.Sub)
}

func Test_Authenticate_CodeExchangeError(t *testing.T) {
	setProvider()
	var testReturnCode int
	testResponseJSON := `{"error":"bad_verification_code","error_description":"The code passed is incorrect or expired.","error_uri":"https://developer.github.com/v3/oauth/#bad-verification-code"}`
	// mock a server for token exchange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(testReturnCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testResponseJSON))
	}))
	defer server.Close()

	testConfigCopy := testConfig

	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	testReturnCode = 500
	user, err := Authenticate(testConfigCopy, resp, request)
	Error(t, err)
	EqualError(t, err, "error: expected http status 200 on token exchange, but got 500")
	Equal(t, "", user.Sub)

	testReturnCode = 200
	user, err = Authenticate(testConfigCopy, resp, request)
	Error(t, err)
	EqualError(t, err, `error: got "bad_verification_code" on token exchange`)
	Equal(t, "", user.Sub)

	testReturnCode = 200
	testResponseJSON = `{"foo": "bar"}`
	user, err = Authenticate(testConfigCopy, resp, request)
	Error(t, err)
	EqualError(t, err, `error: no access_token on token exchange`)
	Equal(t, "", user.Sub)

}

func Test_Authentication_ProviderError(t *testing.T) {
	setProvider()
	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.URL, _ = url.Parse("http://localhost/callback?error=provider_login_error")

	_, err := Authenticate(testConfig, resp, request)

	Error(t, err)
	Equal(t, "error: provider_login_error", err.Error())
}

func Test_Authentication_StateError(t *testing.T) {
	setProvider()
	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=XXXXXXX")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	_, err := Authenticate(testConfig, resp, request)

	Error(t, err)
	Equal(t, "error: oauth state param could not be verified", err.Error())
}

func Test_Authentication_NoCodeError(t *testing.T) {
	setProvider()
	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?state=theState")

	_, err := Authenticate(testConfig, resp, request)

	Error(t, err)
	Equal(t, "error: no auth code provided", err.Error())
}

func Test_Authentication_Provider500(t *testing.T) {
	setProvider()
	// mock a server for token exchange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	testConfigCopy := testConfig

	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	_, err := Authenticate(testConfigCopy, resp, request)

	Error(t, err)
	Equal(t, "error: expected http status 200 on token exchange, but got 500", err.Error())
}

func Test_Authentication_ProviderNetworkError(t *testing.T) {
	setProvider()
	testConfigCopy := testConfig

	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	_, err := Authenticate(testConfigCopy, resp, request)

	Error(t, err)
	Contains(t, err.Error(), "invalid port")
}

func Test_Authentication_TokenParseError(t *testing.T) {
	setProvider()
	// mock a server for token exchange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_t`))

	}))
	defer server.Close()

	testConfigCopy := testConfig

	resp := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", testConfig.RedirectURI, nil)
	request.Header.Set("Cookie", "oauthState=theState")
	request.URL, _ = url.Parse("http://localhost/callback?code=theCode&state=theState")

	_, err := Authenticate(testConfigCopy, resp, request)

	Error(t, err)
	Equal(t, "error on parsing oauth token: unexpected end of JSON input", err.Error())
}
*/
