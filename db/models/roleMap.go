package models

// RoleMap represents a mapping between roles and users in the Go IAM system.
// This provides a many-to-many relationship between roles and users,
// allowing efficient querying of user-role associations.
type RoleMap struct {
	RoleId string   `bson:"role_id"` // ID of the role in the mapping
	UserId []string `bson:"user_id"` // Array of user IDs assigned to this role
}

// RoleMapModel provides database access patterns and field mappings for RoleMap entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type RoleMapModel struct {
	iam              // Embedded struct providing DbName() method
	RoleIdKey string // BSON field key for role ID
	UserIdKey string // BSON field key for user ID array
}

// Name returns the MongoDB collection name for role mappings.
// This implements the DbCollection interface.
func (u RoleMapModel) Name() string {
	return "roleMap"
}

// GetRoleMap returns a properly initialized RoleMapModel with all field mappings.
// This function provides a singleton pattern for accessing role map model operations.
//
// Returns a RoleMapModel instance with all BSON field keys mapped to their respective field names.
func GetRoleMap() RoleMapModel {
	return RoleMapModel{
		RoleIdKey: "role_id",
		UserIdKey: "user_id",
	}
}
