package project

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore implements Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetAll(ctx context.Context) ([]sdk.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Project), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*sdk.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockStore) Create(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockStore) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func TestNewService(t *testing.T) {
	mockStore := &MockStore{}

	service := NewService(mockStore)

	assert.NotNil(t, service)
	assert.Implements(t, (*Service)(nil), service)
}

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockStore)
		expectedError string
	}{
		{
			name: "successful get all",
			setupMock: func(ms *MockStore) {
				ms.On("GetAll", mock.Anything).Return([]sdk.Project{
					{Id: "project1", Name: "Project 1"},
					{Id: "project2", Name: "Project 2"},
				}, nil)
			},
		},
		{
			name: "store error",
			setupMock: func(ms *MockStore) {
				ms.On("GetAll", mock.Anything).Return(nil, errors.New("store error"))
			},
			expectedError: "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}

			tt.setupMock(mockStore)

			service := NewService(mockStore)

			result, err := service.GetAll(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, 2)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name          string
		projectID     string
		setupMock     func(*MockStore)
		expectedError string
	}{
		{
			name:      "successful get",
			projectID: "project1",
			setupMock: func(ms *MockStore) {
				ms.On("Get", mock.Anything, "project1").Return(&sdk.Project{
					Id:   "project1",
					Name: "Test Project",
				}, nil)
			},
		},
		{
			name:          "empty project ID",
			projectID:     "",
			setupMock:     func(ms *MockStore) {},
			expectedError: "project not found",
		},
		{
			name:      "project not found",
			projectID: "nonexistent",
			setupMock: func(ms *MockStore) {
				ms.On("Get", mock.Anything, "nonexistent").Return(nil, errors.New("project not found"))
			},
			expectedError: "project not found",
		},
		{
			name:      "store error",
			projectID: "project1",
			setupMock: func(ms *MockStore) {
				ms.On("Get", mock.Anything, "project1").Return(nil, errors.New("store error"))
			},
			expectedError: "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}

			tt.setupMock(mockStore)

			service := NewService(mockStore)

			result, err := service.Get(context.Background(), tt.projectID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.projectID, result.Id)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_GetByName(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		setupMock     func(*MockStore)
		expectedError string
	}{
		{
			name:        "successful get by name",
			projectName: "Test Project",
			setupMock: func(ms *MockStore) {
				ms.On("GetByName", mock.Anything, "Test Project").Return(&sdk.Project{
					Id:   "project1",
					Name: "Test Project",
				}, nil)
			},
		},
		{
			name:          "empty project name",
			projectName:   "",
			setupMock:     func(ms *MockStore) {},
			expectedError: "project not found",
		},
		{
			name:        "project not found",
			projectName: "Nonexistent Project",
			setupMock: func(ms *MockStore) {
				ms.On("GetByName", mock.Anything, "Nonexistent Project").Return(nil, errors.New("project not found"))
			},
			expectedError: "project not found",
		},
		{
			name:        "store error",
			projectName: "Test Project",
			setupMock: func(ms *MockStore) {
				ms.On("GetByName", mock.Anything, "Test Project").Return(nil, errors.New("store error"))
			},
			expectedError: "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}

			tt.setupMock(mockStore)

			service := NewService(mockStore)

			result, err := service.GetByName(context.Background(), tt.projectName)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.projectName, result.Name)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		project       *sdk.Project
		setupMock     func(*MockStore)
		expectedError string
	}{
		{
			name: "successful create",
			project: &sdk.Project{
				Name:        "New Project",
				Description: "Test project description",
			},
			setupMock: func(ms *MockStore) {
				ms.On("Create", mock.Anything, mock.MatchedBy(func(p *sdk.Project) bool {
					return p.Name == "New Project" && p.Description == "Test project description"
				})).Return(nil)
			},
		},
		{
			name: "store error",
			project: &sdk.Project{
				Name:        "New Project",
				Description: "Test project description",
			},
			setupMock: func(ms *MockStore) {
				ms.On("Create", mock.Anything, mock.Anything).Return(errors.New("creation failed"))
			},
			expectedError: "creation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}

			tt.setupMock(mockStore)

			service := NewService(mockStore)

			err := service.Create(context.Background(), tt.project)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name          string
		project       *sdk.Project
		setupMock     func(*MockStore)
		expectedError string
	}{
		{
			name: "successful update",
			project: &sdk.Project{
				Id:          "project1",
				Name:        "Updated Project",
				Description: "Updated project description",
			},
			setupMock: func(ms *MockStore) {
				ms.On("Update", mock.Anything, mock.MatchedBy(func(p *sdk.Project) bool {
					return p.Id == "project1" && p.Name == "Updated Project"
				})).Return(nil)
			},
		},
		{
			name: "store error",
			project: &sdk.Project{
				Id:          "project1",
				Name:        "Updated Project",
				Description: "Updated project description",
			},
			setupMock: func(ms *MockStore) {
				ms.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed"))
			},
			expectedError: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}

			tt.setupMock(mockStore)

			service := NewService(mockStore)

			err := service.Update(context.Background(), tt.project)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_ErrorConstants(t *testing.T) {
	// Test that the error constants are properly defined
	assert.NotNil(t, ErrProjectNotFound)
	assert.Equal(t, "project not found", ErrProjectNotFound.Error())
}

func TestService_IntegrationScenarios(t *testing.T) {
	t.Run("create and get project flow", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		project := &sdk.Project{
			Id:          "project1",
			Name:        "Integration Test Project",
			Description: "Test project for integration",
		}

		// Mock create operation
		mockStore.On("Create", mock.Anything, project).Return(nil)

		// Mock get operation
		mockStore.On("Get", mock.Anything, "project1").Return(project, nil)

		// Create project
		err := service.Create(context.Background(), project)
		assert.NoError(t, err)

		// Get project
		result, err := service.Get(context.Background(), "project1")
		assert.NoError(t, err)
		assert.Equal(t, project.Id, result.Id)
		assert.Equal(t, project.Name, result.Name)
		assert.Equal(t, project.Description, result.Description)

		mockStore.AssertExpectations(t)
	})

	t.Run("get by name after creation", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		project := &sdk.Project{
			Id:          "project2",
			Name:        "Named Project",
			Description: "Test project for name lookup",
		}

		// Mock create operation
		mockStore.On("Create", mock.Anything, project).Return(nil)

		// Mock get by name operation
		mockStore.On("GetByName", mock.Anything, "Named Project").Return(project, nil)

		// Create project
		err := service.Create(context.Background(), project)
		assert.NoError(t, err)

		// Get project by name
		result, err := service.GetByName(context.Background(), "Named Project")
		assert.NoError(t, err)
		assert.Equal(t, project.Id, result.Id)
		assert.Equal(t, project.Name, result.Name)

		mockStore.AssertExpectations(t)
	})
}
