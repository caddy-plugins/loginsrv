package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/stretchr/testify/assert"
)

func Test_UserInfo_Valid(t *testing.T) {
	Error(t, UserInfo{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Time{})}}.Valid())
	Error(t, UserInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
		},
	}.Valid())
	NoError(t, UserInfo{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second))}}.Valid())
}

func Test_UserInfo_AsMap(t *testing.T) {
	u := UserInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   `json:"sub"`,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Second)),
		},
		Picture:   `json:"picture,omitempty"`,
		Name:      `json:"name,omitempty"`,
		Email:     `json:"email,omitempty"`,
		Origin:    `json:"origin,omitempty"`,
		Refreshes: 42,
		Domain:    `json:"domain,omitempty"`,
		Groups:    []string{`json:"groups,omitempty"`},
	}

	givenJson, _ := json.Marshal(u.AsMap())
	given := UserInfo{}
	err := json.Unmarshal(givenJson, &given)
	NoError(t, err)
	Equal(t, u, given)
}

func Test_UserInfo_AsMap_Minimal(t *testing.T) {
	u := UserInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: `json:"sub"`,
		},
	}

	givenJson, _ := json.Marshal(u.AsMap())
	given := UserInfo{}
	err := json.Unmarshal(givenJson, &given)
	NoError(t, err)
	Equal(t, u, given)
}
