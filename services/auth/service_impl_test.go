package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
)

// Mock services for testing helper functions that need dependencies
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) Expire(ctx context.Context, key string, ttl time.Duration) error {
	args := m.Called(ctx, key, ttl)
	return args.Error(0)
}

// Additional mock services for main interface methods
type MockAuthProviderService struct {
	mock.Mock
}

func (m *MockAuthProviderService) GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.AuthProvider), args.Error(1)
}

func (m *MockAuthProviderService) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.AuthProvider, error) {
	args := m.Called(ctx, id, dontCheckProjects)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthProvider), args.Error(1)
}

func (m *MockAuthProviderService) Create(ctx context.Context, provider *sdk.AuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockAuthProviderService) Update(ctx context.Context, provider *sdk.AuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockAuthProviderService) GetProvider(ctx context.Context, v sdk.AuthProvider) (sdk.ServiceProvider, error) {
	args := m.Called(ctx, v)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sdk.ServiceProvider), args.Error(1)
}

type MockServiceProvider struct {
	mock.Mock
}

func (m *MockServiceProvider) GetAuthCodeUrl(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockServiceProvider) VerifyCode(ctx context.Context, code string) (*sdk.AuthToken, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthToken), args.Error(1)
}

func (m *MockServiceProvider) RefreshToken(refreshToken string) (*sdk.AuthToken, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthToken), args.Error(1)
}

func (m *MockServiceProvider) GetIdentity(token string) ([]sdk.AuthIdentity, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.AuthIdentity), args.Error(1)
}

func (m *MockServiceProvider) HasRefreshTokenFlow() bool {
	return true
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(claims map[string]interface{}, expiryTimeInSeconds int64) (string, error) {
	args := m.Called(claims, expiryTimeInSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type MockEncryptService struct {
	mock.Mock
}

func (m *MockEncryptService) Encrypt(rawMessage string) (string, error) {
	args := m.Called(rawMessage)
	return args.String(0), args.Error(1)
}

func (m *MockEncryptService) Decrypt(encryptedMessage string) (string, error) {
	args := m.Called(encryptedMessage)
	return args.String(0), args.Error(1)
}

// Helper function to create a fully mocked service
func setupFullTestService() (*service, *MockAuthProviderService, *services.MockClientService, *MockCacheService, *MockJWTService, *MockEncryptService, *services.MockUserService) {
	mockAuthProvider := &MockAuthProviderService{}
	mockClient := &services.MockClientService{}
	mockCache := &MockCacheService{}
	mockJWT := &MockJWTService{}
	mockEncrypt := &MockEncryptService{}
	mockUser := &services.MockUserService{}

	svc := &service{
		authP:      mockAuthProvider,
		clientSvc:  mockClient,
		cacheSvc:   mockCache,
		jwtSvc:     mockJWT,
		encSvc:     mockEncrypt,
		usrSvc:     mockUser,
		tokenTTL:   86400, // 24 hours
		refetchTTL: 3600,  // 1 hour
	}

	return svc, mockAuthProvider, mockClient, mockCache, mockJWT, mockEncrypt, mockUser
}

// TestNewService tests the NewService constructor function
func TestNewService(t *testing.T) {
	// Create mock services
	mockAuthProvider := &MockAuthProviderService{}
	mockClient := &services.MockClientService{}
	mockCache := &MockCacheService{}
	mockJWT := &MockJWTService{}
	mockEncrypt := &MockEncryptService{}
	mockUser := &services.MockUserService{}

	// Test parameters
	tokenTTL := int64(86400)  // 24 hours
	refetchTTL := int64(3600) // 1 hour

	// Call NewService
	result := NewService(
		mockAuthProvider,
		mockClient,
		mockCache,
		mockJWT,
		mockEncrypt,
		mockUser,
		tokenTTL,
		refetchTTL,
	)

	// Verify the result
	require.NotNil(t, result)

	// Verify all dependencies are properly set
	assert.Equal(t, mockAuthProvider, result.authP)
	assert.Equal(t, mockClient, result.clientSvc)
	assert.Equal(t, mockCache, result.cacheSvc)
	assert.Equal(t, mockJWT, result.jwtSvc)
	assert.Equal(t, mockEncrypt, result.encSvc)
	assert.Equal(t, mockUser, result.usrSvc)
	assert.Equal(t, tokenTTL, result.tokenTTL)
	assert.Equal(t, refetchTTL, result.refetchTTL)

	// Verify the returned type is correct
	assert.IsType(t, &service{}, result)
}

// TestCacheClientSecret tests the client secret caching
func TestCacheClientSecret(t *testing.T) {
	ctx := context.Background()
	mockCache := &MockCacheService{}

	svc := &service{
		cacheSvc: mockCache,
	}

	clientId := "test-client"
	secret := "test-secret"

	mockCache.On("Set", ctx, "client-test-client", "test-secret", time.Hour*24*365).Return(nil)

	svc.cacheClientSecret(ctx, clientId, secret)

	mockCache.AssertExpectations(t)
}

// TestCacheClientSecretError tests cache error handling
func TestCacheClientSecretError(t *testing.T) {
	ctx := context.Background()
	mockCache := &MockCacheService{}

	svc := &service{
		cacheSvc: mockCache,
	}

	clientId := "test-client"
	secret := "test-secret"

	mockCache.On("Set", ctx, "client-test-client", "test-secret", time.Hour*24*365).Return(errors.New("cache error"))

	// This should not panic even if cache fails
	svc.cacheClientSecret(ctx, clientId, secret)

	mockCache.AssertExpectations(t)
}

// TestGetClientSecret tests client secret retrieval
func TestGetClientSecret(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		clientId       string
		setupMocks     func(*MockCacheService, *services.MockClientService)
		expectedSecret string
		expectedError  string
	}{
		{
			name:     "success - from cache",
			clientId: "test-client",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-test-client").Return("cached-secret", nil)
			},
			expectedSecret: "cached-secret",
		},
		{
			name:     "success - from database",
			clientId: "test-client",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("not found"))
				mockClient.On("Get", ctx, "test-client", true).Return(&sdk.Client{
					Id:     "test-client",
					Secret: "db-secret",
				}, nil)
				mockCache.On("Set", ctx, "client-test-client", "db-secret", time.Hour*24*365).Return(nil)
			},
			expectedSecret: "db-secret",
		},
		{
			name:     "success - from database with cache error",
			clientId: "test-client",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("not found"))
				mockClient.On("Get", ctx, "test-client", true).Return(&sdk.Client{
					Id:     "test-client",
					Secret: "db-secret",
				}, nil)
				mockCache.On("Set", ctx, "client-test-client", "db-secret", time.Hour*24*365).Return(errors.New("cache error"))
			},
			expectedSecret: "db-secret", // Should still return secret even if cache fails
		},
		{
			name:     "error - client not found",
			clientId: "invalid-client",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-invalid-client").Return("", errors.New("not found"))
				mockClient.On("Get", ctx, "invalid-client", true).Return((*sdk.Client)(nil), errors.New("client not found"))
			},
			expectedError: "couldn't get the client even from db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := &MockCacheService{}
			mockClient := &services.MockClientService{}

			svc := &service{
				cacheSvc:  mockCache,
				clientSvc: mockClient,
			}

			tt.setupMocks(mockCache, mockClient)

			secret, err := svc.getClientSecret(ctx, tt.clientId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, secret)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSecret, secret)
			}

			mockCache.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestHandlePrivateClient tests private client validation
func TestHandlePrivateClient(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		clientId      string
		clientSecret  string
		setupMocks    func(*MockCacheService, *services.MockClientService)
		expectedError string
	}{
		{
			name:         "success",
			clientId:     "test-client",
			clientSecret: "correct-secret",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-test-client").Return("correct-secret", nil)
				mockClient.On("VerifySecret", "correct-secret", "correct-secret").Return(nil)
			},
		},
		{
			name:         "error - invalid client secret",
			clientId:     "test-client",
			clientSecret: "wrong-secret",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-test-client").Return("correct-secret", nil)
				mockClient.On("VerifySecret", mock.Anything, mock.Anything).Return(errors.New("invalid client secret"))
			},
			expectedError: "invalid client secret",
		},
		{
			name:         "error - client not found",
			clientId:     "invalid-client",
			clientSecret: "secret",
			setupMocks: func(mockCache *MockCacheService, mockClient *services.MockClientService) {
				mockCache.On("Get", ctx, "client-invalid-client").Return("", errors.New("not found"))
				mockClient.On("Get", ctx, "invalid-client", true).Return((*sdk.Client)(nil), errors.New("client not found"))
			},
			expectedError: "error getting client secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := &MockCacheService{}
			mockClient := &services.MockClientService{}

			svc := &service{
				cacheSvc:  mockCache,
				clientSvc: mockClient,
			}

			tt.setupMocks(mockCache, mockClient)

			err := svc.handlePrivateClient(ctx, tt.clientId, tt.clientSecret)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockCache.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestGenerateCodeChallengeS256 tests the PKCE code challenge generation
func TestGenerateCodeChallengeS256(t *testing.T) {
	tests := []struct {
		name          string
		codeChallenge string
	}{
		{
			name:          "basic test",
			codeChallenge: "test-challenge",
		},
		{
			name:          "empty string",
			codeChallenge: "",
		},
		{
			name:          "long string",
			codeChallenge: "very-long-code-challenge-with-many-characters-to-test-hash-function",
		},
		{
			name:          "special characters",
			codeChallenge: "test-challenge!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateCodeChallengeS256(tt.codeChallenge)

			// Verify it's properly base64 encoded
			_, err := base64.RawURLEncoding.DecodeString(result)
			assert.NoError(t, err)

			// Verify it matches manual calculation
			hash := sha256.Sum256([]byte(tt.codeChallenge))
			expected := base64.RawURLEncoding.EncodeToString(hash[:])
			assert.Equal(t, expected, result)

			// Verify consistent results
			result2 := generateCodeChallengeS256(tt.codeChallenge)
			assert.Equal(t, result, result2)

			// Verify the result is not empty (unless input is empty)
			if tt.codeChallenge != "" {
				assert.NotEmpty(t, result)
			}
		})
	}
}

// TestGenerateCodeChallengeS256EdgeCases tests edge cases
func TestGenerateCodeChallengeS256EdgeCases(t *testing.T) {
	t.Run("unicode characters", func(t *testing.T) {
		input := "test-challenge-ñáéíóú-中文-🔐"
		result := generateCodeChallengeS256(input)

		// Should not panic and should produce valid base64
		_, err := base64.RawURLEncoding.DecodeString(result)
		assert.NoError(t, err)

		// Should be reproducible
		result2 := generateCodeChallengeS256(input)
		assert.Equal(t, result, result2)
	})

	t.Run("very long input", func(t *testing.T) {
		// Create a very long string
		longInput := ""
		for i := 0; i < 1000; i++ {
			longInput += "test-challenge-"
		}

		result := generateCodeChallengeS256(longInput)

		// Should not panic and should produce valid base64
		_, err := base64.RawURLEncoding.DecodeString(result)
		assert.NoError(t, err)

		// SHA256 always produces 32 bytes, so base64 should be predictable length
		assert.Equal(t, 43, len(result)) // 32 bytes -> 43 chars in base64 without padding
	})
}

// TestHandlePublicClient tests the PKCE validation logic
func TestHandlePublicClient(t *testing.T) {
	svc := &service{} // No dependencies needed for this function

	// Helper to generate the expected hash
	generateExpectedHash := func(challenge string) string {
		hash := sha256.Sum256([]byte(challenge))
		return base64.RawURLEncoding.EncodeToString(hash[:])
	}

	originalChallenge := "test-verifier-string"
	expectedHash := generateExpectedHash(originalChallenge)

	tests := []struct {
		name          string
		clientId      string
		codeChallenge string // This is the hash that should match
		token         sdk.AuthToken
		expectedError string
	}{
		{
			name:          "success",
			clientId:      "test-client",
			codeChallenge: expectedHash, // Provide the correct hash
			token: sdk.AuthToken{
				ClientId:            "test-client",
				CodeChallenge:       originalChallenge, // Store the original verifier
				CodeChallengeMethod: "S256",
			},
		},
		{
			name:          "error - invalid code challenge method",
			clientId:      "test-client",
			codeChallenge: expectedHash,
			token: sdk.AuthToken{
				ClientId:            "test-client",
				CodeChallenge:       originalChallenge,
				CodeChallengeMethod: "plain",
			},
			expectedError: "invalid code challenge",
		},
		{
			name:          "error - invalid code verifier",
			clientId:      "test-client",
			codeChallenge: "wrong-hash", // Wrong hash
			token: sdk.AuthToken{
				ClientId:            "test-client",
				CodeChallenge:       originalChallenge, // Correct original verifier
				CodeChallengeMethod: "S256",
			},
			expectedError: "invalid code verifier",
		},
		{
			name:          "error - invalid client id",
			clientId:      "wrong-client",
			codeChallenge: expectedHash,
			token: sdk.AuthToken{
				ClientId:            "test-client",
				CodeChallenge:       originalChallenge,
				CodeChallengeMethod: "S256",
			},
			expectedError: "invalid client id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.handlePublicClient(tt.clientId, tt.codeChallenge, tt.token)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestGetLoginUrl tests the GetLoginUrl method - focusing on error cases that don't require complex mocking
func TestGetLoginUrl(t *testing.T) {
	ctx := context.Background()
	svc, mockAuthProvider, mockClient, mockCache, _, mockEncrypt, _ := setupFullTestService()

	tests := []struct {
		name                string
		clientId            string
		authProviderId      string
		state               string
		redirectUrl         string
		codeChallengeMethod string
		codeChallenge       string
		setupMocks          func()
		expectedError       string
	}{
		{
			name:                "error - client not found when no auth provider specified",
			clientId:            "invalid-client",
			authProviderId:      "",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				mockClient.On("Get", ctx, "invalid-client", true).Return((*sdk.Client)(nil), errors.New("client not found"))
			},
			expectedError: "error fetching client details",
		},
		{
			name:                "error - auth provider not found",
			clientId:            "test-client",
			authProviderId:      "invalid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				mockAuthProvider.On("Get", ctx, "invalid-provider", true).Return((*sdk.AuthProvider)(nil), errors.New("provider not found"))
			},
			expectedError: "error fetching auth provider details",
		},
		{
			name:                "error - auth provider not found with default provider lookup",
			clientId:            "test-client",
			authProviderId:      "", // Empty to trigger default provider lookup
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Client found but default auth provider not found
				client := &sdk.Client{
					Id:                    "test-client",
					DefaultAuthProviderId: "default-provider-id",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockAuthProvider.On("Get", ctx, "default-provider-id", true).Return((*sdk.AuthProvider)(nil), errors.New("default provider not found"))
			},
			expectedError: "error fetching auth provider details",
		},
		{
			name:                "error - service provider creation fails",
			clientId:            "test-client",
			authProviderId:      "valid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Auth provider found but service provider creation fails
				authProvider := &sdk.AuthProvider{
					Id:        "valid-provider",
					ProjectId: "project-123",
				}
				mockAuthProvider.On("Get", ctx, "valid-provider", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(nil, errors.New("service provider creation failed"))
			},
			expectedError: "error getting service provider",
		},
		{
			name:                "error - state caching fails due to encryption failure",
			clientId:            "test-client",
			authProviderId:      "valid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Auth provider and service provider succeed but state caching fails
				authProvider := &sdk.AuthProvider{
					Id:        "valid-provider",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("Get", ctx, "valid-provider", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// State encryption fails
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("", errors.New("encryption failed"))
			},
			expectedError: "error caching the state",
		},
		{
			name:                "error - state caching fails due to cache set failure",
			clientId:            "test-client",
			authProviderId:      "valid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Auth provider and service provider succeed, encryption succeeds, but cache set fails
				authProvider := &sdk.AuthProvider{
					Id:        "valid-provider",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("Get", ctx, "valid-provider", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// State encryption succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-state", nil)
				// Cache set fails
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-state", mock.Anything).Return(errors.New("cache set failed"))
			},
			expectedError: "error caching the state",
		},
		{
			name:                "success - with explicit auth provider",
			clientId:            "test-client",
			authProviderId:      "valid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Full successful flow with explicit auth provider
				authProvider := &sdk.AuthProvider{
					Id:        "valid-provider",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("Get", ctx, "valid-provider", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("GetAuthCodeUrl", mock.AnythingOfType("string")).Return("https://auth-provider.com/oauth/authorize?state=cached-state-id")
				// State caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-state", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-state", mock.Anything).Return(nil)
			},
			expectedError: "", // Should succeed
		},
		{
			name:                "success - with default auth provider lookup",
			clientId:            "test-client",
			authProviderId:      "", // Empty to trigger default provider lookup
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "S256",
			codeChallenge:       "test-challenge",
			setupMocks: func() {
				// Full successful flow with default auth provider lookup
				client := &sdk.Client{
					Id:                    "test-client",
					DefaultAuthProviderId: "default-provider-id",
				}
				authProvider := &sdk.AuthProvider{
					Id:        "default-provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockAuthProvider.On("Get", ctx, "default-provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("GetAuthCodeUrl", mock.AnythingOfType("string")).Return("https://auth-provider.com/oauth/authorize?state=cached-state-id")
				// State caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-state", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-state", mock.Anything).Return(nil)
			},
			expectedError: "", // Should succeed
		},
		{
			name:                "success - without code challenge",
			clientId:            "test-client",
			authProviderId:      "valid-provider",
			state:               "test-state",
			redirectUrl:         "http://localhost:3000/callback",
			codeChallengeMethod: "",
			codeChallenge:       "",
			setupMocks: func() {
				// Full successful flow without PKCE
				authProvider := &sdk.AuthProvider{
					Id:        "valid-provider",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("Get", ctx, "valid-provider", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("GetAuthCodeUrl", mock.AnythingOfType("string")).Return("https://auth-provider.com/oauth/authorize?state=cached-state-id")
				// State caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-state", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-state", mock.Anything).Return(nil)
			},
			expectedError: "", // Should succeed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockAuthProvider.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil
			mockEncrypt.ExpectedCalls = nil

			tt.setupMocks()

			url, err := svc.GetLoginUrl(ctx, tt.clientId, tt.authProviderId, tt.state, tt.redirectUrl, tt.codeChallengeMethod, tt.codeChallenge)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, url)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.Contains(t, url, "https://auth-provider.com/oauth/authorize")
				assert.Contains(t, url, "state=")
			}

			mockAuthProvider.AssertExpectations(t)
			mockClient.AssertExpectations(t)
			mockCache.AssertExpectations(t)
			mockEncrypt.AssertExpectations(t)
		})
	}
}

// TestRedirect tests the Redirect method - focusing on error cases
func TestRedirect(t *testing.T) {
	ctx := context.Background()
	svc, mockAuthProvider, mockClient, mockCache, _, mockEncrypt, _ := setupFullTestService()

	tests := []struct {
		name          string
		code          string
		state         string
		setupMocks    func()
		expectedError string
	}{
		{
			name:  "error - invalid state",
			code:  "auth-code",
			state: "invalid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-invalid-state").Return("", errors.New("state not found"))
			},
			expectedError: "error getting the state from cache",
		},
		{
			name:  "error - state decryption fails",
			code:  "auth-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("", errors.New("decryption failed"))
			},
			expectedError: "error getting the state from cache",
		},
		{
			name:  "error - invalid state format",
			code:  "auth-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				// State with only 4 parts instead of 6
				mockEncrypt.On("Decrypt", "encrypted-state").Return("state:client:provider:url", nil)
			},
			expectedError: "invalid state. expected to have 6 parts but got 4",
		},
		{
			name:  "error - invalid code challenge method",
			code:  "auth-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				// State with invalid code challenge method (not S256)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:SHA1:challenge", nil)
			},
			expectedError: "invalid code challenge",
		},
		{
			name:  "error - auth provider getToken fails",
			code:  "auth-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// getToken will fail when trying to get auth provider
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return((*sdk.AuthProvider)(nil), errors.New("auth provider not found"))
			},
			expectedError: "error getting the token",
		},
		{
			name:  "error - auth provider service provider creation fails",
			code:  "auth-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// Auth provider found but service provider creation fails
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(nil, errors.New("service provider creation failed"))
			},
			expectedError: "error getting the token",
		},
		{
			name:  "error - token verification fails",
			code:  "invalid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// Auth provider and service provider setup succeed but token verification fails
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "invalid-code").Return((*sdk.AuthToken)(nil), errors.New("invalid authorization code"))
			},
			expectedError: "error getting the token",
		},
		{
			name:  "error - auth token caching fails",
			code:  "valid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// Token verification succeeds but caching fails
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				authToken := &sdk.AuthToken{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					ExpiresAt:      time.Now().Add(time.Hour),
					AuthProviderID: "provider-id",
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "valid-code").Return(authToken, nil)
				// Caching fails during encryption
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("", errors.New("encryption failed"))
			},
			expectedError: "error caching the token",
		},
		{
			name:  "error - redirect URL validation fails",
			code:  "valid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// Token operations succeed but redirect URL validation fails
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				authToken := &sdk.AuthToken{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					ExpiresAt:      time.Now().Add(time.Hour),
					AuthProviderID: "provider-id",
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "valid-code").Return(authToken, nil)
				// Auth token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-auth-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-auth-token", mock.Anything).Return(nil)
				// Client validation fails - client not found
				mockClient.On("Get", ctx, "client-id", true).Return((*sdk.Client)(nil), errors.New("client not found"))
			},
			expectedError: "error getting the callback url",
		},
		{
			name:  "error - redirect URL not in allowed list",
			code:  "valid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//unauthorized.com:S256:challenge", nil)
				// Token operations succeed but redirect URL not in client's allowed URLs
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				authToken := &sdk.AuthToken{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					ExpiresAt:      time.Now().Add(time.Hour),
					AuthProviderID: "provider-id",
				}
				client := &sdk.Client{
					Id:           "client-id",
					RedirectURLs: []string{"http://authorized.com", "https://app.example.com"},
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "valid-code").Return(authToken, nil)
				// Auth token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-auth-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-auth-token", mock.Anything).Return(nil)
				// Client found but redirect URL not in allowed list
				mockClient.On("Get", ctx, "client-id", true).Return(client, nil)
			},
			expectedError: "callback url not found in the client details",
		},
		{
			name:  "error - URL decoding fails in state",
			code:  "valid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				// State with invalid URL encoding that will fail to decode
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com%ZZ:S256:challenge", nil)
			},
			expectedError: "error getting the state from cache",
		},
		{
			name:  "success - state invalidation called on successful redirect",
			code:  "valid-code",
			state: "valid-state",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com:S256:challenge", nil)
				// Full successful flow to verify state invalidation is called
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				authToken := &sdk.AuthToken{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					ExpiresAt:      time.Now().Add(time.Hour),
					AuthProviderID: "provider-id",
				}
				client := &sdk.Client{
					Id:           "client-id",
					RedirectURLs: []string{"http://callback.com"},
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "valid-code").Return(authToken, nil)
				// Auth token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-auth-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-auth-token", mock.Anything).Return(nil)
				// Client found and redirect URL is valid
				mockClient.On("Get", ctx, "client-id", true).Return(client, nil)
				// State invalidation should be called - expect it to succeed
				mockCache.On("Delete", ctx, "state-valid-state").Return(nil)
			},
			expectedError: "", // This should succeed
		},
		{
			name:  "warning - state invalidation fails but redirect succeeds",
			code:  "valid-code",
			state: "valid-state-fail-invalidation",
			setupMocks: func() {
				mockCache.On("Get", ctx, "state-valid-state-fail-invalidation").Return("encrypted-state", nil)
				mockEncrypt.On("Decrypt", "encrypted-state").Return("original-state:client-id:provider-id:http%3A//callback.com::challenge", nil)
				// Full successful flow but state invalidation fails
				authProvider := &sdk.AuthProvider{
					Id:        "provider-id",
					ProjectId: "project-123",
				}
				mockServiceProvider := &MockServiceProvider{}
				authToken := &sdk.AuthToken{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					ExpiresAt:      time.Now().Add(time.Hour),
					AuthProviderID: "provider-id",
				}
				client := &sdk.Client{
					Id:           "client-id",
					RedirectURLs: []string{"http://callback.com"},
				}
				mockAuthProvider.On("Get", ctx, "provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				mockServiceProvider.On("VerifyCode", ctx, "valid-code").Return(authToken, nil)
				// Auth token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-auth-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-auth-token", mock.Anything).Return(nil)
				// Client found and redirect URL is valid
				mockClient.On("Get", ctx, "client-id", true).Return(client, nil)
				// State invalidation fails but should not affect the overall result
				mockCache.On("Delete", ctx, "state-valid-state-fail-invalidation").Return(errors.New("cache deletion failed"))
			},
			expectedError: "", // This should still succeed despite invalidation failure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCache.ExpectedCalls = nil
			mockEncrypt.ExpectedCalls = nil
			mockAuthProvider.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil

			tt.setupMocks()

			result, err := svc.Redirect(ctx, tt.code, tt.state)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.RedirectUrl)
				// Verify the redirect URL contains the expected parameters
				assert.Contains(t, result.RedirectUrl, "code=")
				assert.Contains(t, result.RedirectUrl, "state=")
			}

			mockCache.AssertExpectations(t)
			mockEncrypt.AssertExpectations(t)
			mockAuthProvider.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestClientCallback tests the ClientCallback method - focusing on error cases
func TestClientCallback(t *testing.T) {
	ctx := context.Background()
	svc, _, mockClient, mockCache, mockJWT, mockEncrypt, _ := setupFullTestService()

	// Pre-calculate the values for the client ID mismatch test
	originalVerifier := "original-verifier"
	hash := sha256.Sum256([]byte(originalVerifier))
	calculatedHash := base64.RawURLEncoding.EncodeToString(hash[:])

	tests := []struct {
		name          string
		code          string
		codeChallenge string
		clientId      string
		clientSecret  string
		setupMocks    func()
		expectedError string
	}{
		{
			name:          "error - invalid code",
			code:          "invalid-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				mockCache.On("Get", ctx, "auth-code-invalid-code").Return("", errors.New("code not found"))
			},
			expectedError: "error getting the token from cache",
		},
		{
			name:          "error - auth token decryption fails",
			code:          "valid-code-decrypt-fail",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				mockCache.On("Get", ctx, "auth-code-valid-code-decrypt-fail").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return("", errors.New("decryption failed"))
			},
			expectedError: "error decrypting the access token",
		},
		{
			name:          "error - auth token JSON unmarshal fails",
			code:          "valid-code-unmarshal-fail",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				mockCache.On("Get", ctx, "auth-code-valid-code-unmarshal-fail").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return("invalid-json", nil)
			},
			expectedError: "error decoding the token",
		},
		{
			name:          "error - private client - client not found",
			code:          "valid-code",
			codeChallenge: "",
			clientId:      "unknown-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				// Mock successful auth token retrieval
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then client not found in DB
				mockCache.On("Get", ctx, "client-unknown-client").Return("", errors.New("cache miss"))
				mockClient.On("Get", ctx, "unknown-client", true).Return((*sdk.Client)(nil), errors.New("client not found"))
			},
			expectedError: "error handling private client",
		},
		{
			name:          "error - private client - invalid client secret",
			code:          "valid-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "wrong-secret",
			setupMocks: func() {
				// Mock successful auth token retrieval
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then found in DB but wrong secret
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("cache miss"))
				client := &sdk.Client{
					Id:     "test-client",
					Secret: "correct-secret",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockClient.On("VerifySecret", mock.Anything, mock.Anything).Return(errors.New("invalid client secret"))
				mockCache.On("Set", ctx, "client-test-client", "correct-secret", mock.Anything).Return(nil)
			},
			expectedError: "invalid client secret",
		},
		{
			name:          "error - public client - invalid code challenge method",
			code:          "valid-code",
			codeChallenge: "test-challenge",
			clientId:      "test-client",
			clientSecret:  "",
			setupMocks: func() {
				// Mock successful auth token retrieval with invalid code challenge method
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"wrong-hash","code_challenge_method":"SHA1"}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
			},
			expectedError: "invalid code challenge",
		},
		{
			name:          "error - public client - invalid code verifier",
			code:          "valid-code",
			codeChallenge: "wrong-verifier",
			clientId:      "test-client",
			clientSecret:  "",
			setupMocks: func() {
				// Mock successful auth token retrieval with wrong code challenge
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"different-hash","code_challenge_method":"S256"}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
			},
			expectedError: "invalid code verifier",
		},
		{
			name:          "error - public client - client ID mismatch",
			code:          "valid-code",
			codeChallenge: calculatedHash, // Use the pre-calculated hash
			clientId:      "different-client",
			clientSecret:  "",
			setupMocks: func() {
				// Mock successful auth token retrieval with matching code challenge but different client ID
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"original-client","code_challenge":"` + originalVerifier + `","code_challenge_method":"S256"}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
			},
			expectedError: "invalid client id",
		},
		{
			name:          "error - access token caching fails",
			code:          "valid-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				// Mock successful auth token retrieval and private client validation
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then found in DB with correct secret
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("cache miss"))
				client := &sdk.Client{
					Id:     "test-client",
					Secret: "test-secret",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockClient.On("VerifySecret", "test-secret", "test-secret").Return(nil)
				mockCache.On("Set", ctx, "client-test-client", "test-secret", mock.Anything).Return(nil)
				// Access token caching fails during encryption
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("", errors.New("encryption failed"))
			},
			expectedError: "error encrypting the access token",
		},
		{
			name:          "error - JWT token generation fails",
			code:          "valid-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				// Mock successful auth token retrieval and private client validation
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then found in DB with correct secret
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("cache miss"))
				client := &sdk.Client{
					Id:     "test-client",
					Secret: "test-secret",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockClient.On("VerifySecret", "test-secret", "test-secret").Return(nil)
				mockCache.On("Set", ctx, "client-test-client", "test-secret", mock.Anything).Return(nil)
				// Access token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-access-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-access-token", mock.Anything).Return(nil)
				// JWT generation fails
				mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("int64")).Return("", errors.New("JWT generation failed"))
			},
			expectedError: "error generating the access token",
		},
		{
			name:          "error - auth token invalidation fails",
			code:          "valid-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				// Mock successful flow until invalidation
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-valid-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then found in DB with correct secret
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("cache miss"))
				client := &sdk.Client{
					Id:     "test-client",
					Secret: "test-secret",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockClient.On("VerifySecret", "test-secret", "test-secret").Return(nil)
				mockCache.On("Set", ctx, "client-test-client", "test-secret", mock.Anything).Return(nil)
				// Access token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-access-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-access-token", mock.Anything).Return(nil)
				// JWT generation succeeds
				mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("int64")).Return("jwt-token", nil)
				// Auth token invalidation fails
				mockCache.On("Delete", ctx, "auth-code-valid-code").Return(errors.New("cache deletion failed"))
			},
			expectedError: "error invalidating the auth code",
		},
		{
			name:          "success - private client flow",
			code:          "success-code",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "test-secret",
			setupMocks: func() {
				// Mock complete successful private client flow
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-success-code").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// Client secret cache miss, then found in DB with correct secret
				mockCache.On("Get", ctx, "client-test-client").Return("", errors.New("cache miss"))
				client := &sdk.Client{
					Id:     "test-client",
					Secret: "test-secret",
				}
				mockClient.On("Get", ctx, "test-client", true).Return(client, nil)
				mockClient.On("VerifySecret", "test-secret", "test-secret").Return(nil)
				mockCache.On("Set", ctx, "client-test-client", "test-secret", mock.Anything).Return(nil)
				// Access token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-access-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-access-token", mock.Anything).Return(nil)
				// JWT generation succeeds
				mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("int64")).Return("jwt-token", nil)
				// Auth token invalidation succeeds
				mockCache.On("Delete", ctx, "auth-code-success-code").Return(nil)
			},
			expectedError: "", // Should succeed
		},
		{
			name:          "success - public client flow",
			code:          "success-code-public",
			codeChallenge: calculatedHash, // Use the pre-calculated hash
			clientId:      "test-client",
			clientSecret:  "",
			setupMocks: func() {
				// Mock complete successful public client flow
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"` + originalVerifier + `","code_challenge_method":"S256"}`
				mockCache.On("Get", ctx, "auth-code-success-code-public").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// No client secret validation needed for public client
				// Access token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-access-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-access-token", mock.Anything).Return(nil)
				// JWT generation succeeds
				mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("int64")).Return("jwt-token", nil)
				// Auth token invalidation succeeds
				mockCache.On("Delete", ctx, "auth-code-success-code-public").Return(nil)
			},
			expectedError: "", // Should succeed
		},
		{
			name:          "success - neither private nor public client",
			code:          "success-code-neither",
			codeChallenge: "",
			clientId:      "test-client",
			clientSecret:  "",
			setupMocks: func() {
				// Mock complete successful flow with neither private nor public client validation
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","client_id":"test-client","code_challenge":"","code_challenge_method":""}`
				mockCache.On("Get", ctx, "auth-code-success-code-neither").Return("encrypted-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-data").Return(tokenJSON, nil)
				// No client validation needed
				// Access token caching succeeds
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-access-token", nil)
				mockCache.On("Set", ctx, mock.AnythingOfType("string"), "encrypted-access-token", mock.Anything).Return(nil)
				// JWT generation succeeds
				mockJWT.On("GenerateToken", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("int64")).Return("jwt-token", nil)
				// Auth token invalidation succeeds
				mockCache.On("Delete", ctx, "auth-code-success-code-neither").Return(nil)
			},
			expectedError: "", // Should succeed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCache.ExpectedCalls = nil
			mockEncrypt.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			mockJWT.ExpectedCalls = nil

			tt.setupMocks()

			result, err := svc.ClientCallback(ctx, tt.code, tt.codeChallenge, tt.clientId, tt.clientSecret)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "jwt-token", result.AccessToken)
			}

			mockCache.AssertExpectations(t)
			mockEncrypt.AssertExpectations(t)
			mockClient.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}

// Mock metadata types for testing
type MockEmailMetadata struct {
	Email string
}

func (m MockEmailMetadata) UpdateUserDetails(user *sdk.User) {
	user.Email = m.Email
}

type MockPhoneMetadata struct {
	Phone string
}

func (m MockPhoneMetadata) UpdateUserDetails(user *sdk.User) {
	user.Phone = m.Phone
}

type MockNameMetadata struct {
	Name string
}

func (m MockNameMetadata) UpdateUserDetails(user *sdk.User) {
	user.Name = m.Name
}

// TestGetIdentity tests the GetIdentity method
func TestGetIdentity(t *testing.T) {
	ctx := context.Background()
	svc, mockAuthProvider, _, mockCache, mockJWT, mockEncrypt, mockUser := setupFullTestService()

	tests := []struct {
		name          string
		accessToken   string
		setupMocks    func()
		expectedError string
		expectedUser  *sdk.User
	}{
		{
			name:        "success - full flow",
			accessToken: "valid-token-full-flow",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-full-flow").Return(claims, nil)

				// First it tries to get user from cache (this fails)
				mockCache.On("Get", ctx, "token-valid-token-full-flow").Return("", errors.New("user not found in cache"))

				// Then it gets access token from cache (this succeeds)
				encryptedToken := "encrypted-token-data"
				mockCache.On("Get", ctx, "access-token-token-123").Return(encryptedToken, nil)

				// Decrypt the access token
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				//`{"access_token":"at_123","refresh_token":"rt_123","auth_provider_id":"google-provider-id","expires_at":"2025-08-31T10:00:00Z"}`
				mockEncrypt.On("Decrypt", mock.Anything).Return(string(tokenJSON), nil)

				// Mock auth provider
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)

				// Mock service provider
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)

				// Mock GetIdentity call
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "test@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)

				// Mock GetOrCreateUser flow (internal method)
				// First it tries GetByEmail (user doesn't exist, returns ErrorUserNotFound)
				mockUser.On("GetByEmail", ctx, "test@example.com", "test-project").Return((*sdk.User)(nil), user.ErrorUserNotFound)

				// Then it calls Create to create the user
				mockUser.On("Create", ctx, mock.MatchedBy(func(user *sdk.User) bool {
					return user.Email == "test@example.com" && user.Name == "Test User" && user.ProjectId == "test-project"
				})).Run(func(args mock.Arguments) {
					// Simulate setting the ID after creation
					u := args.Get(1).(*sdk.User)
					u.Id = "user-123"
					u.Enabled = true // Make sure user is enabled
				}).Return(nil)

				// Cache user details
				mockEncrypt.On("Encrypt", mock.Anything).Return("encrypted-user-data", nil)
				mockCache.On("Set", ctx, "token-valid-token-full-flow", "encrypted-user-data", mock.Anything).Return(nil)
			},
			expectedUser: &sdk.User{
				Id:        "user-123",
				Email:     "test@example.com",
				Name:      "Test User",
				ProjectId: "test-project",
			},
		},
		{
			name:        "success - user cache hit",
			accessToken: "valid-token-cache-hit",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-cache-hit").Return(claims, nil)

				// Mock successful user cache hit - getUserFromCache succeeds
				encryptedUserData := "encrypted-user-cache-data"
				mockCache.On("Get", ctx, "token-valid-token-cache-hit").Return(encryptedUserData, nil)

				// Mock successful decryption of cached user data
				userJSON := `{"id":"cached-user-123","email":"cached@example.com","name":"Cached User","project_id":"test-project","enabled":true}`
				mockEncrypt.On("Decrypt", encryptedUserData).Return(userJSON, nil)

				// No other service calls should be made since we return early from cache hit
			},
			expectedUser: &sdk.User{
				Id:        "cached-user-123",
				Email:     "cached@example.com",
				Name:      "Cached User",
				ProjectId: "test-project",
				Enabled:   true,
			},
		},
		{
			name:        "error - user cache hit but decryption fails",
			accessToken: "valid-token-cache-decrypt-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-cache-decrypt-fail").Return(claims, nil)

				// Mock user cache hit but decryption fails
				encryptedUserData := "corrupted-encrypted-user-data"
				mockCache.On("Get", ctx, "token-valid-token-cache-decrypt-fail").Return(encryptedUserData, nil)
				mockEncrypt.On("Decrypt", encryptedUserData).Return("", errors.New("decryption failed"))

				// Since getUserFromCache fails, it should continue with normal flow
				// Mock access token cache hit for the normal flow
				encryptedToken := "encrypted-token-data"
				mockCache.On("Get", ctx, "access-token-token-123").Return(encryptedToken, nil)

				// Decrypt the access token successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockEncrypt.On("Decrypt", encryptedToken).Return(string(tokenJSON), nil)

				// Mock auth provider
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)

				// Mock service provider
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)

				// Mock GetIdentity call
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "fallback@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Fallback User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)

				// Mock existing user found
				existingUser := &sdk.User{
					Id:        "existing-user-456",
					Email:     "fallback@example.com",
					Name:      "Fallback User",
					ProjectId: "test-project",
					Enabled:   true,
				}
				mockUser.On("GetByEmail", ctx, "fallback@example.com", "test-project").Return(existingUser, nil)

				// Cache user details
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-user-data", nil)
				mockCache.On("Set", ctx, "token-valid-token-cache-decrypt-fail", "encrypted-user-data", mock.Anything).Return(nil)
			},
			expectedUser: &sdk.User{
				Id:        "existing-user-456",
				Email:     "fallback@example.com",
				Name:      "Fallback User",
				ProjectId: "test-project",
				Enabled:   true,
			},
		},
		{
			name:        "error - user cache hit but JSON decode fails",
			accessToken: "valid-token-cache-json-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-cache-json-fail").Return(claims, nil)

				// Mock user cache hit but JSON decode fails
				encryptedUserData := "encrypted-user-cache-data"
				mockCache.On("Get", ctx, "token-valid-token-cache-json-fail").Return(encryptedUserData, nil)
				mockEncrypt.On("Decrypt", encryptedUserData).Return("invalid-json-data", nil)

				// Since getUserFromCache fails, it should continue with normal flow
				// Mock access token cache hit for the normal flow
				encryptedToken := "encrypted-token-data"
				mockCache.On("Get", ctx, "access-token-token-123").Return(encryptedToken, nil)

				// Decrypt the access token successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				// tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","auth_provider_id":"google-provider-id","expires_at":"2025-08-31T10:00:00Z"}`
				mockEncrypt.On("Decrypt", encryptedToken).Return(string(tokenJSON), nil)

				// Mock auth provider
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)

				// Mock service provider
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)

				// Mock GetIdentity call
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "recovery@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Recovery User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)

				// Mock existing user found
				existingUser := &sdk.User{
					Id:        "recovery-user-789",
					Email:     "recovery@example.com",
					Name:      "Recovery User",
					ProjectId: "test-project",
					Enabled:   true,
				}
				mockUser.On("GetByEmail", ctx, "recovery@example.com", "test-project").Return(existingUser, nil)

				// Cache user details
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("encrypted-user-data", nil)
				mockCache.On("Set", ctx, "token-valid-token-cache-json-fail", "encrypted-user-data", mock.Anything).Return(nil)
			},
			expectedUser: &sdk.User{
				Id:        "recovery-user-789",
				Email:     "recovery@example.com",
				Name:      "Recovery User",
				ProjectId: "test-project",
				Enabled:   true,
			},
		},
		{
			name:        "error - invalid token",
			accessToken: "invalid-token",
			setupMocks: func() {
				mockJWT.On("ValidateToken", "invalid-token").Return(map[string]interface{}(nil), errors.New("invalid token"))
			},
			expectedError: "error validating the access token",
		},
		{
			name:        "error - missing id claim",
			accessToken: "token-without-id",
			setupMocks: func() {
				claims := map[string]interface{}{
					"email": "test@example.com",
				}
				mockJWT.On("ValidateToken", "token-without-id").Return(claims, nil)
			},
			expectedError: "error getting the access token id from claims",
		},
		{
			name:        "error - access token not found in cache",
			accessToken: "valid-token-no-access-token",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-no-access-token").Return(claims, nil)
				// First it tries to get user from cache using the access token directly (this can fail)
				mockCache.On("Get", ctx, "token-valid-token-no-access-token").Return("", errors.New("user not found in cache"))
				// Then it tries to get access token from cache using the id from claims (this should also fail)
				mockCache.On("Get", ctx, "access-token-token-123").Return("", errors.New("access token not found in cache"))
			},
			expectedError: "error getting the token from cache",
		},
		{
			name:        "error - access token decryption fails",
			accessToken: "valid-token-decrypt-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-decrypt-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-decrypt-fail").Return("", errors.New("user not found in cache"))
				// Access token found but decryption fails
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return("", errors.New("decryption failed"))
			},
			expectedError: "error decrypting the access token",
		},
		{
			name:        "error - access token JSON unmarshal fails",
			accessToken: "valid-token-unmarshal-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-unmarshal-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-unmarshal-fail").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted but invalid JSON
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return("invalid-json", nil)
			},
			expectedError: "error decoding the token",
		},
		{
			name:        "error - auth provider not found",
			accessToken: "valid-token-no-provider",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-no-provider").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-no-provider").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","auth_provider_id":"unknown-provider","expires_at":"2025-08-31T10:00:00Z"}`
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(tokenJSON, nil)
				// Auth provider not found
				mockAuthProvider.On("Get", ctx, "unknown-provider", true).Return((*sdk.AuthProvider)(nil), errors.New("provider not found"))
			},
			expectedError: "error getting the identity from auth provider error fetching auth provider details provider not found",
		},
		{
			name:        "error - service provider creation fails",
			accessToken: "valid-token-sp-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-sp-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-sp-fail").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				tokenJSON := `{"access_token":"at_123","refresh_token":"rt_123","auth_provider_id":"google-provider-id","expires_at":"2025-08-31T10:00:00Z"}`
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(tokenJSON, nil)
				// Auth provider found but service provider creation fails
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(nil, errors.New("service provider creation failed"))
			},
			expectedError: "error getting service provider",
		},
		{
			name:        "error - service provider GetIdentity fails",
			accessToken: "valid-token-getidentity-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-getidentity-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-getidentity-fail").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity fails
				mockServiceProvider.On("GetIdentity", "at_123").Return(([]sdk.AuthIdentity)(nil), errors.New("failed to get identity from provider"))
			},
			expectedError: "error getting the identity from service provider",
		},
		{
			name:        "error - user has no email or phone",
			accessToken: "valid-token-no-contact",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-no-contact").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-no-contact").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity returns identity with no email or phone
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}}, // Only name, no email/phone
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)
			},
			expectedError: "email or phone is required",
		},
		{
			name:        "error - user creation fails",
			accessToken: "valid-token-create-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-create-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-create-fail").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity returns valid identity
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "test@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)
				// User doesn't exist
				mockUser.On("GetByEmail", ctx, "test@example.com", "test-project").Return((*sdk.User)(nil), user.ErrorUserNotFound)
				// User creation fails
				mockUser.On("Create", ctx, mock.MatchedBy(func(user *sdk.User) bool {
					return user.Email == "test@example.com" && user.Name == "Test User" && user.ProjectId == "test-project"
				})).Return(errors.New("user creation failed"))
			},
			expectedError: "error creating the user",
		},
		{
			name:        "error - user is disabled",
			accessToken: "valid-token-user-disabled",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-user-disabled").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-user-disabled").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity returns valid identity
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "test@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)
				// User exists but is disabled
				disabledUser := &sdk.User{
					Id:        "user-123",
					Email:     "test@example.com",
					Name:      "Test User",
					ProjectId: "test-project",
					Enabled:   false, // User is disabled
				}
				mockUser.On("GetByEmail", ctx, "test@example.com", "test-project").Return(disabledUser, nil)
			},
			expectedError: "user is disabled",
		},
		{
			name:        "error - user is expired",
			accessToken: "valid-token-user-expired",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-user-expired").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-user-expired").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity returns valid identity
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "test@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)
				// User exists but is expired
				yesterday := time.Now().Add(-24 * time.Hour)
				expiredUser := &sdk.User{
					Id:        "user-123",
					Email:     "test@example.com",
					Name:      "Test User",
					ProjectId: "test-project",
					Enabled:   true,
					Expiry:    &yesterday, // User is expired
				}
				mockUser.On("GetByEmail", ctx, "test@example.com", "test-project").Return(expiredUser, nil)
			},
			expectedError: "user expired",
		},
		{
			name:        "error - cache user details fails",
			accessToken: "valid-token-cache-fail",
			setupMocks: func() {
				claims := map[string]interface{}{
					"id": "token-123",
				}
				mockJWT.On("ValidateToken", "valid-token-cache-fail").Return(claims, nil)
				// User not found in cache
				mockCache.On("Get", ctx, "token-valid-token-cache-fail").Return("", errors.New("user not found in cache"))
				// Access token found and decrypted successfully
				token := &sdk.AuthToken{
					AccessToken:    "at_123",
					RefreshToken:   "rt_123",
					AuthProviderID: "google-provider-id",
					ExpiresAt:      time.Now().Add(24 * time.Hour),
				}
				tokenJSON, err := json.Marshal(token)
				if err != nil {
					log.Printf("Error encoding token JSON: %v", err)
				}
				mockCache.On("Get", ctx, "access-token-token-123").Return("encrypted-token-data", nil)
				mockEncrypt.On("Decrypt", "encrypted-token-data").Return(string(tokenJSON), nil)
				// Auth provider and service provider found
				authProvider := &sdk.AuthProvider{
					Id:        "google-provider-id",
					ProjectId: "test-project",
				}
				mockAuthProvider.On("Get", ctx, "google-provider-id", true).Return(authProvider, nil)
				mockServiceProvider := &MockServiceProvider{}
				mockAuthProvider.On("GetProvider", ctx, *authProvider).Return(mockServiceProvider, nil)
				// GetIdentity returns valid identity
				identities := []sdk.AuthIdentity{
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockEmailMetadata{Email: "test@example.com"}},
					{Type: sdk.AuthIdentityTypeEmail, Metadata: MockNameMetadata{Name: "Test User"}},
				}
				mockServiceProvider.On("GetIdentity", "at_123").Return(identities, nil)
				// User creation succeeds
				mockUser.On("GetByEmail", ctx, "test@example.com", "test-project").Return((*sdk.User)(nil), user.ErrorUserNotFound)
				mockUser.On("Create", ctx, mock.MatchedBy(func(user *sdk.User) bool {
					return user.Email == "test@example.com" && user.Name == "Test User" && user.ProjectId == "test-project"
				})).Run(func(args mock.Arguments) {
					u := args.Get(1).(*sdk.User)
					u.Id = "user-123"
					u.Enabled = true
				}).Return(nil)
				// Cache user details fails during encryption
				mockEncrypt.On("Encrypt", mock.AnythingOfType("string")).Return("", errors.New("encryption failed"))
			},
			expectedError: "error encrypting the user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockJWT.ExpectedCalls = nil
			mockCache.ExpectedCalls = nil
			mockEncrypt.ExpectedCalls = nil
			mockAuthProvider.ExpectedCalls = nil
			mockUser.ExpectedCalls = nil

			tt.setupMocks()

			user, err := svc.GetIdentity(ctx, tt.accessToken)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Id, user.Id)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.ProjectId, user.ProjectId)
			}

			mockJWT.AssertExpectations(t)
			mockCache.AssertExpectations(t)
			mockEncrypt.AssertExpectations(t)
			mockAuthProvider.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	}
}

// TestHandleEvent tests the HandleEvent method
func TestHandleEvent(t *testing.T) {
	ctx := context.Background()
	svc, _, _, mockCache, _, _, _ := setupFullTestService()

	// Create a mock event
	mockEvent := &MockEvent{
		name: goiamuniverse.EventClientCreated,
		payload: sdk.Client{
			Id:     "test-client",
			Secret: "test-secret",
		},
	}

	tests := []struct {
		name       string
		event      utils.Event[sdk.Client]
		setupMocks func()
	}{
		{
			name:  "client created event",
			event: mockEvent,
			setupMocks: func() {
				mockCache.On("Set", ctx, "client-test-client", "test-secret", time.Hour*24*365).Return(nil)
			},
		},
		{
			name: "client updated event",
			event: &MockEvent{
				name: goiamuniverse.EventClientUpdated,
				payload: sdk.Client{
					Id:     "test-client",
					Secret: "updated-secret",
				},
			},
			setupMocks: func() {
				mockCache.On("Set", ctx, "client-test-client", "updated-secret", time.Hour*24*365).Return(nil)
			},
		},
		{
			name: "unknown event - should be ignored",
			event: &MockEvent{
				name: goiamuniverse.Event("unknown-event"),
				payload: sdk.Client{
					Id:     "test-client",
					Secret: "test-secret",
				},
			},
			setupMocks: func() {
				// No cache calls expected for unknown events
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCache.ExpectedCalls = nil

			tt.setupMocks()

			svc.HandleEvent(tt.event)

			mockCache.AssertExpectations(t)
		})
	}
}

// Mock Event implementation for testing
type MockEvent struct {
	name     goiamuniverse.Event
	payload  sdk.Client
	metadata sdk.Metadata
	context  context.Context
}

func (m *MockEvent) Name() goiamuniverse.Event {
	return m.name
}

func (m *MockEvent) Payload() sdk.Client {
	return m.payload
}

func (m *MockEvent) Metadata() sdk.Metadata {
	return m.metadata
}

func (m *MockEvent) Context() context.Context {
	if m.context == nil {
		return context.Background()
	}
	return m.context
}
