package sdk

import (
	"errors"
	"time"
)

var ErrRoleNotFound = errors.New("role not found")

type Role struct {
	Id        string               `json:"id"`
	ProjectId string               `json:"project_id"`
	Name      string               `json:"name"`
	Resources map[string]Resources `json:"resources"`
	CreatedAt *time.Time           `json:"created_at"`
	CreatedBy string               `json:"created_by"`
	UpdatedAt *time.Time           `json:"updated_at"`
	UpdatedBy string               `json:"updated_by"`
}

type Resources struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type RoleQuery struct {
	ProjectId   string `json:"project_id"`
	SearchQuery string `json:"search_query"`
}

type RoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Role  `json:"data,omitempty"`
}

type RoleListResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *[]Role `json:"data,omitempty"`
}
