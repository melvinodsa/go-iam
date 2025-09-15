package models

import (
	"time"
)

// Role represents a role entity in the Go IAM system.
// Roles define collections of permissions that can be assigned to users.
// Each role belongs to a project and can have access to multiple resources.
type Role struct {
	Id          string               `bson:"id"`          // Unique identifier for the role
	ProjectId   string               `bson:"project_id"`  // ID of the project this role belongs to
	Name        string               `bson:"name"`        // Human-readable name of the role
	Description string               `bson:"description"` // Detailed description of the role's purpose
	Resources   map[string]Resources `bson:"resources"`   // Map of resources this role has access to
	Enabled     bool                 `bson:"enabled"`     // Whether the role is currently active
	CreatedAt   time.Time            `bson:"created_at"`  // Timestamp when the role was created
	CreatedBy   string               `bson:"created_by"`  // User who created the role
	UpdatedAt   time.Time            `bson:"updated_at"`  // Timestamp when the role was last updated
	UpdatedBy   string               `bson:"updated_by"`  // User who last updated the role
}

// RoleModel provides database access patterns and field mappings for Role entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type RoleModel struct {
	iam                   // Embedded struct providing DbName() method
	IdKey          string // BSON field key for role ID
	ProjectIdKey   string // BSON field key for project ID
	NameKey        string // BSON field key for role name
	DescriptionKey string // BSON field key for role description
	ResourcesKey   string // BSON field key for role resources
	CreatedAtKey   string // BSON field key for creation timestamp
	CreatedByKey   string // BSON field key for creator
	UpdatedAtKey   string // BSON field key for update timestamp
	EnabledKey     string // BSON field key for enabled status
	UpdatedByKey   string // BSON field key for updater
}

// Name returns the MongoDB collection name for roles.
// This implements the DbCollection interface.
func (u RoleModel) Name() string {
	return "roles"
}

// Resources represents a resource that can be associated with a role.
// Resources define the entities that roles can have permissions on.
type Resources struct {
	Id   string `bson:"id"`   // Unique identifier of the resource
	Key  string `bson:"key"`  // Unique key identifier for the resource
	Name string `bson:"name"` // Human-readable name of the resource
}

// GetRoleModel returns a properly initialized RoleModel with all field mappings.
// This function provides a singleton pattern for accessing role model operations.
//
// Returns a RoleModel instance with all BSON field keys mapped to their respective field names.
func GetRoleModel() RoleModel {
	return RoleModel{
		IdKey:          "id",
		ProjectIdKey:   "project_id",
		NameKey:        "name",
		DescriptionKey: "description",
		ResourcesKey:   "resources",
		CreatedAtKey:   "created_at",
		CreatedByKey:   "created_by",
		UpdatedAtKey:   "updated_at",
		UpdatedByKey:   "updated_by",
		EnabledKey:     "enabled",
	}
}
