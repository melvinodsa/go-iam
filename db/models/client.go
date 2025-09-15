package models

import "time"

// Client represents an OAuth2/OIDC client application in the Go IAM system.
// Clients are applications that can authenticate users and access protected resources.
// Each client belongs to a project and can have various configuration options.
type Client struct {
	Id                    string     `bson:"id"`                       // Unique identifier for the client
	Name                  string     `bson:"name"`                     // Human-readable name of the client
	Description           string     `bson:"description"`              // Detailed description of the client's purpose
	Secret                string     `bson:"secret"`                   // Client secret for authentication
	Tags                  []string   `bson:"tags"`                     // Tags for categorizing and filtering clients
	RedirectURLs          []string   `bson:"redirect_urls"`            // Allowed redirect URLs for OAuth2 flows
	DefaultAuthProviderId string     `bson:"default_auth_provider_id"` // Default authentication provider for this client
	GoIamClient           bool       `bson:"go_iam_client"`            // Indicates if this is a Go-IAM system client
	ProjectId             string     `bson:"project_id"`               // ID of the project this client belongs to
	ServiceAccountEmail   string     `bson:"service_account_email"`    // Email for service account authentication
	Scopes                []string   `bson:"scopes"`                   // OAuth2 scopes this client can request
	Enabled               bool       `bson:"enabled"`                  // Whether the client is currently active
	LinkedUserId          string     `bson:"linked_user_id"`           // User ID for service account clients
	CreatedAt             *time.Time `bson:"created_at"`               // Timestamp when the client was created
	CreatedBy             string     `bson:"created_by"`               // User who created the client
	UpdatedAt             *time.Time `bson:"updated_at"`               // Timestamp when the client was last updated
	UpdatedBy             string     `bson:"updated_by"`               // User who last updated the client
}

// ClientModel provides database access patterns and field mappings for Client entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type ClientModel struct {
	iam                    // Embedded struct providing DbName() method
	IdKey           string // BSON field key for client ID
	NameKey         string // BSON field key for client name
	TagsKey         string // BSON field key for client tags
	DescriptionKey  string // BSON field key for client description
	ProjectIdKey    string // BSON field key for project ID
	GoIamClientKey  string // BSON field key for Go-IAM client flag
	LinkedUserIdKey string // BSON field key for linked user ID (service accounts)
	UpdatedAtKey    string // BSON field key for last updated timestamp
}

// Name returns the MongoDB collection name for clients.
// This implements the DbCollection interface.
func (c ClientModel) Name() string {
	return "clients"
}

// GetClientModel returns a properly initialized ClientModel with all field mappings.
// This function provides a singleton pattern for accessing client model operations.
//
// Returns a ClientModel instance with all BSON field keys mapped to their respective field names.
func GetClientModel() ClientModel {
	return ClientModel{
		IdKey:           "id",
		NameKey:         "name",
		TagsKey:         "tags",
		DescriptionKey:  "description",
		ProjectIdKey:    "project_id",
		GoIamClientKey:  "go_iam_client",
		LinkedUserIdKey: "linked_user_id",
		UpdatedAtKey:    "updated_at",
	}
}
