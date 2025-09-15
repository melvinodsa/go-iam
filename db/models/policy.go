package models

import "time"

// Policy represents a resource-based policy that associates roles with resources.
// Policies define fine-grained access control rules that can be applied to users and resources.
type Policy struct {
	Id          string            `bson:"id"`          // Unique identifier for the policy
	Name        string            `bson:"name"`        // Human-readable name of the policy
	Roles       map[string]string `bson:"roles"`       // Map of role IDs to role names associated with this policy
	Description string            `bson:"description"` // Detailed description of the policy's purpose
	CreatedAt   *time.Time        `bson:"created_at"`  // Timestamp when the policy was created
	CreatedBy   string            `bson:"created_by"`  // User who created the policy
}

// PolicyModel provides database access patterns and field mappings for Policy entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type PolicyModel struct {
	iam                   // Embedded struct providing DbName() method
	IdKey          string // BSON field key for policy ID
	NameKey        string // BSON field key for policy name
	RolesKey       string // BSON field key for policy roles
	DescriptionKey string // BSON field key for policy description
}

// Name returns the MongoDB collection name for policies.
// This implements the DbCollection interface.
func (p PolicyModel) Name() string {
	return "policies"
}

// GetPolicyModel returns a properly initialized PolicyModel with all field mappings.
// This function provides a singleton pattern for accessing policy model operations.
//
// Returns a PolicyModel instance with all BSON field keys mapped to their respective field names.
func GetPolicyModel() PolicyModel {
	return PolicyModel{
		IdKey:          "id",
		NameKey:        "name",
		RolesKey:       "roles",
		DescriptionKey: "description",
	}
}
