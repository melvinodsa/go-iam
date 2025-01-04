package models

import "time"

type Client struct {
	Id                    string     `bson:"id"`
	Name                  string     `bson:"name"`
	Description           string     `bson:"description"`
	Secret                string     `bson:"secret"`
	Tags                  []string   `bson:"tags"`
	RedirectURLs          []string   `bson:"redirect_urls"`
	DefaultAuthProviderId string     `bson:"default_auth_provider_id"`
	ProjectId             string     `bson:"project_id"`
	Scopes                []string   `bson:"scopes"`
	Enabled               bool       `bson:"enabled"`
	CreatedAt             *time.Time `bson:"created_at"`
	CreatedBy             string     `bson:"created_by"`
	UpdatedAt             *time.Time `bson:"updated_at"`
	UpdatedBy             string     `bson:"updated_by"`
}

type ClientModel struct {
	iam
	IdKey          string
	NameKey        string
	TagsKey        string
	DescriptionKey string
	ProjectIdKey   string
}

func (c ClientModel) Name() string {
	return "clients"
}

func GetClientModel() ClientModel {
	return ClientModel{
		IdKey:          "id",
		NameKey:        "name",
		TagsKey:        "tags",
		DescriptionKey: "description",
		ProjectIdKey:   "project_id",
	}
}
