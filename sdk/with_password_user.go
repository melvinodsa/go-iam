package sdk

import (
	"errors"
	"time"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type WithPasswordUser struct {
	ProjectID string     `json:"project_id"`         // Unique project ID
	Email     string     `json:"email"`              // Unique email
	Password  string     `json:"password,omitempty"` // Hashed password
	CreatedAt *time.Time `json:"created_at"`         // Timestamp when the user was created
	CreatedBy string     `json:"created_by"`         // User who created this user
	UpdatedAt *time.Time `json:"updated_at"`         // Timestamp when the user was last updated
	UpdatedBy string     `json:"updated_by"`         // User who last updated this user
}
