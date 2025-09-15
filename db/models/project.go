package models

import "time"

// Project represents a project entity in the Go IAM system.
// Projects are organizational units that contain users, clients, roles, and resources.
// They provide isolation and multi-tenancy in the IAM system.
type Project struct {
	Id          string     `bson:"id"`          // Unique identifier for the project
	Name        string     `bson:"name"`        // Human-readable name of the project
	Tags        []string   `bson:"tags"`        // Tags for categorizing and filtering projects
	Description string     `bson:"description"` // Detailed description of the project's purpose
	CreatedAt   *time.Time `bson:"created_at"`  // Timestamp when the project was created
	CreatedBy   string     `bson:"created_by"`  // User who created the project
	UpdatedAt   *time.Time `bson:"updated_at"`  // Timestamp when the project was last updated
	UpdatedBy   string     `bson:"updated_by"`  // User who last updated the project
}

// ProjectModel provides database access patterns and field mappings for Project entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type ProjectModel struct {
	iam                   // Embedded struct providing DbName() method
	IdKey          string // BSON field key for project ID
	NameKey        string // BSON field key for project name
	TagsKey        string // BSON field key for project tags
	DescriptionKey string // BSON field key for project description
}

// Name returns the MongoDB collection name for projects.
// This implements the DbCollection interface.
func (p ProjectModel) Name() string {
	return "projects"
}

// GetProjectModel returns a properly initialized ProjectModel with all field mappings.
// This function provides a singleton pattern for accessing project model operations.
//
// Returns a ProjectModel instance with all BSON field keys mapped to their respective field names.
func GetProjectModel() ProjectModel {
	return ProjectModel{
		IdKey:          "id",
		NameKey:        "name",
		TagsKey:        "tags",
		DescriptionKey: "description",
	}
}
