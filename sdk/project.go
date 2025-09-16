package sdk

import (
	"errors"
	"time"
)

// ErrProjectNotFound is returned when a requested project cannot be found.
var ErrProjectNotFound = errors.New("project not found")

// Project represents a project in the Go IAM system.
// Projects provide multi-tenant isolation, ensuring that users, clients,
// and other resources are scoped to specific organizational units.
type Project struct {
	Id          string     `json:"id"`          // Unique identifier for the project
	Name        string     `json:"name"`        // Display name of the project
	Tags        []string   `json:"tags"`        // Tags for categorizing the project
	Description string     `json:"description"` // Description of the project's purpose
	CreatedAt   *time.Time `json:"created_at"`  // Timestamp when project was created
	CreatedBy   string     `json:"created_by"`  // ID of the user who created this project
	UpdatedAt   *time.Time `json:"updated_at"`  // Timestamp when project was last updated
	UpdatedBy   string     `json:"updated_by"`  // ID of the user who last updated this project
}

// ProjectResponse represents an API response containing a single project.
type ProjectResponse struct {
	Success bool     `json:"success"`        // Indicates if the operation was successful
	Message string   `json:"message"`        // Human-readable message about the operation
	Data    *Project `json:"data,omitempty"` // The project data (present only on success)
}

// ProjectsResponse represents an API response containing a list of projects.
type ProjectsResponse struct {
	Success bool      `json:"success"`        // Indicates if the operation was successful
	Message string    `json:"message"`        // Human-readable message about the operation
	Data    []Project `json:"data,omitempty"` // Array of project data
}

// ProjectType is a utility type for type-safe operations involving projects.
type ProjectType struct{}

// ProjectsTypeVal is a global instance of ProjectType for use in type-safe operations.
var ProjectsTypeVal = ProjectType{}
