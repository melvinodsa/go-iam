package config

import "github.com/melvinodsa/go-iam/sdk"

type Encrypter struct {
	key sdk.MaskedBytes
}

func (e Encrypter) Key() sdk.MaskedBytes {
	return e.key
}
