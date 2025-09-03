package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/mock"
)

// MockClientService is a mock implementation of client.Service
type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, queryParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) GetGoIamClients(ctx context.Context, params sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error) {
	args := m.Called(ctx, id, dontCheckProjects)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Client), args.Error(1)
}

func (m *MockClientService) Create(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Update(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Emit(event utils.Event[sdk.Client]) {
	m.Called(event)
}

func (m *MockClientService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Client], sdk.Client]) {
	m.Called(eventName, subscriber)
}

func (m *MockClientService) VerifySecret(plainSecret, hashedSecret string) error {
	args := m.Called(plainSecret, hashedSecret)
	return args.Error(0)
}
