package config

import "github.com/melvinodsa/go-iam/sdk"

type Jwt struct {
	secret sdk.MaskedBytes
}

func (j Jwt) Secret() sdk.MaskedBytes {
	return j.secret
}
