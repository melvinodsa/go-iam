package models

import (
	"time"
)

type Role struct {
	Id        string      `json:"id"`
	ProjectId string      `json:"project_id"`
	Name      string      `json:"name"`
	Resources []Resources `json:"resources"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt time.Time   `json:"updated_at"`
	UpdatedBy string      `json:"updated_by"`
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
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
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
