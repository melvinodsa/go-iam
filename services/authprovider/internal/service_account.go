package internal

import (
	"context"
	"errors"
	"strings"
	"fmt"

	"github.com/melvinodsa/go-iam/sdk"
)


type ServiceAccountProvider struct {
	clientId string
}

// NewServiceAccountProvider creates a new service account provider

func NewServiceAccountProvider(p sdk.AuthProvider) sdk.ServiceProvider {

	return &ServiceAccountProvider{
		clientId : p.GetParam("@INTERNAL/CLIENT_ID"),
	}
}

// GetAuthCodeUrl - not used for service accounts
func (s *ServiceAccountProvider) GetAuthCodeUrl(state string) string {
	return ""
}

// VerifyCode - not used for service accounts
func (s *ServiceAccountProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	return nil, errors.New("service accounts use client credentials, not authorization codes")
}

// RefreshToken - service accounts don't refresh, they re-authenticate
func (s *ServiceAccountProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	return nil, errors.New("service account tokens cannot be refreshed - please re-authenticate using client credentials")
}

// GetIdentity extracts the client ID from the service account token
func (s *ServiceAccountProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	// Validate token format
	if !strings.HasPrefix(token, "service-account:") {
		return nil, fmt.Errorf("invalid service account token format")
	}
	
	clientId := strings.TrimPrefix(token, "service-account:")
	if clientId == "" {
		return nil, fmt.Errorf("service account token missing client ID")
	}
	
	// Return identity that will be used to fetch/create user
	return []sdk.AuthIdentity{
		{
			Type:     sdk.AuthIdentityType("service_account"),
			Metadata: ServiceAccountIdentity{ClientId: clientId},
		},
	}, nil
	
}

// ServiceAccountIdentity implements sdk.AuthMetadataType for service accounts
type ServiceAccountIdentity struct {
	ClientId string `json:"client_id"`
}

func (s ServiceAccountIdentity) UpdateUserDetails(user *sdk.User) {
	// Set a temporary identifier that will be used in getOrCreateUser
	// The actual user will be fetched using the client's linked_user_id
	user.Email = fmt.Sprintf("service-account-%s@internal", s.ClientId)
}