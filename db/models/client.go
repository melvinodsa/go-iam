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
	GoIamClient           bool       `bson:"go_iam_client"` // Indicates if this is a Go-IAM client
	ProjectId             string     `bson:"project_id"`
	ServiceAccountEmail   string     `bson:"service_account_email"`
	Scopes                []string   `bson:"scopes"`
	Enabled               bool       `bson:"enabled"`
	LinkedUserId          string     `bson:"linked_user_id"`
	CreatedAt             *time.Time `bson:"created_at"`
	CreatedBy             string     `bson:"created_by"`
	UpdatedAt             *time.Time `bson:"updated_at"`
	UpdatedBy             string     `bson:"updated_by"`
}

type ClientModel struct {
	iam
	IdKey           string
	NameKey         string
	TagsKey         string
	DescriptionKey  string
	ProjectIdKey    string
	GoIamClientKey  string // Indicates if this is a Go-IAM client
	LinkedUserIdKey string // For service account authentication
	UpdatedAtKey    string
}

func (c ClientModel) Name() string {
	return "clients"
}

func GetClientModel() ClientModel {
	return ClientModel{
		IdKey:           "id",
		NameKey:         "name",
		TagsKey:         "tags",
		DescriptionKey:  "description",
		ProjectIdKey:    "project_id",
		GoIamClientKey:  "go_iam_client",
		LinkedUserIdKey: "linked_user_id",
		UpdatedAtKey:    "updated_at",
	}
}
