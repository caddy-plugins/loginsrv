package osiam

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddy-plugins/loginsrv/model"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/stretchr/testify/assert"
)

func TestBackend_Authenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(osiamMockHandler))
	defer server.Close()

	// positive case
	backend, err := NewBackend(server.URL, "example-client", "secret")
	NoError(t, err)
	authenticated, userInfo, err := backend.Authenticate("admin", "koala")

	NoError(t, err)
	True(t, authenticated)
	Equal(t,
		model.UserInfo{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: `admin`,
			},
			Origin: "osiam",
		},
		userInfo)

	// wrong client credentials
	backend, err = NewBackend(server.URL, "example-client", "XXX")
	NoError(t, err)
	authenticated, _, err = backend.Authenticate("admin", "koala")
	Error(t, err)
	False(t, authenticated)

	// wrong user credentials
	backend, err = NewBackend(server.URL, "example-client", "secret")
	NoError(t, err)
	authenticated, _, err = backend.Authenticate("admin", "XXX")
	NoError(t, err)
	False(t, authenticated)

}

func TestBackend_AuthenticateErrorCases(t *testing.T) {
	_, err := NewBackend("://", "example-client", "secret")
	Error(t, err)

	_, err = NewBackend("http://example.com", "", "secret")
	Error(t, err)

	_, err = NewBackend("http://example.com", "example-client", "")
	Error(t, err)
}
