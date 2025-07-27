package sdk

import (
	"errors"
	"time"
)

var ErrRoleNotFound = errors.New("role not found")

type Role struct {
	Id          string               `json:"id"`
	ProjectId   string               `json:"project_id"`
	Description string               `json:"description"`
	Name        string               `json:"name"`
	Resources   map[string]Resources `json:"resources"`
	Enabled     bool                 `json:"enabled"`
	CreatedAt   *time.Time           `json:"created_at"`
	CreatedBy   string               `json:"created_by"`
	UpdatedAt   *time.Time           `json:"updated_at"`
	UpdatedBy   string               `json:"updated_by"`
}

type Resources struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type RoleQuery struct {
	ProjectIds  []string `json:"project_ids"`
	SearchQuery string   `json:"search_query"`
	Skip        int64    `json:"skip"`
	Limit       int64    `json:"limit"`
}

type RoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Role  `json:"data,omitempty"`
}

type RoleList struct {
	Roles []Role `json:"roles"`
	Total int64  `json:"total"`
	Skip  int64  `json:"skip"`
	Limit int64  `json:"limit"`
}

type RoleListResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *RoleList `json:"data,omitempty"`
}

const (
	EventRoleUpdated = "role:updated"
)
