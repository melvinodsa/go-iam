package models

import "time"

type WithPasswordUser struct {
	ID        string     `bson:"id"`         // Unique identifier
	ProjectID string     `bson:"project_id"` // Unique project ID
	Email     string     `bson:"email"`      // Unique email
	Password  string     `bson:"password"`   // Hashed password
	CreatedAt *time.Time `bson:"created_at"` // Timestamp when the user was created
	CreatedBy string     `bson:"created_by"` // User who created this user
	UpdatedAt *time.Time `bson:"updated_at"` // Timestamp when the user was last updated
	UpdatedBy string     `bson:"updated_by"` // User who last updated this user
}

type WithPasswordUserModel struct {
	iam
	ProjectIDKey string // BSON key for project_id field
	EmailKey     string // BSON key for email field
	PasswordKey  string // BSON key for hashed password field
	UpdatedAtKey string // BSON key for updated_at field
	UpdatedByKey string // BSON key for updated_by field
}

func GetWithPasswordUserModel() *WithPasswordUserModel {
	return &WithPasswordUserModel{
		ProjectIDKey: "project_id",
		EmailKey:     "email",
		PasswordKey:  "password",
		UpdatedAtKey: "updated_at",
		UpdatedByKey: "updated_by",
	}
}

func (u *WithPasswordUserModel) Name() string {
	return "with_password_users"
}
