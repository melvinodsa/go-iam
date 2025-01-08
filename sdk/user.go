package sdk

import "time"

type User struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"project_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Enabled   bool       `json:"enabled"`
	Expiry    *time.Time `json:"expiry"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
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
