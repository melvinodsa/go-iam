package goiamclient

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/melvinodsa/go-iam/sdk"
)

const (
	// Token validity periods
	AccessTokenDuration  = time.Hour * 24     
	RefreshTokenDuration = time.Hour * 24 * 30 // 30 days refresh token
)

type GoIAMClientProvider struct {
	// No configuration needed as this is internal
}

// NewGoIAMClientProvider creates a new GoIAM client provider
func NewGoIAMClientProvider(p sdk.AuthProvider) sdk.ServiceProvider {
	return &GoIAMClientProvider{}
}

// GetAuthCodeUrl - not used for service accounts but required by interface
func (g *GoIAMClientProvider) GetAuthCodeUrl(state string) string {
	// Service accounts don't use OAuth flow, they use client credentials
	return ""
}

// VerifyCode handles the "authorization code" which for service accounts
// is actually the client credentials validation result
func (g *GoIAMClientProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	// For service accounts, the "code" contains the validated client info
	// This is called after client credentials are verified
	
	// Decode the synthetic code (contains client ID and user ID)
	var codeData ServiceAccountCode
	if err := json.Unmarshal([]byte(code), &codeData); err != nil {
		return nil, fmt.Errorf("invalid service account code: %w", err)
	}
	
	// Generate tokens
	accessToken := g.generateAccessToken(codeData.ClientID, codeData.UserID)
	refreshToken := g.generateRefreshToken(codeData.ClientID, codeData.UserID)
	
	return &sdk.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(AccessTokenDuration),
	}, nil
}

// RefreshToken handles refresh token for service accounts
func (g *GoIAMClientProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	tokenData, err := g.parseServiceAccountToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	if tokenData.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}
	
	newAccessToken := g.generateAccessToken(tokenData.ClientID, tokenData.UserID)
	
	return &sdk.AuthToken{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(AccessTokenDuration),
	}, nil
}

// GetIdentity extracts identity from the service account token
func (g *GoIAMClientProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	tokenData, err := g.parseServiceAccountToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid service account token: %w", err)
	}
	
	return []sdk.AuthIdentity{
		{
			Type:     sdk.AuthIdentityType("goiam_client"),
			Metadata: GoIAMClientIdentity{
				ClientID: tokenData.ClientID,
				UserID:   tokenData.UserID,
			},
		},
	}, nil
}

type ServiceAccountCode struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
}

type ServiceAccountTokenData struct {
	ClientID  string    `json:"client_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

type GoIAMClientIdentity struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
}

func (g GoIAMClientIdentity) UpdateUserDetails(user *sdk.User) {
	// The actual user will be fetched using the UserID
	user.Id = g.UserID
}

func (g *GoIAMClientProvider) generateAccessToken(clientID, userID string) string {
	tokenData := ServiceAccountTokenData{
		ClientID:  clientID,
		UserID:    userID,
		Type:      "access",
		ExpiresAt: time.Now().Add(AccessTokenDuration),
		IssuedAt:  time.Now(),
	}
	
	data, _ := json.Marshal(tokenData)
	return base64.RawURLEncoding.EncodeToString(data)
}

func (g *GoIAMClientProvider) generateRefreshToken(clientID, userID string) string {
	tokenData := ServiceAccountTokenData{
		ClientID:  clientID,
		UserID:    userID,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(RefreshTokenDuration),
		IssuedAt:  time.Now(),
	}
	
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	
	data, _ := json.Marshal(tokenData)
	combined := append(data, randomBytes...)
	return base64.RawURLEncoding.EncodeToString(combined)
}

func (g *GoIAMClientProvider) parseServiceAccountToken(token string) (*ServiceAccountTokenData, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}
	
	// For refresh tokens, remove the random bytes (last 16 bytes)
	if len(data) > 16 {
		data = data[:len(data)-16]
	}
	
	var tokenData ServiceAccountTokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return nil, fmt.Errorf("failed to parse token data: %w", err)
	}
	
	return &tokenData, nil
}