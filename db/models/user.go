package models

import "time"

type User struct {
	Id        string     `bson:"id"`
	ProjectId string     `bson:"project_id"`
	Name      string     `bson:"name"`
	Email     string     `bson:"email"`
	IsEnabled bool       `bson:"is_enabled"`
	Expiry    *time.Time `bson:"expiry"`
	CreatedAt *time.Time `bson:"created_at"`
	CreatedBy string     `bson:"created_by"`
	UpdatedAt *time.Time `bson:"updated_at"`
	UpdatedBy string     `bson:"updated_by"`
}

type UserModel struct {
	iam
	IdKey        string
	EmailKey     string
	IsEnabledKey string
	ExpiryKey    string
}

func (u UserModel) Name() string {
	return "users"
}

func GetUserModel() UserModel {
	return UserModel{
		IdKey:        "id",
		EmailKey:     "email",
		IsEnabledKey: "is_enabled",
		ExpiryKey:    "expiry",
	}
}
