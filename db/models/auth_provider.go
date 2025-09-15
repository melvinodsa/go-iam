package models

import "time"

// AuthProviderType represents the type of authentication provider.
// This defines the specific implementation used for authentication.
type AuthProviderType string

// AuthProvider represents an authentication provider in the Go IAM system.
// Auth providers handle external authentication services like Google, GitHub, etc.
// Each provider belongs to a project and can be configured with custom parameters.
type AuthProvider struct {
	Id        string              `bson:"id"`         // Unique identifier for the auth provider
	Name      string              `bson:"name"`       // Human-readable name of the auth provider
	Icon      string              `bson:"icon"`       // Icon URL or identifier for UI display
	Provider  AuthProviderType    `bson:"provider"`   // Type of authentication provider
	Params    []AuthProviderParam `bson:"params"`     // Configuration parameters for the provider
	ProjectId string              `bson:"project_id"` // ID of the project this provider belongs to
	Enabled   bool                `bson:"enabled"`    // Whether the provider is currently active
	CreatedAt *time.Time          `bson:"created_at"` // Timestamp when the provider was created
	UpdatedAt *time.Time          `bson:"updated_at"` // Timestamp when the provider was last updated
	CreatedBy string              `bson:"created_by"` // User who created the provider
	UpdatedBy string              `bson:"updated_by"` // User who last updated the provider
}

// AuthProviderParam represents a configuration parameter for an authentication provider.
// Parameters can include client IDs, secrets, endpoints, and other provider-specific settings.
type AuthProviderParam struct {
	Label    string `bson:"label"`     // Human-readable label for the parameter
	Value    string `bson:"value"`     // Value of the parameter
	Key      string `bson:"key"`       // Unique key identifier for the parameter
	IsSecret bool   `bson:"is_secret"` // Whether this parameter contains sensitive information
}

// AuthProviderModel provides database access patterns and field mappings for AuthProvider entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type AuthProviderModel struct {
	iam                 // Embedded struct providing DbName() method
	IdKey        string // BSON field key for auth provider ID
	NameKey      string // BSON field key for auth provider name
	ProviderKey  string // BSON field key for provider type
	IsEnabledKey string // BSON field key for enabled status
	ProjectIdKey string // BSON field key for project ID
	ParamsKey    string // BSON field key for provider parameters
}

// Name returns the MongoDB collection name for auth providers.
// This implements the DbCollection interface.
func (a AuthProviderModel) Name() string {
	return "auth_providers"
}

// GetAuthProviderModel returns a properly initialized AuthProviderModel with all field mappings.
// This function provides a singleton pattern for accessing auth provider model operations.
//
// Returns an AuthProviderModel instance with all BSON field keys mapped to their respective field names.
func GetAuthProviderModel() AuthProviderModel {
	return AuthProviderModel{
		IdKey:        "id",
		NameKey:      "name",
		ProviderKey:  "provider",
		IsEnabledKey: "is_enabled",
		ProjectIdKey: "project_id",
		ParamsKey:    "params",
	}
}
