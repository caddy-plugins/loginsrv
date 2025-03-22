package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserInfo holds the parameters returned by the backends.
// This information will be serialized to build the JWT token contents.
type UserInfo struct {
	jwt.RegisteredClaims
	Picture   string   `json:"picture,omitempty"`
	Name      string   `json:"name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Origin    string   `json:"origin,omitempty"`
	Refreshes int      `json:"refs,omitempty"`
	Domain    string   `json:"domain,omitempty"`
	Groups    []string `json:"groups,omitempty"`
}

// Valid lets us use the user info as Claim for jwt-go.
// It checks the token expiry.
func (u UserInfo) Valid() error {
	if u.ExpiresAt != nil && u.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}
	return nil
}

func (u UserInfo) AsMap() map[string]interface{} {
	m := map[string]interface{}{
		"sub": u.Subject,
	}
	if u.ID != "" {
		m["id"] = u.ID
	}
	if u.Picture != "" {
		m["picture"] = u.Picture
	}
	if u.Name != "" {
		m["name"] = u.Name
	}
	if u.Email != "" {
		m["email"] = u.Email
	}
	if u.Origin != "" {
		m["origin"] = u.Origin
	}
	if u.ExpiresAt != nil {
		m["exp"] = u.ExpiresAt.Unix()
	}
	if u.Refreshes != 0 {
		m["refs"] = u.Refreshes
	}
	if u.Domain != "" {
		m["domain"] = u.Domain
	}
	if len(u.Groups) > 0 {
		m["groups"] = u.Groups
	}
	return m
}
