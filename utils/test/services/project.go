package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/mock"
)

// MockProjectService implements project.Service interface for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetAll(ctx context.Context) ([]sdk.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sdk.Project), args.Error(1)
}

func (m *MockProjectService) Get(ctx context.Context, id string) (*sdk.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) Create(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) Update(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}
