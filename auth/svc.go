package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	validator "github.com/eryk-vieira/go-cognito-jwt-validator"
	"github.com/jrhoward/cognito-authflow-handler/config"
)

type Tokens struct {
	Id_token      string `json:"id_token"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Expire        int    `json:"expire"`
	Token_type    string `json:"token_type"`
}

func setCognitoToken(value string, refresh bool) (Tokens, error) {
	oauthServer := config.GetOauthTokenEndpoint()
	body := url.Values{}
	if refresh {
		body.Set("refresh_token", value)
		body.Set("grant_type", "refresh_token")
	} else {
		body.Set("code", value)
		body.Set("grant_type", "authorization_code")
		body.Set("redirect_uri", config.Get("callBackUrl"))
		body.Set("scope", config.Get("scope"))
	}

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, oauthServer, strings.NewReader(body.Encode()))
	clientAuth := base64.StdEncoding.EncodeToString([]byte(get("clientId") + ":" + get("clientSecret")))
	r.Header.Add("Authorization", "Basic "+clientAuth)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	var t Tokens
	resp, err := client.Do(r)
	if err != nil {
		return t, err
	}
	if resp.StatusCode != 200 {
		return t, errors.New(resp.Status)
	}
	json.NewDecoder(resp.Body).Decode(&t)
	return t, nil
}

func validateCognitoToken(tokenString string) error {
	validator := validator.New(&validator.Config{
		Region:          config.Get("awsRegion"),
		CognitoPoolId:   config.Get("poolId"),
		CognitoClientId: get("clientId"),
	})
	err := validator.Validate(tokenString)

	if err != nil {
		return err
	}
	return nil
}

func revokeToken(refreshToken string) error {
	oauthServer := config.GetOauthRevokeEndpoint()
	body := url.Values{}
	body.Set("token", refreshToken)
	body.Set("client_id", get("clientId"))

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, oauthServer, strings.NewReader(body.Encode()))
	clientAuth := base64.StdEncoding.EncodeToString([]byte(get("clientId") + ":" + get("clientSecret")))
	r.Header.Add("Authorization", "Basic "+clientAuth)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}
