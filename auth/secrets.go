package auth

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type secretDatabase struct {
	Cognito struct {
		ClientId     string `yaml:"clientId" env:"CLIENT_ID" env-description:"" env-required:""`
		ClientSecret string `yaml:"clientSecret" env:"CLIENT_SECRET" env-description:"" env-required:""`
	} `yaml:"cognito"`
}

var secrets secretDatabase

func Init(path string) error {
	err := cleanenv.ReadConfig(path, &secrets)
	if err != nil {
		return err
	}
	return nil
}

func get(propertyName string) string {
	switch propertyName {
	case "clientId":
		return secrets.Cognito.ClientId
	case "clientSecret":
		return secrets.Cognito.ClientSecret
	default:
		return ""
	}
}
