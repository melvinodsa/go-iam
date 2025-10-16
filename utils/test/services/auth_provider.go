package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/mock"
)

type MockAuthProviderService struct {
	mock.Mock
}

func (m *MockAuthProviderService) GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sdk.AuthProvider), args.Error(1)
}

func (m *MockAuthProviderService) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.AuthProvider, error) {
	args := m.Called(ctx, id, dontCheckProjects)
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
	return args.Get(0).(sdk.ServiceProvider), args.Error(1)
}

// MockServiceProvider implements sdk.ServiceProvider interface for testing
type MockServiceProvider struct {
	mock.Mock
}

func (m *MockServiceProvider) HasRefreshTokenFlow() bool {
	args := m.Called()
	return args.Bool(0)
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
