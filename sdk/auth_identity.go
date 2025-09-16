package sdk

// AuthIdentityType represents the type of authentication identity.
// This determines how a user is identified during authentication.
type AuthIdentityType string

const (
	// AuthIdentityTypeEmail indicates email-based authentication identity.
	AuthIdentityTypeEmail AuthIdentityType = "email"

	// AuthIdentityTypePhone indicates phone-based authentication identity.
	AuthIdentityTypePhone AuthIdentityType = "phone"
)

// AuthMetadataType is an interface for authentication metadata that can update user details.
// Implementations of this interface provide type-specific logic for updating
// user information based on authentication provider data.
type AuthMetadataType interface {
	// UpdateUserDetails updates the provided user with authentication-specific information.
	UpdateUserDetails(user *User)
}

// AuthIdentity represents a user's authentication identity with associated metadata.
// This structure ties together the identity type with provider-specific metadata
// that can be used to update user information during authentication.
type AuthIdentity struct {
	Type     AuthIdentityType `json:"type"`     // The type of authentication identity
	Metadata AuthMetadataType `json:"metadata"` // Provider-specific metadata for this identity
}

// UpdateUserDetails updates the provided user with information from this auth identity.
// This method delegates to the metadata's UpdateUserDetails method to perform
// type-specific user data updates.
func (a AuthIdentity) UpdateUserDetails(user *User) {
	a.Metadata.UpdateUserDetails(user)
}
