package sdk

// RoleMap represents a mapping between roles and users.
// This structure is used to define which users are assigned to specific roles,
// enabling efficient role-based access control queries and operations.
type RoleMap struct {
	Role_id string   `json:"resource_id"` // ID of the role being mapped
	User_id []string `json:"role_id"`     // Array of user IDs assigned to this role
}
