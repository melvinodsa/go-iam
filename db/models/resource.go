package models

import "time"

type Resource struct {
	ID          string     `bson:"id,omitempty"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	Key         string     `bson:"key"`
	Enabled     bool       `bson:"enabled"`
	CreatedAt   *time.Time `bson:"created_at"`
	CreatedBy   string     `bson:"created_by"`
	UpdatedAt   *time.Time `bson:"updated_at"`
	UpdatedBy   string     `bson:"updated_by"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty"`
}

type ResourceModel struct {
	iam
	IdKey          string
	NameKey        string
	DescriptionKey string
	KeyKey         string
	EnabledKey     string
}

func (r ResourceModel) Name() string {
	return "resources"
}

func GetResourceModel() ResourceModel {
	return ResourceModel{
		NameKey:        "name",
		DescriptionKey: "description",
		KeyKey:         "key",
		EnabledKey:     "enabled",
	}
}
