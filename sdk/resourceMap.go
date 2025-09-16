package sdk

// ResourceMap represents a mapping between resources and roles.
// This structure is used to define which roles have access to specific resources,
// enabling efficient resource-based access control queries and authorization decisions.
type ResourceMap struct {
	Resource_id string   `json:"resource_id"` // ID of the resource being mapped
	Role_id     []string `json:"role_id"`     // Array of role IDs that have access to this resource
}
