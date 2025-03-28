package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/caddy-plugins/loginsrv/model"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type userClaimsProvider struct {
	url        string
	auth       string
	httpClient http.Client
}

func newUserClaimsProvider(url, auth string, timeout time.Duration) (*userClaimsProvider, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}

	return &userClaimsProvider{
		url:        url,
		auth:       auth,
		httpClient: http.Client{Timeout: timeout},
	}, nil
}

func (provider *userClaimsProvider) Claims(userInfo model.UserInfo) (jwt.Claims, error) {
	claimsURL := provider.buildURL(userInfo)
	req, err := http.NewRequest(http.MethodGet, claimsURL, nil)
	if err != nil {
		return nil, fmt.Errorf(`failed to request the user claims API(%s): %w`, provider.url, err)
	}
	if provider.auth != "" {
		req.Header.Add("Authorization", "Bearer "+provider.auth)
	}

	resp, err := provider.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(`failed to request the user claims API(%s): %w`, provider.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return userInfo, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("bad http response code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	remoteClaims := map[string]interface{}{}
	err = decoder.Decode(&remoteClaims)
	if err != nil {
		return nil, fmt.Errorf(`failed to decode the user claims(%s): %w`, provider.url, err)
	}

	return jwt.MapClaims(mergeClaims(userInfo, remoteClaims)), nil
}

func (provider *userClaimsProvider) buildURL(userInfo model.UserInfo) string {
	// error can be ignored, it was already checked in validateURL
	u, _ := url.Parse(provider.url)

	query := u.Query()

	query.Add("sub", url.QueryEscape(userInfo.Subject))
	if userInfo.ID != "" {
		query.Add("id", url.QueryEscape(userInfo.ID))
	}
	if userInfo.Origin != "" {
		query.Add("origin", url.QueryEscape(userInfo.Origin))
	}
	if userInfo.Domain != "" {
		query.Add("domain", url.QueryEscape(userInfo.Domain))
	}
	if userInfo.Email != "" {
		query.Add("email", url.QueryEscape(userInfo.Email))
	}
	if len(userInfo.Groups) > 0 {
		for _, group := range userInfo.Groups {
			query.Add("group", url.QueryEscape(group))
		}
	}

	u.RawQuery = query.Encode()

	return u.String()
}

func mergeClaims(userInfo model.UserInfo, remoteClaims map[string]interface{}) customClaims {
	claims := customClaims(userInfo.AsMap())
	claims.merge(remoteClaims)
	return claims
}

func validateURL(s string) error {
	_, err := url.Parse(s)
	return errors.Wrap(err, "invalid claims provider url")
}
