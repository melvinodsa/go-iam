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
	ProjectIds  []string `json:"project_ids"`
	SearchQuery string   `json:"search_query"`
	Skip        int64    `json:"skip"`
	Limit       int64    `json:"limit"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *User  `json:"data,omitempty"`
}

type DashboardUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		User  *User `json:"user"`
		Setup struct {
			ClientAdded bool   `json:"client_added"`
			ClientId    string `json:"client_id"`
		} `json:"setup"`
	} `json:"data,omitempty"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
	Skip  int64  `json:"skip"`
	Limit int64  `json:"limit"`
}

type UserListResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *UserList `json:"data"` // Changed to map for consistency
}

type UserRoleUpdate struct {
	ToBeAdded   []string `json:"to_be_added"`
	ToBeRemoved []string `json:"to_be_removed"`
}
