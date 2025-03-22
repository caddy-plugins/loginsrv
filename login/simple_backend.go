package login

import (
	"errors"

	"github.com/caddy-plugins/loginsrv/model"
	"github.com/golang-jwt/jwt/v5"
)

// SimpleProviderName const with the providers name
const SimpleProviderName = "simple"

func init() {
	RegisterProvider(
		&ProviderDescription{
			Name:     SimpleProviderName,
			HelpText: "Simple login backend opts: user1=password,user2=password,..",
		},
		SimpleBackendFactory)
}

// SimpleBackendFactory returns a new configured SimpleBackend
func SimpleBackendFactory(config map[string]string) (Backend, error) {
	userPassword := map[string]string{}
	for k, v := range config {
		userPassword[k] = v
	}
	if len(userPassword) == 0 {
		return nil, errors.New("no users provided for simple backend")
	}
	return NewSimpleBackend(userPassword), nil
}

// SimpleBackend working on a map of username password pairs
type SimpleBackend struct {
	userPassword map[string]string
}

// NewSimpleBackend creates a new SIMPLE Backend and verifies the parameters.
func NewSimpleBackend(userPassword map[string]string) *SimpleBackend {
	return &SimpleBackend{
		userPassword: userPassword,
	}
}

// Authenticate the user
func (sb *SimpleBackend) Authenticate(username, password string) (bool, model.UserInfo, error) {
	if p, exist := sb.userPassword[username]; exist && p == password {
		return true, model.UserInfo{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: username,
			},
			Origin: SimpleProviderName,
		}, nil
	}
	return false, model.UserInfo{}, nil
}
