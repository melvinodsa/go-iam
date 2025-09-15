package models

import "time"

// User represents a user entity in the Go IAM system.
// Users are the primary subjects of authentication and authorization,
// with assigned roles, resources, and policies that determine their access rights.
type User struct {
	Id             string                  `bson:"id"`                         // Unique identifier for the user
	ProjectId      string                  `bson:"project_id"`                 // ID of the project this user belongs to
	Name           string                  `bson:"name"`                       // Display name of the user
	Email          string                  `bson:"email"`                      // Email address of the user
	Phone          string                  `bson:"phone"`                      // Phone number of the user
	Enabled        bool                    `bson:"enabled"`                    // Whether the user account is active
	ProfilePic     string                  `bson:"profile_pic"`                // URL or path to the user's profile picture
	Expiry         *time.Time              `bson:"expiry"`                     // Optional expiration date for the user account
	Roles          map[string]UserRoles    `bson:"roles"`                      // Roles assigned to the user
	Resources      map[string]UserResource `bson:"resources"`                  // Resources the user has access to
	Policies       map[string]UserPolicy   `bson:"policies"`                   // Policies applied to the user
	LinkedClientId string                  `bson:"linked_client_id,omitempty"` // Client ID for service account users
	CreatedAt      *time.Time              `bson:"created_at"`                 // Timestamp when the user was created
	CreatedBy      string                  `bson:"created_by"`                 // User who created this user
	UpdatedAt      *time.Time              `bson:"updated_at"`                 // Timestamp when the user was last updated
	UpdatedBy      string                  `bson:"updated_by"`                 // User who last updated this user
}

// UserPolicy represents a policy assignment to a user with dynamic value mapping.
// Policies define fine-grained permissions and can have configurable arguments.
type UserPolicy struct {
	Name    string            `bson:"name,omitempty"`    // Name of the policy
	Mapping UserPolicyMapping `bson:"mapping,omitempty"` // Dynamic value mappings for policy arguments
}

// UserPolicyMapping contains argument mappings for policy execution.
// This allows policies to have dynamic values based on user context.
type UserPolicyMapping struct {
	Arguments map[string]UserPolicyMappingValue `bson:"arguments,omitempty"` // Argument name to value mappings
}

// UserPolicyMappingValue represents a mapped value for policy arguments.
// Currently supports static values, but can be extended for dynamic values.
type UserPolicyMappingValue struct {
	Static string `bson:"static,omitempty"` // Static value for the policy argument
}

// UserResource represents a resource that a user has access to.
// Resources can have associated roles and policies that define the user's permissions.
type UserResource struct {
	RoleIds   map[string]bool `bson:"role_ids"`   // Map of role IDs assigned to this resource
	PolicyIds map[string]bool `bson:"policy_ids"` // Map of policy IDs applied to this resource
	Key       string          `bson:"key"`        // Unique key identifier for the resource
	Name      string          `bson:"name"`       // Human-readable name of the resource
}

// UserRoles represents a role assignment to a user.
// Roles define collections of permissions that can be assigned to users.
type UserRoles struct {
	Id   string `bson:"id"`   // Unique identifier of the role
	Name string `bson:"name"` // Human-readable name of the role
}

// UserModel provides database access patterns and field mappings for User entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
// UserModel provides database access patterns and field mappings for User entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type UserModel struct {
	iam                 // Embedded struct providing DbName() method
	IdKey        string // BSON field key for user ID
	NameKey      string // BSON field key for user name
	EmailKey     string // BSON field key for user email
	PhoneKey     string // BSON field key for user phone
	EnabledKey   string // BSON field key for enabled status
	RolesIdKey   string // BSON field key for user roles
	PoliciesKey  string // BSON field key for user policies
	ResourcesKey string // BSON field key for user resources
	IsEnabledKey string // BSON field key for enabled status (alternative)
	ProjectIDKey string // BSON field key for project ID
	ExpiryKey    string // BSON field key for account expiry
}

// Name returns the MongoDB collection name for users.
// This implements the DbCollection interface.
func (u UserModel) Name() string {
	return "users"
}

// GetUserModel returns a properly initialized UserModel with all field mappings.
// This function provides a singleton pattern for accessing user model operations.
//
// Returns a UserModel instance with all BSON field keys mapped to their respective field names.
func GetUserModel() UserModel {
	return UserModel{
		IdKey:        "id",
		NameKey:      "name",
		EmailKey:     "email",
		PhoneKey:     "phone",
		EnabledKey:   "enabled",
		RolesIdKey:   "roles",
		ResourcesKey: "resources",
		PoliciesKey:  "policies",
		IsEnabledKey: "is_enabled",
		ProjectIDKey: "project_id",
		ExpiryKey:    "expiry",
	}
}
