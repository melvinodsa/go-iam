package config

import "github.com/melvinodsa/go-iam/sdk"

type Redis struct {
	Host     string          `json:"host"`
	Password sdk.MaskedBytes `json:"password"`
	DB       int             `json:"db"`
}
