package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/mock"
)

type MockResourceService struct {
	mock.Mock
}

func (m *MockResourceService) Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sdk.ResourceList), args.Error(1)
}

func (m *MockResourceService) Get(ctx context.Context, id string) (*sdk.Resource, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.Resource), args.Error(1)
}

func (m *MockResourceService) Create(ctx context.Context, resource *sdk.Resource) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

func (m *MockResourceService) Update(ctx context.Context, resource *sdk.Resource) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

func (m *MockResourceService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockResourceService) Emit(event utils.Event[sdk.Resource]) {
	m.Called(event)
}

func (m *MockResourceService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Resource], sdk.Resource]) {
	m.Called(eventName, subscriber)
}
