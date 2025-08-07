package internal

import (
	"context"
	"errors"
	"strings"

	"github.com/melvinodsa/go-iam/sdk"
)


type ServiceAccountProvider struct {
	// No configuration needed for internal service accounts
}

// NewServiceAccountProvider creates a new service account provider

func NewServiceAccountProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	// We don't need any params for internal service accounts
	// but keeping the signature consistent with other providers
	return &ServiceAccountProvider{}
}

// GetAuthCodeUrl returns empty string as service accounts don't use OAuth flow
func (s *ServiceAccountProvider) GetAuthCodeUrl(state string) string {
	// Service accounts authenticate with client credentials, not OAuth
	return ""
}

// VerifyCode returns error as service accounts don't use authorization codes
func (s *ServiceAccountProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	return nil, errors.New("service accounts use client credentials, not authorization codes")
}

// RefreshToken returns error as service account tokens are not refreshable
func (s *ServiceAccountProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	// Service accounts must re-authenticate with client credentials when token expires
	return nil, errors.New("service account tokens cannot be refreshed - please re-authenticate using client credentials")
}

// GetIdentity extracts the user ID from the service account token
// The token format is "service-account:{userId}"
func (s *ServiceAccountProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	// Validate token format
	if !strings.HasPrefix(token, "service-account:") {
		return nil, errors.New("invalid service account token format")
	}
	
	// Extract user ID
	userId := strings.TrimPrefix(token, "service-account:")
	if userId == "" {
		return nil, errors.New("service account token missing user ID")
	}
	
	// Return identity with the user ID
	// The actual user details will be fetched from the database in GetIdentity
	return []sdk.AuthIdentity{
		{
			Type:     sdk.AuthIdentityTypeEmail, // Using email type for consistency
			Metadata: ServiceAccountIdentity{UserId: userId},
		},
	}, nil
}

// ServiceAccountIdentity implements sdk.AuthMetadataType for service accounts
type ServiceAccountIdentity struct {
	UserId string `json:"user_id"`
}

// UpdateUserDetails sets the user ID for service account users
// The complete user details are fetched from database in the auth service
func (s ServiceAccountIdentity) UpdateUserDetails(user *sdk.User) {
	// For service accounts, we only set the ID
	// The complete user details will be fetched from the database
	user.Id = s.UserId
}