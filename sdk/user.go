package sdk

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	Id        string                  `json:"id"`
	ProjectId string                  `json:"project_id"`
	Name      string                  `json:"name"`
	Email     string                  `json:"email"`
	Phone     string                  `json:"phone"`
	Enabled   bool                    `json:"enabled"`
	Expiry    *time.Time              `json:"expiry"`
	Roles     map[string]UserRole     `json:"roles"`
	Resources map[string]UserResource `json:"resources"`
	Policies  map[string]string       `json:"policies"`
	CreatedAt *time.Time              `json:"created_at"`
	CreatedBy string                  `json:"created_by"`
	UpdatedAt *time.Time              `json:"updated_at"`
	UpdatedBy string                  `json:"updated_by"`
}

type UserRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserResource struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type UserQuery struct {
	ProjectId   string `json:"project_id"`
	SearchQuery string `json:"search_query"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *User  `json:"data,omitempty"`
}

type UserListResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *map[string]User `json:"data,omitempty"` // Changed to map for consistency
}
