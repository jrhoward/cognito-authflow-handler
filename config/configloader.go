package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type configDatabase struct {
	Server struct {
		Port                  string `yaml:"port" env-description:"port this server is listening on" env-required:""`
		Domain                string `yaml:"domain" env-description:"" env-required:""`
		AuthHandlerRedirect   string `yaml:"authHandlerRedirect" env-description:"" env-required:""`
		LogoutHandlerRedirect string `yaml:"logoutHandlerRedirect" env-description:"" env-required:""`
		CookieMaxAge          int    `yaml:"cookieMaxAge" env-description:"" env-required:""`
		IdCookieName          string `yaml:"idCookieName" env-description:"" env-required:""`
		RefreshCookieName     string `yaml:"refreshCookieName" env-description:"" env-required:""`
	} `yaml:"server"`
	Cognito struct {
		OauthServer string `yaml:"oauthServer" env-description:"" env-required:""`
		CallBackUrl string `yaml:"callBackUrl" env-description:"" env-required:""`
		PoolId      string `yaml:"poolId" env-description:"" env-required:""`
		Scope       string `yaml:"scope" env-description:"" env-required:""`
		AwsRegion   string `yaml:"awsRegion" env-description:"" env-required:""`
	} `yaml:"cognito"`
}

var config configDatabase
var serverHost string
var oauthTokenEndpoint string
var oauthRevokeEndpoint string

func Init(path string) error {
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		return err
	}
	setServerHost()
	setOauthTokenEndpoint()
	setOauthRevokeEndpoint()
	return nil
}

func setServerHost() {
	serverHost = config.Server.Domain + ":" + config.Server.Port
}

func GetServerHost() string {
	return serverHost
}

func setOauthRevokeEndpoint() {
	oauthRevokeEndpoint = config.Cognito.OauthServer + "/oauth2/revoke"
}

func GetOauthRevokeEndpoint() string {
	return oauthRevokeEndpoint
}

func setOauthTokenEndpoint() {
	oauthTokenEndpoint = config.Cognito.OauthServer + "/oauth2/token"
}

func GetOauthTokenEndpoint() string {
	return oauthTokenEndpoint
}

func GetCookieMaxAge() int {
	return config.Server.CookieMaxAge
}

func Get(propertyName string) string {
	switch propertyName {
	case "domain":
		return config.Server.Domain
	case "port":
		return config.Server.Port
	case "authHandlerRedirect":
		return config.Server.AuthHandlerRedirect
	case "logoutHandlerRedirect":
		return config.Server.LogoutHandlerRedirect
	case "oauthServer":
		return config.Cognito.OauthServer
	case "callBackUrl":
		return config.Cognito.CallBackUrl
	case "poolId":
		return config.Cognito.PoolId
	case "scope":
		return config.Cognito.Scope
	case "awsRegion":
		return config.Cognito.AwsRegion
	case "idCookieName":
		return config.Server.IdCookieName
	case "refreshCookieName":
		return config.Server.RefreshCookieName
	default:
		return ""

	}
}
