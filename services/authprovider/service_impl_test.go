package authprovider

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore implements Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sdk.AuthProvider), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*sdk.AuthProvider, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.AuthProvider), args.Error(1)
}

func (m *MockStore) Create(ctx context.Context, provider *sdk.AuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, provider *sdk.AuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

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

func (m *MockProjectService) Create(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) Update(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name     string
		store    Store
		project  project.Service
		expected Service
	}{
		{
			name:     "success_create_new_service",
			store:    &MockStore{},
			project:  &MockProjectService{},
			expected: &service{s: &MockStore{}, p: &MockProjectService{}},
		},
		{
			name:     "success_with_nil_dependencies",
			store:    nil,
			project:  nil,
			expected: &service{s: nil, p: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewService(tt.store, tt.project)

			// Check that the service is not nil
			assert.NotNil(t, result)

			// Check that it's the correct type
			svc, ok := result.(*service)
			assert.True(t, ok)

			// Check internal fields are set correctly
			assert.Equal(t, tt.store, svc.s)
			assert.Equal(t, tt.project, svc.p)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		contextSetup   func() context.Context
		params         sdk.AuthProviderQueryParams
		mockSetup      func(*MockStore)
		expectedResult []sdk.AuthProvider
		expectedError  error
	}{
		{
			name: "success_with_projects_in_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1", "project2"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			params: sdk.AuthProviderQueryParams{},
			mockSetup: func(m *MockStore) {
				expectedParams := sdk.AuthProviderQueryParams{
					ProjectIds: []string{"project1", "project2"},
				}
				m.On("GetAll", mock.Anything, expectedParams).Return([]sdk.AuthProvider{
					{Id: "ap1", Name: "Provider 1", ProjectId: "project1"},
					{Id: "ap2", Name: "Provider 2", ProjectId: "project2"},
				}, nil)
			},
			expectedResult: []sdk.AuthProvider{
				{Id: "ap1", Name: "Provider 1", ProjectId: "project1"},
				{Id: "ap2", Name: "Provider 2", ProjectId: "project2"},
			},
			expectedError: nil,
		},
		{
			name: "success_with_empty_projects",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			params: sdk.AuthProviderQueryParams{},
			mockSetup: func(m *MockStore) {
				expectedParams := sdk.AuthProviderQueryParams{
					ProjectIds: []string{},
				}
				m.On("GetAll", mock.Anything, expectedParams).Return([]sdk.AuthProvider{}, nil)
			},
			expectedResult: []sdk.AuthProvider{},
			expectedError:  nil,
		},
		{
			name: "success_with_no_context_projects",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: nil, // No projects in context
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			params: sdk.AuthProviderQueryParams{},
			mockSetup: func(m *MockStore) {
				expectedParams := sdk.AuthProviderQueryParams{
					ProjectIds: nil,
				}
				m.On("GetAll", mock.Anything, expectedParams).Return([]sdk.AuthProvider{}, nil)
			},
			expectedResult: []sdk.AuthProvider{},
			expectedError:  nil,
		},
		{
			name: "error_from_store",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			params: sdk.AuthProviderQueryParams{},
			mockSetup: func(m *MockStore) {
				expectedParams := sdk.AuthProviderQueryParams{
					ProjectIds: []string{"project1"},
				}
				m.On("GetAll", mock.Anything, expectedParams).Return([]sdk.AuthProvider{}, errors.New("store error"))
			},
			expectedResult: []sdk.AuthProvider{},
			expectedError:  errors.New("store error"),
		},
		{
			name: "success_preserves_existing_params",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			params: sdk.AuthProviderQueryParams{
				ProjectIds: []string{"existing"},
			},
			mockSetup: func(m *MockStore) {
				expectedParams := sdk.AuthProviderQueryParams{
					ProjectIds: []string{"project1"}, // Should override existing
				}
				m.On("GetAll", mock.Anything, expectedParams).Return([]sdk.AuthProvider{
					{Id: "ap1", Name: "Provider 1", ProjectId: "project1"},
				}, nil)
			},
			expectedResult: []sdk.AuthProvider{
				{Id: "ap1", Name: "Provider 1", ProjectId: "project1"},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			mockProject := &MockProjectService{}
			svc := NewService(mockStore, mockProject)

			ctx := tt.contextSetup()
			tt.mockSetup(mockStore)

			result, err := svc.GetAll(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name              string
		contextSetup      func() context.Context
		id                string
		dontCheckProjects bool
		mockSetup         func(*MockStore)
		expectedResult    *sdk.AuthProvider
		expectedError     error
	}{
		{
			name: "success_dont_check_projects",
			contextSetup: func() context.Context {
				// For dontCheckProjects=true, we still need proper metadata setup
				metadata := sdk.Metadata{
					ProjectIds: []string{}, // Empty projects is fine since we're not checking
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: true,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{
					Id:        "ap1",
					Name:      "Provider 1",
					ProjectId: "project1",
				}, nil)
			},
			expectedResult: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			expectedError: nil,
		},
		{
			name: "success_check_projects_allowed",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1", "project2"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: false,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{
					Id:        "ap1",
					Name:      "Provider 1",
					ProjectId: "project1",
				}, nil)
			},
			expectedResult: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			expectedError: nil,
		},
		{
			name: "error_check_projects_not_allowed",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project2", "project3"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: false,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{
					Id:        "ap1",
					Name:      "Provider 1",
					ProjectId: "project1",
				}, nil)
			},
			expectedResult: nil,
			expectedError:  ErrAuthProviderNotFound,
		},
		{
			name: "error_from_store",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{}, // Empty projects for store error test
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: true,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{}, errors.New("store error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("store error"),
		},
		{
			name: "error_with_empty_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: false,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{
					Id:        "ap1",
					Name:      "Provider 1",
					ProjectId: "project1",
				}, nil)
			},
			expectedResult: nil,
			expectedError:  ErrAuthProviderNotFound,
		},
		{
			name: "error_with_nil_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: nil, // No projects in context
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			id:                "ap1",
			dontCheckProjects: false,
			mockSetup: func(m *MockStore) {
				m.On("Get", mock.Anything, "ap1").Return(&sdk.AuthProvider{
					Id:        "ap1",
					Name:      "Provider 1",
					ProjectId: "project1",
				}, nil)
			},
			expectedResult: nil,
			expectedError:  ErrAuthProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			mockProject := &MockProjectService{}
			svc := NewService(mockStore, mockProject)

			ctx := tt.contextSetup()
			tt.mockSetup(mockStore)

			result, err := svc.Get(ctx, tt.id, tt.dontCheckProjects)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		provider      *sdk.AuthProvider
		mockSetup     func(*MockStore)
		expectedError error
	}{
		{
			name: "success_valid_project",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1", "project2"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			mockSetup: func(m *MockStore) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(p *sdk.AuthProvider) bool {
					return p.Id == "ap1" && p.Name == "Provider 1" && p.ProjectId == "project1"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error_invalid_project",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project2", "project3"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_empty_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_nil_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: nil, // No projects in context for Create
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_from_store",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1",
				ProjectId: "project1",
			},
			mockSetup: func(m *MockStore) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(p *sdk.AuthProvider) bool {
					return p.Id == "ap1" && p.Name == "Provider 1" && p.ProjectId == "project1"
				})).Return(errors.New("store error"))
			},
			expectedError: errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			mockProject := &MockProjectService{}
			svc := NewService(mockStore, mockProject)

			ctx := tt.contextSetup()
			tt.mockSetup(mockStore)

			err := svc.Create(ctx, tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
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
		contextSetup  func() context.Context
		provider      *sdk.AuthProvider
		mockSetup     func(*MockStore)
		expectedError error
	}{
		{
			name: "success_valid_project",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1", "project2"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1 Updated",
				ProjectId: "project1",
			},
			mockSetup: func(m *MockStore) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(p *sdk.AuthProvider) bool {
					return p.Id == "ap1" && p.Name == "Provider 1 Updated" && p.ProjectId == "project1"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error_invalid_project",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project2", "project3"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1 Updated",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_empty_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1 Updated",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_nil_projects_context",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: nil, // No projects in context for Update
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1 Updated",
				ProjectId: "project1",
			},
			mockSetup:     func(m *MockStore) {},
			expectedError: project.ErrProjectNotFound,
		},
		{
			name: "error_from_store",
			contextSetup: func() context.Context {
				metadata := sdk.Metadata{
					ProjectIds: []string{"project1"},
				}
				return middlewares.AddMetadata(context.Background(), metadata)
			},
			provider: &sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Provider 1 Updated",
				ProjectId: "project1",
			},
			mockSetup: func(m *MockStore) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(p *sdk.AuthProvider) bool {
					return p.Id == "ap1" && p.Name == "Provider 1 Updated" && p.ProjectId == "project1"
				})).Return(errors.New("store error"))
			},
			expectedError: errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			mockProject := &MockProjectService{}
			svc := NewService(mockStore, mockProject)

			ctx := tt.contextSetup()
			tt.mockSetup(mockStore)

			err := svc.Update(ctx, tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestService_GetProvider(t *testing.T) {
	tests := []struct {
		name           string
		authProvider   sdk.AuthProvider
		expectedResult sdk.ServiceProvider
		expectedError  error
	}{
		{
			name: "success_google_provider",
			authProvider: sdk.AuthProvider{
				Id:       "ap1",
				Name:     "Google Provider",
				Provider: sdk.AuthProviderTypeGoogle,
				Params: []sdk.AuthProviderParam{
					{Key: "@GOOGLE/CLIENT_ID", Value: "client123"},
					{Key: "@GOOGLE/CLIENT_SECRET", Value: "secret123"},
					{Key: "@GOOGLE/REDIRECT_URL", Value: "http://localhost:8080/callback"},
				},
				ProjectId: "project1",
			},
			expectedResult: nil, // We can't easily compare the Google provider instance
			expectedError:  nil,
		},
		{
			name: "error_unknown_provider",
			authProvider: sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Unknown Provider",
				Provider:  "UNKNOWN",
				ProjectId: "project1",
			},
			expectedResult: nil,
			expectedError:  errors.New("unknown auth provider: UNKNOWN"),
		},
		{
			name: "success_google_provider_empty_params",
			authProvider: sdk.AuthProvider{
				Id:        "ap1",
				Name:      "Google Provider",
				Provider:  sdk.AuthProviderTypeGoogle,
				Params:    []sdk.AuthProviderParam{},
				ProjectId: "project1",
			},
			expectedResult: nil, // We can't easily compare the Google provider instance
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStore{}
			mockProject := &MockProjectService{}
			svc := NewService(mockStore, mockProject)

			result, err := svc.GetProvider(context.Background(), tt.authProvider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				// For Google provider, we just check that it's not nil
				// since comparing the actual instance is complex
				if tt.authProvider.Provider == sdk.AuthProviderTypeGoogle {
					assert.NotNil(t, result)
				} else {
					assert.Equal(t, tt.expectedResult, result)
				}
			}
		})
	}
}
