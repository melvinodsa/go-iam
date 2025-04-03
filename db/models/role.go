package models

import (
	"time"
)

type Role struct {
	Id        string               `bson:"id"`
	ProjectId string               `bson:"project_id"`
	Name      string               `bson:"name"`
	Resources map[string]Resources `bson:"resources"`
	CreatedAt time.Time            `bson:"created_at"`
	CreatedBy string               `bson:"created_by"`
	UpdatedAt time.Time            `bson:"updated_at"`
	UpdatedBy string               `bson:"updated_by"`
}

type RoleModel struct {
	iam
	IdKey        string
	ProjectIdKey string
	NameKey      string
	ResourcesKey string
	CreatedAtKey string
	CreatedByKey string
	UpdatedAtKey string
	UpdatedByKey string
}

func (u RoleModel) Name() string {
	return "roles"
}

type Resources struct {
	Id   string `bson:"id"`
	Key  string `bson:"key"`
	Name string `bson:"name"`
}

func GetRoleModel() RoleModel {
	return RoleModel{
		IdKey:        "id",
		ProjectIdKey: "project_id",
		NameKey:      "name",
		ResourcesKey: "resources",
		CreatedAtKey: "created_at",
		CreatedByKey: "created_by",
		UpdatedAtKey: "updated_at",
		UpdatedByKey: "updated_by",
	}
}
