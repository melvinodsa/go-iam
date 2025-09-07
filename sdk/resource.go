package sdk

import (
	"errors"
	"time"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type Resource struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Key         string     `json:"key"`
	Enabled     bool       `json:"enabled"`
	ProjectId   string     `json:"project_id"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   string     `json:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ResourceQuery struct {
	ProjectIds  []string `json:"project_ids,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Key         string   `json:"key,omitempty"`
	Skip        int64    `json:"skip"`
	Limit       int64    `json:"limit"`
}

type ResourceResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *Resource `json:"data,omitempty"`
}

type ResourceList struct {
	Resources []Resource `json:"resources"`
	Total     int64      `json:"total"`
	Skip      int64      `json:"skip"`
	Limit     int64      `json:"limit"`
}

type ResourcesResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *ResourceList `json:"data,omitempty"`
}
