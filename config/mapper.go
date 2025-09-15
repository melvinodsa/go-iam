package config

// Context keys used for storing resource and role mappings in Fiber request contexts.
// These keys are used throughout the application to access cached mapping data.
var (
	// ResourceMapContextKey is the context key used to store resource mappings
	// in Fiber request contexts. Used for caching resource data per request.
	ResourceMapContextKey = "resourceMap"

	// RoleMapContextKey is the context key used to store role mappings
	// in Fiber request contexts. Used for caching role data per request.
	RoleMapContextKey = "roleMap"
)
