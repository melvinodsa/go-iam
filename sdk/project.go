package sdk

import (
	"errors"
	"time"
)

var ErrProjectNotFound = errors.New("project not found")

type Project struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Tags        []string   `json:"tags"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   string     `json:"updated_by"`
}

type ProjectResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    *Project `json:"data,omitempty"`
}

type ProjectsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    []Project `json:"data,omitempty"`
}

type ProjectType struct{}

var ProjectsTypeVal = ProjectType{}
