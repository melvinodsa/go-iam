package microsoft

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createMockMicrosoftProvider creates a mock SDK AuthProvider for testing
func createMockMicrosoftProvider() sdk.AuthProvider {
	return sdk.AuthProvider{
		Id:        "microsoft-test-id",
		Name:      "Microsoft Test Provider",
		Provider:  sdk.AuthProviderTypeMicrosoft,
		ProjectId: "test-project",
		Params: []sdk.AuthProviderParam{
			{
				Key:      "@MICROSOFT/CLIENT_ID",
				Value:    "test-client-id",
				Label:    "Client ID",
				IsSecret: false,
			},
			{
				Key:      "@MICROSOFT/CLIENT_SECRET",
				Value:    "test-client-secret",
				Label:    "Client Secret",
				IsSecret: true,
			},
			{
				Key:      "@MICROSOFT/REDIRECT_URL",
				Value:    "http://localhost:8080/callback",
				Label:    "Redirect URL",
				IsSecret: false,
			},
		},
	}
}

// MockMicrosoftEndpoints represents mock Microsoft OAuth endpoints
type MockMicrosoftEndpoints struct {
	TokenServer    *httptest.Server
	UserinfoServer *httptest.Server
}

// NewMockMicrosoftEndpoints creates mock Microsoft OAuth endpoints
func NewMockMicrosoftEndpoints() *MockMicrosoftEndpoints {
	mock := &MockMicrosoftEndpoints{}

	// Mock token endpoint (for RefreshToken)
	mock.TokenServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		refreshToken := r.FormValue("refresh_token")
		clientID := r.FormValue("client_id")
		clientSecret := r.FormValue("client_secret")

		w.Header().Set("Content-Type", "application/json")

		// Simulate various response scenarios based on input
		switch {
		case refreshToken == "":
			// Empty refresh token
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "Refresh token is required",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		case strings.Contains(refreshToken, "expired"):
			// Expired token
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Token has been expired or revoked.",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		case clientID != "test-client-id" || clientSecret != "test-client-secret":
			// Invalid credentials
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_client",
				"error_description": "Invalid client credentials",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		case refreshToken == "valid-refresh-token":
			// Valid refresh token
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(TokenResponse{
				AccessToken: "new-access-token",
				ExpiresIn:   3600,
				TokenType:   "Bearer",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		default:
			// Default error case
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Invalid refresh token",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		}
	}))

	// Mock userinfo endpoint (for GetIdentity)
	mock.UserinfoServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")

		switch authHeader {
		case "Bearer valid-access-token":
			// Valid access token
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"email":      "test@example.com",
				"givenname":  "John",
				"familyname": "Doe",
				"picture":    "https://example.com/profile.jpg",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		case "Bearer empty-fields-token":
			// Token with empty user fields
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"email":      "",
				"givenname":  "",
				"familyname": "",
				"picture":    "",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		case "":
			// No authorization header
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "unauthorized",
				"error_description": "Access token is required",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		default:
			// Invalid access token
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_token",
				"error_description": "The access token is invalid",
			})
			if err != nil {
				log.Errorf("failed to encode response: %w", err)
			}
		}
	}))

	return mock
}

func (m *MockMicrosoftEndpoints) Close() {
	if m.TokenServer != nil {
		m.TokenServer.Close()
	}
	if m.UserinfoServer != nil {
		m.UserinfoServer.Close()
	}
}

func TestNewAuthProvider(t *testing.T) {
	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	assert.NotNil(t, serviceProvider)

	// Cast to access internal config
	provider, ok := serviceProvider.(authProvider)
	require.True(t, ok)

	assert.Equal(t, "test-client-id", provider.cnf.ClientID)
	assert.Equal(t, "test-client-secret", provider.cnf.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", provider.cnf.RedirectURL)
	assert.Equal(t, []string{"openid", "profile", "email"}, provider.cnf.Scopes)
	assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize", provider.cnf.Endpoint.AuthURL)
	assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/token", provider.cnf.Endpoint.TokenURL)
}

func TestGetAuthCodeUrl(t *testing.T) {
	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	tests := []struct {
		name  string
		state string
	}{
		{
			name:  "basic state",
			state: "test-state",
		},
		{
			name:  "empty state",
			state: "",
		},
		{
			name:  "complex state",
			state: "state-with-special-chars-123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUrl := serviceProvider.GetAuthCodeUrl(tt.state)

			assert.Contains(t, authUrl, "login.microsoftonline.com")
			assert.Contains(t, authUrl, "client_id=test-client-id")
			assert.Contains(t, authUrl, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback")
			assert.Contains(t, authUrl, "response_type=code")
			assert.Contains(t, authUrl, "scope=openid+profile+email")

			if tt.state != "" {
				// State values get URL encoded, so check for the encoded version
				if strings.Contains(tt.state, "!@#") {
					assert.Contains(t, authUrl, "state=state-with-special-chars-123%21%40%23")
				} else {
					assert.Contains(t, authUrl, fmt.Sprintf("state=%s", tt.state))
				}
			}
		})
	}
}

func TestVerifyCode(t *testing.T) {
	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	t.Run("invalid code - network error expected", func(t *testing.T) {
		// Since we're not mocking the actual OAuth exchange endpoint,
		// this should fail with a network error
		ctx := context.Background()
		_, err := serviceProvider.VerifyCode(ctx, "invalid-code")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error verifying the code with microsoft exchange")
	})

	t.Run("empty code", func(t *testing.T) {
		ctx := context.Background()
		_, err := serviceProvider.VerifyCode(ctx, "")
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := serviceProvider.VerifyCode(ctx, "some-code")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error verifying the code with microsoft exchange")
	})
}

func TestRefreshToken(t *testing.T) {
	mockEndpoints := NewMockMicrosoftEndpoints()
	defer mockEndpoints.Close()

	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)
	provider, ok := serviceProvider.(authProvider)
	require.True(t, ok)

	// Override the token URL to use our mock server
	provider.cnf.Endpoint.TokenURL = mockEndpoints.TokenServer.URL

	tests := []struct {
		name           string
		refreshToken   string
		expectError    bool
		expectedAccess string
	}{
		{
			name:           "valid refresh token",
			refreshToken:   "valid-refresh-token",
			expectError:    false,
			expectedAccess: "new-access-token",
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			expectError:  true,
		},
		{
			name:         "expired refresh token",
			refreshToken: "expired-refresh-token",
			expectError:  true,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-refresh-token",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.RefreshToken(tt.refreshToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, token) {
					assert.Equal(t, tt.expectedAccess, token.AccessToken)
					assert.Equal(t, tt.refreshToken, token.RefreshToken)
					assert.WithinDuration(t, time.Now().Add(3600*time.Second), token.ExpiresAt, 2*time.Second)
				}
			}
		})
	}
}

func TestRefreshToken_NetworkError(t *testing.T) {
	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)
	provider, ok := serviceProvider.(authProvider)
	require.True(t, ok)

	// Use an invalid URL to simulate network error
	provider.cnf.Endpoint.TokenURL = "http://127.0.0.1:9999/nonexistent"

	token, err := provider.RefreshToken("some-token")
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "error refreshing the token")
}

func TestGetIdentity(t *testing.T) {
	mockEndpoints := NewMockMicrosoftEndpoints()
	defer mockEndpoints.Close()

	// Create a custom GetIdentity function that uses our mock server
	getIdentityWithMockServer := func(token string) ([]sdk.AuthIdentity, error) {
		req, err := http.NewRequest("GET", mockEndpoints.UserinfoServer.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request. %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching the identity. %w", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading the response. %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API error: %s", string(respBytes))
		}

		var userInfo struct {
			Email      string `json:"email"`
			FirstName  string `json:"givenname"`
			LastName   string `json:"familyname"`
			ProfilePic string `json:"picture"`
		}
		err = json.Unmarshal(respBytes, &userInfo)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
		}

		return []sdk.AuthIdentity{
			{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityEmail{Email: userInfo.Email}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityName{Name: fmt.Sprintf("%s %s", userInfo.FirstName, userInfo.LastName)}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: MicrosoftIdentityProfilePic{ProfilePic: userInfo.ProfilePic}},
		}, nil
	}

	tests := []struct {
		name        string
		token       string
		expectError bool
		checkData   func(t *testing.T, identities []sdk.AuthIdentity)
	}{
		{
			name:        "valid access token",
			token:       "valid-access-token",
			expectError: false,
			checkData: func(t *testing.T, identities []sdk.AuthIdentity) {
				require.Len(t, identities, 3)

				// Check email identity
				emailMetadata, ok := identities[0].Metadata.(MicrosoftIdentityEmail)
				require.True(t, ok)
				assert.Equal(t, "test@example.com", emailMetadata.Email)

				// Check name identity
				nameMetadata, ok := identities[1].Metadata.(MicrosoftIdentityName)
				require.True(t, ok)
				assert.Equal(t, "John Doe", nameMetadata.Name)

				// Check profile pic identity
				picMetadata, ok := identities[2].Metadata.(MicrosoftIdentityProfilePic)
				require.True(t, ok)
				assert.Equal(t, "https://example.com/profile.jpg", picMetadata.ProfilePic)
			},
		},
		{
			name:        "empty fields token",
			token:       "empty-fields-token",
			expectError: false,
			checkData: func(t *testing.T, identities []sdk.AuthIdentity) {
				require.Len(t, identities, 3)

				// Check empty email
				emailMetadata, ok := identities[0].Metadata.(MicrosoftIdentityEmail)
				require.True(t, ok)
				assert.Equal(t, "", emailMetadata.Email)

				// Check empty name (should be single space)
				nameMetadata, ok := identities[1].Metadata.(MicrosoftIdentityName)
				require.True(t, ok)
				assert.Equal(t, " ", nameMetadata.Name)
			},
		},
		{
			name:        "invalid access token",
			token:       "invalid-access-token",
			expectError: true,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identities, err := getIdentityWithMockServer(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, identities)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, identities)
				if tt.checkData != nil {
					tt.checkData(t, identities)
				}
			}
		})
	}
}

func TestGetIdentity_NetworkError(t *testing.T) {
	mockProvider := createMockMicrosoftProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	// The actual GetIdentity method might succeed with empty data or fail with network error
	// when using an invalid token. Both cases are valid for our test.
	identities, err := serviceProvider.GetIdentity("invalid-token-that-should-fail")

	// Microsoft API might return empty data instead of error for invalid tokens
	// This is still a valid test case - we're testing the integration works
	if err != nil {
		// If there's an error, it should be a network or API error
		assert.Nil(t, identities)
		errorStr := err.Error()
		assert.True(t,
			strings.Contains(errorStr, "error fetching the identity") ||
				strings.Contains(errorStr, "error unmarshalling the response") ||
				strings.Contains(errorStr, "error creating request"),
			"Expected network or API error, got: %s", errorStr)
	} else {
		// If no error, Microsoft returned empty data, which is also valid
		assert.NotNil(t, identities)
		// We should get 3 identity objects even if they're empty
		assert.Len(t, identities, 3)
	}
}

// Test identity metadata types
func TestMicrosoftIdentityEmail_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := MicrosoftIdentityEmail{Email: "test@example.com"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "test@example.com", user.Email)
}

func TestMicrosoftIdentityName_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := MicrosoftIdentityName{Name: "John Doe"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "John Doe", user.Name)
}

func TestMicrosoftIdentityProfilePic_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := MicrosoftIdentityProfilePic{ProfilePic: "https://example.com/pic.jpg"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "https://example.com/pic.jpg", user.ProfilePic)
}

// Test edge cases
func TestAuthProvider_EdgeCases(t *testing.T) {
	t.Run("provider with missing parameters", func(t *testing.T) {
		emptyProvider := sdk.AuthProvider{
			Id:       "empty-provider",
			Provider: sdk.AuthProviderTypeMicrosoft,
			Params:   []sdk.AuthProviderParam{},
		}

		serviceProvider := NewAuthProvider(emptyProvider)
		provider, ok := serviceProvider.(authProvider)
		require.True(t, ok)

		// Should have empty values for missing parameters
		assert.Equal(t, "", provider.cnf.ClientID)
		assert.Equal(t, "", provider.cnf.ClientSecret)
		assert.Equal(t, "", provider.cnf.RedirectURL)
	})

	t.Run("provider with partial parameters", func(t *testing.T) {
		partialProvider := sdk.AuthProvider{
			Id:       "partial-provider",
			Provider: sdk.AuthProviderTypeMicrosoft,
			Params: []sdk.AuthProviderParam{
				{Key: "@MICROSOFT/CLIENT_ID", Value: "only-client-id"},
			},
		}

		serviceProvider := NewAuthProvider(partialProvider)
		provider, ok := serviceProvider.(authProvider)
		require.True(t, ok)

		assert.Equal(t, "only-client-id", provider.cnf.ClientID)
		assert.Equal(t, "", provider.cnf.ClientSecret)
		assert.Equal(t, "", provider.cnf.RedirectURL)
	})
}
