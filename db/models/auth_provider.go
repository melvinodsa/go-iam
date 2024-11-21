package models

type AuthProvider struct {
	Id        string `bson:"id"`
	Provider  string `bson:"provider"`
	IsEnabled bool   `bson:"is_enabled"`
}

type AuthProviderModel struct {
	iam
	IdKey        string
	ProviderKey  string
	IsEnabledKey string
}

func (a AuthProviderModel) Name() string {
	return "auth_providers"
}

func GetAuthProviderModel() AuthProviderModel {
	return AuthProviderModel{
		IdKey:        "id",
		ProviderKey:  "provider",
		IsEnabledKey: "is_enabled",
	}
}
