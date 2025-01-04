package models

import "time"

type AuthProviderType string

type AuthProvider struct {
	Id        string              `bson:"id"`
	Name      string              `bson:"name"`
	Icon      string              `bson:"icon"`
	Provider  AuthProviderType    `bson:"provider"`
	Params    []AuthProviderParam `bson:"params"`
	ProjectId string              `bson:"project_id"`
	Enabled   bool                `bson:"enabled"`
	CreatedAt *time.Time          `bson:"created_at"`
	UpdatedAt *time.Time          `bson:"updated_at"`
	CreatedBy string              `bson:"created_by"`
	UpdatedBy string              `bson:"updated_by"`
}

type AuthProviderParam struct {
	Label    string `bson:"label"`
	Value    string `bson:"value"`
	Key      string `bson:"key"`
	IsSecret bool   `bson:"is_secret"`
}

type AuthProviderModel struct {
	iam
	IdKey        string
	ProviderKey  string
	IsEnabledKey string
	ProjectIdKey string
}

func (a AuthProviderModel) Name() string {
	return "auth_providers"
}

func GetAuthProviderModel() AuthProviderModel {
	return AuthProviderModel{
		IdKey:        "id",
		ProviderKey:  "provider",
		IsEnabledKey: "is_enabled",
		ProjectIdKey: "project_id",
	}
}
