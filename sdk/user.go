// Package sdk provides client types and structures for the Go IAM API.
//
// This package contains all the data types, request/response structures,
// and utility types needed to interact with the Go IAM identity and access
// management system. It serves as the client-side SDK for applications
// that need to authenticate users, manage permissions, and perform
// identity-related operations.
//
// The package is organized into several key areas:
//   - Authentication: Login, token management, and OAuth2/OIDC flows
//   - User Management: User profiles, roles, and permissions
//   - Client Management: OAuth2 clients and service accounts
//   - Project Management: Multi-tenant project isolation
//   - Authorization: Roles, resources, and policies
//   - Providers: External authentication provider integration
//
// All types in this package are designed to be JSON-serializable and
// are used for both API requests and responses in the Go IAM system.
package sdk

import (
	"errors"
	"time"
)

// ErrUserNotFound is returned when a requested user cannot be found in the system.
var ErrUserNotFound = errors.New("user not found")

// User represents a user entity in the Go IAM system.
// Users are the primary identities that can authenticate and be granted
// permissions within projects. Each user belongs to a specific project
// and can have roles, resources, and policies assigned to them.
type User struct {
	Id             string                  `json:"id"`                         // Unique identifier for the user
	ProjectId      string                  `json:"project_id"`                 // ID of the project this user belongs to
	Name           string                  `json:"name"`                       // Display name of the user
	Email          string                  `json:"email"`                      // Email address (unique within project)
	Phone          string                  `json:"phone"`                      // Phone number (optional)
	Enabled        bool                    `json:"enabled"`                    // Whether the user account is active
	ProfilePic     string                  `json:"profile_pic"`                // URL to the user's profile picture
	LinkedClientId string                  `json:"linked_client_id,omitempty"` // Associated client ID for service accounts
	Expiry         *time.Time              `json:"expiry"`                     // Account expiration time (optional)
	Roles          map[string]UserRole     `json:"roles"`                      // Assigned roles mapped by role ID
	Resources      map[string]UserResource `json:"resources"`                  // Associated resources mapped by resource key
	Policies       map[string]UserPolicy   `json:"policies"`                   // Applied policies mapped by policy name
	CreatedAt      *time.Time              `json:"created_at"`                 // Timestamp when user was created
	CreatedBy      string                  `json:"created_by"`                 // ID of the user who created this user
	UpdatedAt      *time.Time              `json:"updated_at"`                 // Timestamp when user was last updated
	UpdatedBy      string                  `json:"updated_by"`                 // ID of the user who last updated this user
}

// UserPolicy represents a policy assigned to a user with optional argument mappings.
// Policies define permission rules that can be dynamically configured through
// argument substitution.
type UserPolicy struct {
	Name    string            `json:"name"`              // Name of the policy
	Mapping UserPolicyMapping `json:"mapping,omitempty"` // Argument mappings for policy customization
}

// UserPolicyMapping contains argument mappings for dynamic policy evaluation.
// This allows policies to be parameterized with user-specific or context-specific values.
type UserPolicyMapping struct {
	Arguments map[string]UserPolicyMappingValue `json:"arguments,omitempty"` // Named arguments for policy substitution
}

// UserPolicyMappingValue represents a value that can be substituted into a policy.
// Currently supports static string values, but can be extended for dynamic values.
type UserPolicyMappingValue struct {
	Static string `json:"static,omitempty"` // Static string value for substitution
}

// UserRole represents a role assigned to a user.
// Roles are collections of permissions that can be granted to users.
type UserRole struct {
	Id   string `json:"id"`   // Unique identifier of the role
	Name string `json:"name"` // Display name of the role
}

// UserResource represents a resource associated with a user along with
// the roles and policies that apply to that resource.
type UserResource struct {
	RoleIds   map[string]bool `json:"role_ids"`   // Set of role IDs that apply to this resource
	PolicyIds map[string]bool `json:"policy_ids"` // Set of policy IDs that apply to this resource
	Key       string          `json:"key"`        // Unique key identifying the resource
	Name      string          `json:"name"`       // Display name of the resource
}

// AddUserResourceRequest represents a request to associate a resource with a user.
// This includes specifying which role and/or policy should apply to the resource.
type AddUserResourceRequest struct {
	RoleId   string `json:"role_id"`   // ID of the role to apply to the resource
	PolicyId string `json:"policy_id"` // ID of the policy to apply to the resource
	Key      string `json:"key"`       // Unique key of the resource to associate
	Name     string `json:"name"`      // Display name of the resource
}

// UserQuery represents search and filtering criteria for user queries.
// This is used for listing users with various filters and pagination.
type UserQuery struct {
	ProjectIds  []string `json:"project_ids"`  // Filter by specific project IDs
	RoleId      string   `json:"role_id"`      // Filter by users having a specific role
	SearchQuery string   `json:"search_query"` // Text search across user fields
	Skip        int64    `json:"skip"`         // Number of records to skip (pagination)
	Limit       int64    `json:"limit"`        // Maximum number of records to return
}

// UserResponse represents a standard API response containing a single user.
type UserResponse struct {
	Success bool   `json:"success"`        // Indicates if the operation was successful
	Message string `json:"message"`        // Human-readable message about the operation
	Data    *User  `json:"data,omitempty"` // The user data (present only on success)
}

// DashboardUserResponse represents a specialized response for dashboard user data.
// This includes additional setup information beyond the basic user data.
type DashboardUserResponse struct {
	Success bool   `json:"success"` // Indicates if the operation was successful
	Message string `json:"message"` // Human-readable message about the operation
	Data    struct {
		User  *User `json:"user"` // The user data
		Setup struct {
			ClientAdded bool   `json:"client_added"` // Whether a client has been set up for this user
			ClientId    string `json:"client_id"`    // ID of the associated client (if any)
		} `json:"setup"` // Setup and configuration information
	} `json:"data,omitempty"` // Combined user and setup data
}

// UserList represents a paginated list of users with metadata.
type UserList struct {
	Users []User `json:"users"` // Array of user objects
	Total int64  `json:"total"` // Total number of users matching the query (before pagination)
	Skip  int64  `json:"skip"`  // Number of records skipped
	Limit int64  `json:"limit"` // Maximum number of records returned
}

// UserListResponse represents an API response containing a list of users.
type UserListResponse struct {
	Success bool      `json:"success"` // Indicates if the operation was successful
	Message string    `json:"message"` // Human-readable message about the operation
	Data    *UserList `json:"data"`    // The paginated user list data
}

// UserRoleUpdate represents changes to be made to a user's role assignments.
// This supports both adding and removing roles in a single operation.
type UserRoleUpdate struct {
	ToBeAdded   []string `json:"to_be_added"`   // Array of role IDs to assign to the user
	ToBeRemoved []string `json:"to_be_removed"` // Array of role IDs to remove from the user
}

// UserPolicyUpdate represents changes to be made to a user's policy assignments.
// This supports both adding and removing policies in a single operation.
type UserPolicyUpdate struct {
	ToBeAdded   map[string]UserPolicy `json:"to_be_added"`   // Map of policy names to UserPolicy objects to assign
	ToBeRemoved []string              `json:"to_be_removed"` // Array of policy names to remove from the user
}

// UserType is a utility type for type-safe operations involving users.
type UserType struct{}

// UserTypeVal is a global instance of UserType for use in type-safe operations.
var UserTypeVal = UserType{}
