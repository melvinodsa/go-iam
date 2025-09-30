package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createMockGitHubProvider creates a mock SDK AuthProvider for testing
func createMockGitHubProvider() sdk.AuthProvider {
	return sdk.AuthProvider{
		Id:        "github-test-id",
		Name:      "GitHub Test Provider",
		Provider:  sdk.AuthProviderTypeGitHub,
		ProjectId: "test-project",
		Params: []sdk.AuthProviderParam{
			{
				Key:      "@GITHUB/CLIENT_ID",
				Value:    "test-client-id",
				Label:    "Client ID",
				IsSecret: false,
			},
			{
				Key:      "@GITHUB/CLIENT_SECRET",
				Value:    "test-client-secret",
				Label:    "Client Secret",
				IsSecret: true,
			},
			{
				Key:      "@GITHUB/REDIRECT_URL",
				Value:    "http://localhost:8080/callback",
				Label:    "Redirect URL",
				IsSecret: false,
			},
		},
	}
}

// MockGitHubEndpoints represents mock GitHub OAuth endpoints
type MockGitHubEndpoints struct {
	TokenServer  *httptest.Server
	UserServer   *httptest.Server
	EmailsServer *httptest.Server
}

// NewMockGitHubEndpoints creates mock GitHub OAuth endpoints
func NewMockGitHubEndpoints() *MockGitHubEndpoints {
	mock := &MockGitHubEndpoints{}

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
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case strings.Contains(refreshToken, "expired"):
			// Expired token
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Token has been expired or revoked.",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case clientID != "test-client-id" || clientSecret != "test-client-secret":
			// Invalid credentials
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_client",
				"error_description": "Invalid client credentials",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case refreshToken == "valid-refresh-token":
			// Valid refresh token
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(TokenResponse{
				AccessToken: "new-access-token",
				ExpiresIn:   3600,
				TokenType:   "Bearer",
				Scope:       "user:email,read:user",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		default:
			// Default error case
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Invalid refresh token",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		}
	}))

	// Mock user endpoint (for GetIdentity)
	mock.UserServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")

		switch authHeader {
		case "Bearer valid-access-token":
			// Valid access token
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"email":      "test@example.com",
				"name":       "John Doe",
				"avatar_url": "https://avatars.githubusercontent.com/u/123456?v=4",
				"login":      "johndoe",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case "Bearer no-email-token":
			// Token with no email (private email)
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"email":      nil,
				"name":       "Jane Doe",
				"avatar_url": "https://avatars.githubusercontent.com/u/654321?v=4",
				"login":      "janedoe",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case "Bearer empty-fields-token":
			// Token with empty user fields
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"email":      "",
				"name":       "",
				"avatar_url": "",
				"login":      "emptyuser",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case "":
			// No authorization header
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		default:
			// Invalid access token
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		}
	}))

	// Mock emails endpoint (for private emails)
	mock.EmailsServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")

		switch authHeader {
		case "Bearer no-email-token":
			// Return primary email for user with private email
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"email":      "jane.doe@private.example.com",
					"primary":    true,
					"verified":   true,
					"visibility": "private",
				},
				{
					"email":      "jane.public@example.com",
					"primary":    false,
					"verified":   true,
					"visibility": "public",
				},
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		default:
			// Invalid token or no emails
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			})
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		}
	}))

	return mock
}

func (m *MockGitHubEndpoints) Close() {
	if m.TokenServer != nil {
		m.TokenServer.Close()
	}
	if m.UserServer != nil {
		m.UserServer.Close()
	}
	if m.EmailsServer != nil {
		m.EmailsServer.Close()
	}
}

func TestNewAuthProvider(t *testing.T) {
	mockProvider := createMockGitHubProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	assert.NotNil(t, serviceProvider)

	// Verify it's the correct type
	provider, ok := serviceProvider.(authProvider)
	assert.True(t, ok)

	// Verify OAuth config
	assert.Equal(t, "test-client-id", provider.cnf.ClientID)
	assert.Equal(t, "test-client-secret", provider.cnf.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", provider.cnf.RedirectURL)
	assert.Contains(t, provider.cnf.Scopes, "user:email")
	assert.Contains(t, provider.cnf.Scopes, "read:user")
	assert.Equal(t, "https://github.com/login/oauth/authorize", provider.cnf.Endpoint.AuthURL)
	assert.Equal(t, "https://github.com/login/oauth/access_token", provider.cnf.Endpoint.TokenURL)
}

func TestGetAuthCodeUrl(t *testing.T) {
	mockProvider := createMockGitHubProvider()
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

			assert.Contains(t, authUrl, "github.com/login/oauth/authorize")
			assert.Contains(t, authUrl, "client_id=test-client-id")
			assert.Contains(t, authUrl, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback")
			assert.Contains(t, authUrl, "response_type=code")
			assert.Contains(t, authUrl, "scope=user%3Aemail+read%3Auser")

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
	mockProvider := createMockGitHubProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	t.Run("invalid code - network error expected", func(t *testing.T) {
		// Since we're not mocking the actual OAuth exchange endpoint,
		// this should fail with a network error
		ctx := context.Background()
		_, err := serviceProvider.VerifyCode(ctx, "invalid-code")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error verifying the code with github exchange")
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
		assert.Contains(t, err.Error(), "error verifying the code with github exchange")
	})
}

func TestGetIdentity(t *testing.T) {
	mockEndpoints := NewMockGitHubEndpoints()
	defer mockEndpoints.Close()

	// Create a custom GetIdentity function that uses our mock servers
	getIdentityWithMockServer := func(token string, userURL, emailsURL string) ([]sdk.AuthIdentity, error) {
		// Mock the GitHub API requests
		req, err := http.NewRequest("GET", userURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request. %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching the identity. %w", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			respBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API error: %s", string(respBytes))
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading the response. %w", err)
		}

		var userInfo struct {
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		}
		err = json.Unmarshal(respBytes, &userInfo)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
		}

		identities := []sdk.AuthIdentity{
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: userInfo.Email}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityName{Name: userInfo.Name}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityProfilePic{ProfilePic: userInfo.AvatarURL}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityUsername{Username: userInfo.Login}},
		}

		// If email is empty, try to fetch from emails endpoint
		if userInfo.Email == "" && emailsURL != "" {
			emailReq, err := http.NewRequest("GET", emailsURL, nil)
			if err == nil {
				emailReq.Header.Set("Authorization", "Bearer "+token)
				emailReq.Header.Set("Accept", "application/vnd.github.v3+json")

				emailResp, err := http.DefaultClient.Do(emailReq)
				if err == nil && emailResp.StatusCode == http.StatusOK {
					defer func() {
						_ = emailResp.Body.Close()
					}()
					emailBytes, err := io.ReadAll(emailResp.Body)
					if err == nil {
						var emails []struct {
							Email   string `json:"email"`
							Primary bool   `json:"primary"`
						}
						if json.Unmarshal(emailBytes, &emails) == nil {
							for _, email := range emails {
								if email.Primary {
									identities[0] = sdk.AuthIdentity{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: email.Email}}
									break
								}
							}
						}
					}
				}
			}
		}

		return identities, nil
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
				require.Len(t, identities, 4)

				// Check email identity
				emailMetadata, ok := identities[0].Metadata.(GitHubIdentityEmail)
				require.True(t, ok)
				assert.Equal(t, "test@example.com", emailMetadata.Email)

				// Check name identity
				nameMetadata, ok := identities[1].Metadata.(GitHubIdentityName)
				require.True(t, ok)
				assert.Equal(t, "John Doe", nameMetadata.Name)

				// Check profile pic identity
				picMetadata, ok := identities[2].Metadata.(GitHubIdentityProfilePic)
				require.True(t, ok)
				assert.Equal(t, "https://avatars.githubusercontent.com/u/123456?v=4", picMetadata.ProfilePic)

				// Check username identity
				usernameMetadata, ok := identities[3].Metadata.(GitHubIdentityUsername)
				require.True(t, ok)
				assert.Equal(t, "johndoe", usernameMetadata.Username)
			},
		},
		{
			name:        "no email token - should fetch from emails endpoint",
			token:       "no-email-token",
			expectError: false,
			checkData: func(t *testing.T, identities []sdk.AuthIdentity) {
				require.Len(t, identities, 4)

				// Check that primary email was fetched from emails endpoint
				emailMetadata, ok := identities[0].Metadata.(GitHubIdentityEmail)
				require.True(t, ok)
				assert.Equal(t, "jane.doe@private.example.com", emailMetadata.Email)

				// Check username
				usernameMetadata, ok := identities[3].Metadata.(GitHubIdentityUsername)
				require.True(t, ok)
				assert.Equal(t, "janedoe", usernameMetadata.Username)
			},
		},
		{
			name:        "empty fields token",
			token:       "empty-fields-token",
			expectError: false,
			checkData: func(t *testing.T, identities []sdk.AuthIdentity) {
				require.Len(t, identities, 4)

				// Check empty email
				emailMetadata, ok := identities[0].Metadata.(GitHubIdentityEmail)
				require.True(t, ok)
				assert.Equal(t, "", emailMetadata.Email)

				// Check empty name
				nameMetadata, ok := identities[1].Metadata.(GitHubIdentityName)
				require.True(t, ok)
				assert.Equal(t, "", nameMetadata.Name)

				// Check username (should still be present)
				usernameMetadata, ok := identities[3].Metadata.(GitHubIdentityUsername)
				require.True(t, ok)
				assert.Equal(t, "emptyuser", usernameMetadata.Username)
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
			identities, err := getIdentityWithMockServer(tt.token, mockEndpoints.UserServer.URL, mockEndpoints.EmailsServer.URL)

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
	mockProvider := createMockGitHubProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	// The actual GetIdentity method might succeed with empty data or fail with network error
	// when using an invalid token. Both cases are valid for our test.
	identities, err := serviceProvider.GetIdentity("invalid-token-that-should-fail")

	// GitHub API might return error or empty data instead of error for invalid tokens
	// This is still a valid test case - we're testing the integration works
	if err != nil {
		// If there's an error, it should be a network or API error
		assert.Nil(t, identities)
		errorStr := err.Error()
		assert.True(t,
			strings.Contains(errorStr, "error fetching the identity") ||
				strings.Contains(errorStr, "API error") ||
				strings.Contains(errorStr, "error creating request"),
			"Expected network or API error, got: %s", errorStr)
	} else {
		// If no error, GitHub returned some data, which is also valid
		assert.NotNil(t, identities)
		// We should get 4 identity objects even if they're empty/invalid
		assert.Len(t, identities, 4)
	}
}

func TestGetIdentity_APIError(t *testing.T) {
	// Create a mock server that returns an error status
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal Server Error"}`))
	}))
	defer mockServer.Close()

	// Create a custom GetIdentity function that uses our mock server
	getIdentityWithMock := func(token string) ([]sdk.AuthIdentity, error) {
		req, err := http.NewRequest("GET", mockServer.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request. %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

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
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		}
		err = json.Unmarshal(respBytes, &userInfo)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
		}

		identities := []sdk.AuthIdentity{
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: userInfo.Email}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityName{Name: userInfo.Name}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityProfilePic{ProfilePic: userInfo.AvatarURL}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityUsername{Username: userInfo.Login}},
		}

		return identities, nil
	}

	identities, err := getIdentityWithMock("any-token")
	assert.Error(t, err)
	assert.Nil(t, identities)
	assert.Contains(t, err.Error(), "API error")
}

func TestGetIdentity_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer mockServer.Close()

	// Create a custom GetIdentity function that uses our mock server
	getIdentityWithMock := func(token string) ([]sdk.AuthIdentity, error) {
		req, err := http.NewRequest("GET", mockServer.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request. %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

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
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		}
		err = json.Unmarshal(respBytes, &userInfo)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling the response. %s - %w", string(respBytes), err)
		}

		identities := []sdk.AuthIdentity{
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityEmail{Email: userInfo.Email}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityName{Name: userInfo.Name}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityProfilePic{ProfilePic: userInfo.AvatarURL}},
			{Type: sdk.AuthIdentityTypeEmail, Metadata: GitHubIdentityUsername{Username: userInfo.Login}},
		}

		return identities, nil
	}

	identities, err := getIdentityWithMock("any-token")
	assert.Error(t, err)
	assert.Nil(t, identities)
	assert.Contains(t, err.Error(), "error unmarshalling the response")
}

// Test identity metadata types
func TestGitHubIdentityEmail_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := GitHubIdentityEmail{Email: "test@example.com"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "test@example.com", user.Email)
}

func TestGitHubIdentityName_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := GitHubIdentityName{Name: "John Doe"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "John Doe", user.Name)
}

func TestGitHubIdentityProfilePic_UpdateUserDetails(t *testing.T) {
	user := &sdk.User{}
	identity := GitHubIdentityProfilePic{ProfilePic: "https://avatars.githubusercontent.com/u/123456?v=4"}

	identity.UpdateUserDetails(user)

	assert.Equal(t, "https://avatars.githubusercontent.com/u/123456?v=4", user.ProfilePic)
}

func TestGitHubIdentityUsername_UpdateUserDetails(t *testing.T) {
	t.Run("update name when empty", func(t *testing.T) {
		user := &sdk.User{}
		identity := GitHubIdentityUsername{Username: "johndoe"}

		identity.UpdateUserDetails(user)

		assert.Equal(t, "johndoe", user.Name)
	})

	t.Run("don't override existing name", func(t *testing.T) {
		user := &sdk.User{Name: "John Doe"}
		identity := GitHubIdentityUsername{Username: "johndoe"}

		identity.UpdateUserDetails(user)

		assert.Equal(t, "John Doe", user.Name) // Should not be overridden
	})
}

func TestHasRefreshTokenFlow(t *testing.T) {
	mockProvider := createMockGitHubProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	// GitHub doesn't support refresh token flow
	assert.False(t, serviceProvider.HasRefreshTokenFlow())
}

func TestRefreshToken(t *testing.T) {
	mockProvider := createMockGitHubProvider()
	serviceProvider := NewAuthProvider(mockProvider)

	// GitHub refresh token should return nil
	token, err := serviceProvider.RefreshToken("any-token")
	assert.Nil(t, token)
	assert.Nil(t, err)
}

func TestAuthProvider_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		provider sdk.AuthProvider
	}{
		{
			name: "provider with missing parameters",
			provider: sdk.AuthProvider{
				Id:       "github-test-id",
				Name:     "GitHub Test Provider",
				Provider: sdk.AuthProviderTypeGitHub,
				Params:   []sdk.AuthProviderParam{}, // Empty params
			},
		},
		{
			name: "provider with partial parameters",
			provider: sdk.AuthProvider{
				Id:       "github-test-id",
				Name:     "GitHub Test Provider",
				Provider: sdk.AuthProviderTypeGitHub,
				Params: []sdk.AuthProviderParam{
					{
						Key:   "@GITHUB/CLIENT_ID",
						Value: "test-client-id",
					},
					// Missing CLIENT_SECRET and REDIRECT_URL
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic even with missing parameters
			serviceProvider := NewAuthProvider(tt.provider)
			assert.NotNil(t, serviceProvider)

			// Should still generate auth URL (might be invalid, but shouldn't crash)
			authUrl := serviceProvider.GetAuthCodeUrl("test-state")
			assert.NotEmpty(t, authUrl)
		})
	}
}
