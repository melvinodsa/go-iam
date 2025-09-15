package models

import "time"

// Resource represents a resource entity in the Go IAM system.
// Resources are entities that can be protected and accessed through the IAM system.
// They can be assigned to roles and have policies applied to control access.
type Resource struct {
	ID          string     `bson:"id,omitempty"`         // Unique identifier for the resource
	Name        string     `bson:"name"`                 // Human-readable name of the resource
	Description string     `bson:"description"`          // Detailed description of the resource
	Key         string     `bson:"key"`                  // Unique key identifier for the resource
	ProjectId   string     `bson:"project_id"`           // ID of the project this resource belongs to
	Enabled     bool       `bson:"enabled"`              // Whether the resource is currently active
	CreatedAt   *time.Time `bson:"created_at"`           // Timestamp when the resource was created
	CreatedBy   string     `bson:"created_by"`           // User who created the resource
	UpdatedAt   *time.Time `bson:"updated_at"`           // Timestamp when the resource was last updated
	UpdatedBy   string     `bson:"updated_by"`           // User who last updated the resource
	DeletedAt   *time.Time `bson:"deleted_at,omitempty"` // Timestamp when the resource was soft deleted
}

// ResourceModel provides database access patterns and field mappings for Resource entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type ResourceModel struct {
	iam                   // Embedded struct providing DbName() method
	IdKey          string // BSON field key for resource ID
	NameKey        string // BSON field key for resource name
	DescriptionKey string // BSON field key for resource description
	KeyKey         string // BSON field key for resource key
	EnabledKey     string // BSON field key for enabled status
	ProjectIdKey   string // BSON field key for project ID
}

// Name returns the MongoDB collection name for resources.
// This implements the DbCollection interface.
func (r ResourceModel) Name() string {
	return "resources"
}

// GetResourceModel returns a properly initialized ResourceModel with all field mappings.
// This function provides a singleton pattern for accessing resource model operations.
//
// Returns a ResourceModel instance with all BSON field keys mapped to their respective field names.
func GetResourceModel() ResourceModel {
	return ResourceModel{
		IdKey:          "id",
		NameKey:        "name",
		DescriptionKey: "description",
		KeyKey:         "key",
		EnabledKey:     "enabled",
		ProjectIdKey:   "project_id",
	}
}
