package models

import "time"

type User struct {
	Id         string                  `bson:"id"`
	ProjectId  string                  `bson:"project_id"`
	Name       string                  `bson:"name"`
	Email      string                  `bson:"email"`
	Phone      string                  `bson:"phone"`
	Enabled    bool                    `bson:"enabled"`
	ProfilePic string                  `bson:"profile_pic"`
	Expiry     *time.Time              `bson:"expiry"`
	Roles      map[string]UserRoles    `bson:"roles"`
	Resources  map[string]UserResource `bson:"resources"`
	Policies   map[string]string       `bson:"policies"`
	CreatedAt  *time.Time              `bson:"created_at"`
	CreatedBy  string                  `bson:"created_by"`
	UpdatedAt  *time.Time              `bson:"updated_at"`
	UpdatedBy  string                  `bson:"updated_by"`
}

type UserResource struct {
	Id   string `bson:"id"`
	Key  string `bson:"key"`
	Name string `bson:"name"`
}

type UserRoles struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type UserModel struct {
	iam
	IdKey         string
	NameKey       string
	EmailKey      string
	PhoneKey      string
	EnabledKey    string
	RolesIdKey    string
	PoliciesKey   string
	ResourceIdKey string
	IsEnabledKey  string
	ProjectIDKey  string
	ExpiryKey     string
}

func (u UserModel) Name() string {
	return "users"
}

func GetUserModel() UserModel {
	return UserModel{
		IdKey:         "id",
		NameKey:       "name",
		EmailKey:      "email",
		PhoneKey:      "phone",
		EnabledKey:    "enabled",
		RolesIdKey:    "roles",
		ResourceIdKey: "resources",
		PoliciesKey:   "policies",
		IsEnabledKey:  "is_enabled",
		ProjectIDKey:  "project_id",
		ExpiryKey:     "expiry",
	}
}
