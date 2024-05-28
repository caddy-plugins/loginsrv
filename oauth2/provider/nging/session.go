package nging

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/admpub/goth"
	"golang.org/x/oauth2"
)

// Session stores data during the auth process with Github.
type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	rootURL      string // http://***/
}

func (s *Session) SetRootURL(rootURL string) {
	s.rootURL = rootURL
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Github provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Github and return the access token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	var options []oauth2.AuthCodeOption
	if !p.isFullCallbackURL {
		options = append(options, oauth2.SetAuthURLParam(`redirect_uri`, s.rootURL+strings.TrimPrefix(p.CallbackURL, `/`)))
	}
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), params.Get("code"), options...)
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.Expiry = token.Expiry
	return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
