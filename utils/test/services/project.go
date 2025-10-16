package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/mock"
)

// MockProjectService implements project service interface for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) Create(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) Get(ctx context.Context, id string) (*sdk.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) Update(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) GetAll(ctx context.Context) ([]sdk.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Project), args.Error(1)
}
