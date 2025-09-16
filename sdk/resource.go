package sdk

import (
	"errors"
	"time"
)

var (
	// ErrResourceNotFound is returned when a requested resource cannot be found.
	ErrResourceNotFound = errors.New("resource not found")
)

// Resource represents a resource in the Go IAM system.
// Resources are entities that can be protected by access control policies.
// They can represent anything from API endpoints to data objects, files,
// or any other system component that needs authorization.
type Resource struct {
	ID          string     `json:"id"`                   // Unique identifier for the resource
	Name        string     `json:"name"`                 // Display name of the resource
	Description string     `json:"description"`          // Description of what this resource represents
	Key         string     `json:"key"`                  // Unique key identifying the resource type/category
	Enabled     bool       `json:"enabled"`              // Whether this resource is active
	ProjectId   string     `json:"project_id"`           // ID of the project this resource belongs to
	CreatedAt   *time.Time `json:"created_at"`           // Timestamp when resource was created
	CreatedBy   string     `json:"created_by"`           // ID of the user who created this resource
	UpdatedAt   *time.Time `json:"updated_at"`           // Timestamp when resource was last updated
	UpdatedBy   string     `json:"updated_by"`           // ID of the user who last updated this resource
	DeletedAt   *time.Time `json:"deleted_at,omitempty"` // Timestamp when resource was deleted (soft delete)
}

// ResourceQuery represents search and filtering criteria for resource queries.
// This is used for listing resources with various filters and pagination.
type ResourceQuery struct {
	ProjectIds  []string `json:"project_ids,omitempty"` // Filter by specific project IDs
	Name        string   `json:"name,omitempty"`        // Filter by resource name (partial match)
	Description string   `json:"description,omitempty"` // Filter by resource description (partial match)
	Key         string   `json:"key,omitempty"`         // Filter by resource key (partial match)
	Skip        int64    `json:"skip"`                  // Number of records to skip (pagination)
	Limit       int64    `json:"limit"`                 // Maximum number of records to return
}

// ResourceResponse represents an API response containing a single resource.
type ResourceResponse struct {
	Success bool      `json:"success"`        // Indicates if the operation was successful
	Message string    `json:"message"`        // Human-readable message about the operation
	Data    *Resource `json:"data,omitempty"` // The resource data (present only on success)
}

// ResourceList represents a paginated list of resources with metadata.
type ResourceList struct {
	Resources []Resource `json:"resources"` // Array of resource objects
	Total     int64      `json:"total"`     // Total number of resources matching the query (before pagination)
	Skip      int64      `json:"skip"`      // Number of records skipped
	Limit     int64      `json:"limit"`     // Maximum number of records returned
}

// ResourcesResponse represents an API response containing a list of resources.
type ResourcesResponse struct {
	Success bool          `json:"success"`        // Indicates if the operation was successful
	Message string        `json:"message"`        // Human-readable message about the operation
	Data    *ResourceList `json:"data,omitempty"` // The paginated resource list data
}
