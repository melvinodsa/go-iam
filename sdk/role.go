package sdk

import (
	"errors"
	"time"
)

// ErrRoleNotFound is returned when a requested role cannot be found.
var ErrRoleNotFound = errors.New("role not found")

// Role represents a role in the Go IAM system.
// Roles are collections of permissions that can be assigned to users.
// Each role is associated with specific resources and defines what
// actions can be performed on those resources.
type Role struct {
	Id          string               `json:"id"`          // Unique identifier for the role
	ProjectId   string               `json:"project_id"`  // ID of the project this role belongs to
	Description string               `json:"description"` // Description of the role's purpose
	Name        string               `json:"name"`        // Display name of the role
	Resources   map[string]Resources `json:"resources"`   // Map of resource keys to resource definitions
	Enabled     bool                 `json:"enabled"`     // Whether this role is active
	CreatedAt   *time.Time           `json:"created_at"`  // Timestamp when role was created
	CreatedBy   string               `json:"created_by"`  // ID of the user who created this role
	UpdatedAt   *time.Time           `json:"updated_at"`  // Timestamp when role was last updated
	UpdatedBy   string               `json:"updated_by"`  // ID of the user who last updated this role
}

// Resources represents a resource definition within a role.
// This defines which specific resource a role grants access to.
type Resources struct {
	Id   string `json:"id"`   // Unique identifier of the resource
	Key  string `json:"key"`  // Key identifying the resource type/category
	Name string `json:"name"` // Display name of the resource
}

// RoleQuery represents search and filtering criteria for role queries.
// This is used for listing roles with various filters and pagination.
type RoleQuery struct {
	ProjectIds  []string `json:"project_ids"`  // Filter by specific project IDs
	SearchQuery string   `json:"search_query"` // Text search across role fields
	Skip        int64    `json:"skip"`         // Number of records to skip (pagination)
	Limit       int64    `json:"limit"`        // Maximum number of records to return
}

// RoleResponse represents an API response containing a single role.
type RoleResponse struct {
	Success bool   `json:"success"`        // Indicates if the operation was successful
	Message string `json:"message"`        // Human-readable message about the operation
	Data    *Role  `json:"data,omitempty"` // The role data (present only on success)
}

// RoleList represents a paginated list of roles with metadata.
type RoleList struct {
	Roles []Role `json:"roles"` // Array of role objects
	Total int64  `json:"total"` // Total number of roles matching the query (before pagination)
	Skip  int64  `json:"skip"`  // Number of records skipped
	Limit int64  `json:"limit"` // Maximum number of records returned
}

// RoleListResponse represents an API response containing a list of roles.
type RoleListResponse struct {
	Success bool      `json:"success"`        // Indicates if the operation was successful
	Message string    `json:"message"`        // Human-readable message about the operation
	Data    *RoleList `json:"data,omitempty"` // The paginated role list data
}
