package models

import "time"

// Policy represents a resource-based policy that associates roles with resources

type Policy struct {
	Id          string            `bson:"id"`
	Name        string            `bson:"name"`
	Roles       map[string]string `bson:"roles"`
	Description string            `bson:"description"`
	CreatedAt   *time.Time        `bson:"created_at"`
	CreatedBy   string            `bson:"created_by"`
}

type PolicyModel struct {
	iam
	IdKey          string
	NameKey        string
	RolesKey       string
	DescriptionKey string
}

func (p PolicyModel) Name() string {
	return "policies"
}

func GetPolicyModel() PolicyModel {
	return PolicyModel{
		IdKey:          "id",
		NameKey:        "name",
		RolesKey:       "roles",
		DescriptionKey: "description",
	}
}
