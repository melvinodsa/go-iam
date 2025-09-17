package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider sdk.AuthProvider
		expected authProvider
	}{
		{
			name: "creates provider with default scopes",
			provider: mockAuthProvider(map[string]string{
				"@OIDC/CLIENT_ID":         "test-client-id",
				"@OIDC/CLIENT_SECRET":     "test-client-secret",
				"@OIDC/REDIRECT_URL":      "https://example.com/callback",
				"@OIDC/AUTHORIZATION_URL": "https://provider.example.com/auth",
				"@OIDC/TOKEN_URL":         "https://provider.example.com/token",
				"@OIDC/USERINFO_URL":      "https://provider.example.com/userinfo",
				"@OIDC/ISSUER":            "https://provider.example.com",
			}),
			expected: authProvider{
				userInfoURL:  "https://provider.example.com/userinfo",
				issuer:       "https://provider.example.com",
				providerName: "Test Provider",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewAuthProvider(tt.provider)
			assert.NotNil(t, provider)

			oidcProvider, ok := provider.(authProvider)
			require.True(t, ok)

			assert.Equal(t, "test-client-id", oidcProvider.cnf.ClientID)
			assert.Equal(t, "test-client-secret", oidcProvider.cnf.ClientSecret)
			assert.Equal(t, "https://example.com/callback", oidcProvider.cnf.RedirectURL)
			assert.Equal(t, "https://provider.example.com/auth", oidcProvider.cnf.Endpoint.AuthURL)
			assert.Equal(t, "https://provider.example.com/token", oidcProvider.cnf.Endpoint.TokenURL)

			if tt.provider.GetParam("@OIDC/SCOPES") != "" {
				assert.Equal(t, []string{"openid", "profile", "email", "custom_scope"}, oidcProvider.cnf.Scopes)
			} else {
				assert.Equal(t, []string{"openid", "profile", "email"}, oidcProvider.cnf.Scopes)
			}
		})
	}
}

func TestAuthProvider_HasRefreshTokenFlow(t *testing.T) {
	provider := createTestProvider()
	assert.True(t, provider.HasRefreshTokenFlow())
}

func TestAuthProvider_GetAuthCodeUrl(t *testing.T) {
	provider := createTestProvider()
	state := "test-state"

	authURL := provider.GetAuthCodeUrl(state)

	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	assert.Equal(t, "provider.example.com", parsedURL.Host)
	assert.Equal(t, "/auth", parsedURL.Path)

	query := parsedURL.Query()
	assert.Equal(t, "test-client-id", query.Get("client_id"))
	assert.Equal(t, state, query.Get("state"))
	assert.Equal(t, "code", query.Get("response_type"))
}

func TestAuthProvider_VerifyCode(t *testing.T) {
	// Mock token server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			t.Errorf("Expected path '/token', got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"access_token":  "access-token-123",
			"refresh_token": "refresh-token-456",
			"expires_in":    3600,
			"token_type":    "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	provider := createTestProviderWithServer(server.URL)

	token, err := provider.VerifyCode(context.Background(), "test-code")

	require.NoError(t, err)
	assert.Equal(t, "access-token-123", token.AccessToken)
	assert.Equal(t, "refresh-token-456", token.RefreshToken)
	assert.True(t, token.ExpiresAt.After(time.Now()))
}

func TestAuthProvider_RefreshToken(t *testing.T) {
	// Mock token server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			t.Errorf("Expected path '/token', got %s", r.URL.Path)
		}

		err := r.ParseForm()
		require.NoError(t, err)

		assert.Equal(t, "refresh_token", r.Form.Get("grant_type"))
		assert.Equal(t, "old-refresh-token", r.Form.Get("refresh_token"))

		response := TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}
		err = json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	provider := createTestProviderWithServer(server.URL)

	token, err := provider.RefreshToken("old-refresh-token")

	require.NoError(t, err)
	assert.Equal(t, "new-access-token", token.AccessToken)
	assert.Equal(t, "new-refresh-token", token.RefreshToken)
	assert.True(t, token.ExpiresAt.After(time.Now()))
}

func TestAuthProvider_GetIdentity(t *testing.T) {
	// Mock userinfo server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/userinfo" {
			t.Errorf("Expected path '/userinfo', got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer test-access-token", auth)

		userInfo := UserInfoResponse{
			Sub:           "user-123",
			Name:          "John Doe",
			GivenName:     "John",
			FamilyName:    "Doe",
			Email:         "john.doe@example.com",
			EmailVerified: true,
			Picture:       "https://example.com/avatar.jpg",
		}
		err := json.NewEncoder(w).Encode(userInfo)
		require.NoError(t, err)
	}))
	defer server.Close()

	provider := createTestProviderWithUserInfoServer(server.URL)

	identities, err := provider.GetIdentity("test-access-token")

	require.NoError(t, err)
	assert.Len(t, identities, 3)

	// Check email identity
	emailIdentity := identities[0]
	assert.Equal(t, sdk.AuthIdentityTypeEmail, emailIdentity.Type)
	emailMeta, ok := emailIdentity.Metadata.(OIDCIdentityEmail)
	require.True(t, ok)
	assert.Equal(t, "john.doe@example.com", emailMeta.Email)

	// Check name identity
	nameIdentity := identities[1]
	assert.Equal(t, sdk.AuthIdentityTypeEmail, nameIdentity.Type)
	nameMeta, ok := nameIdentity.Metadata.(OIDCIdentityName)
	require.True(t, ok)
	assert.Equal(t, "John Doe", nameMeta.Name)

	// Check profile picture identity
	picIdentity := identities[2]
	assert.Equal(t, sdk.AuthIdentityTypeEmail, picIdentity.Type)
	picMeta, ok := picIdentity.Metadata.(OIDCIdentityProfilePic)
	require.True(t, ok)
	assert.Equal(t, "https://example.com/avatar.jpg", picMeta.ProfilePic)
}

func TestIdentityTypes_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}

	// Test email identity
	emailIdentity := OIDCIdentityEmail{Email: "test@example.com"}
	emailIdentity.UpdateUserDetails(user)
	assert.Equal(t, "test@example.com", user.Email)

	// Test name identity
	nameIdentity := OIDCIdentityName{Name: "Test User"}
	nameIdentity.UpdateUserDetails(user)
	assert.Equal(t, "Test User", user.Name)

	// Test profile picture identity
	picIdentity := OIDCIdentityProfilePic{ProfilePic: "https://example.com/pic.jpg"}
	picIdentity.UpdateUserDetails(user)
	assert.Equal(t, "https://example.com/pic.jpg", user.ProfilePic)
}

// Helper functions

func mockAuthProvider(params map[string]string) sdk.AuthProvider {
	return sdk.AuthProvider{
		Name: "Test Provider",
		Params: func() []sdk.AuthProviderParam {
			var paramsList []sdk.AuthProviderParam
			for key, value := range params {
				paramsList = append(paramsList, sdk.AuthProviderParam{
					Key:   key,
					Value: value,
				})
			}
			return paramsList
		}(),
	}
}

func createTestProvider() authProvider {
	provider := mockAuthProvider(map[string]string{
		"@OIDC/CLIENT_ID":         "test-client-id",
		"@OIDC/CLIENT_SECRET":     "test-client-secret",
		"@OIDC/REDIRECT_URL":      "https://example.com/callback",
		"@OIDC/AUTHORIZATION_URL": "https://provider.example.com/auth",
		"@OIDC/TOKEN_URL":         "https://provider.example.com/token",
		"@OIDC/USERINFO_URL":      "https://provider.example.com/userinfo",
	})
	return NewAuthProvider(provider).(authProvider)
}

func createTestProviderWithServer(serverURL string) authProvider {
	provider := mockAuthProvider(map[string]string{
		"@OIDC/CLIENT_ID":         "test-client-id",
		"@OIDC/CLIENT_SECRET":     "test-client-secret",
		"@OIDC/REDIRECT_URL":      "https://example.com/callback",
		"@OIDC/AUTHORIZATION_URL": serverURL + "/auth",
		"@OIDC/TOKEN_URL":         serverURL + "/token",
		"@OIDC/USERINFO_URL":      serverURL + "/userinfo",
	})
	return NewAuthProvider(provider).(authProvider)
}

func createTestProviderWithUserInfoServer(serverURL string) authProvider {
	provider := mockAuthProvider(map[string]string{
		"@OIDC/CLIENT_ID":         "test-client-id",
		"@OIDC/CLIENT_SECRET":     "test-client-secret",
		"@OIDC/REDIRECT_URL":      "https://example.com/callback",
		"@OIDC/AUTHORIZATION_URL": "https://provider.example.com/auth",
		"@OIDC/TOKEN_URL":         "https://provider.example.com/token",
		"@OIDC/USERINFO_URL":      serverURL + "/userinfo",
	})
	return NewAuthProvider(provider).(authProvider)
}
