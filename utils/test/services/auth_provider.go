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
