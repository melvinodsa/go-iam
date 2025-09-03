package policy

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

func (m *MockStore) GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.PolicyList), args.Error(1)
}

func TestNewService(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	assert.NotNil(t, svc)
	// Since the service struct is not exported, we test the behavior instead
	assert.Implements(t, (*Service)(nil), svc)
}

func TestService_GetAll_Success(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	ctx := context.Background()
	query := sdk.PolicyQuery{
		Query: "test",
		Limit: 10,
		Skip:  0,
	}

	expectedPolicies := &sdk.PolicyList{
		Policies: []sdk.Policy{
			{
				Id:          "policy1",
				Name:        "Test Policy 1",
				Description: "First test policy",
			},
			{
				Id:          "policy2",
				Name:        "Test Policy 2",
				Description: "Second test policy",
			},
		},
		Total: 2,
		Limit: 10,
		Skip:  0,
	}

	// Setup mock
	store.On("GetAll", ctx, query).Return(expectedPolicies, nil)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedPolicies, result)
	store.AssertExpectations(t)
}

func TestService_GetAll_StoreError(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	ctx := context.Background()
	query := sdk.PolicyQuery{
		Query: "test",
		Limit: 5,
		Skip:  0,
	}

	expectedError := errors.New("store error")

	// Setup mock
	store.On("GetAll", ctx, query).Return(nil, expectedError)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	store.AssertExpectations(t)
}

func TestService_GetAll_EmptyResult(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	ctx := context.Background()
	query := sdk.PolicyQuery{
		Query: "nonexistent",
		Limit: 10,
		Skip:  0,
	}

	expectedPolicies := &sdk.PolicyList{
		Policies: []sdk.Policy{},
		Total:    0,
		Limit:    10,
		Skip:     0,
	}

	// Setup mock
	store.On("GetAll", ctx, query).Return(expectedPolicies, nil)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedPolicies, result)
	assert.Empty(t, result.Policies)
	assert.Equal(t, 0, result.Total)
	store.AssertExpectations(t)
}

func TestService_GetAll_WithPagination(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	ctx := context.Background()
	query := sdk.PolicyQuery{
		Query: "test",
		Limit: 2,
		Skip:  5,
	}

	expectedPolicies := &sdk.PolicyList{
		Policies: []sdk.Policy{
			{
				Id:          "policy6",
				Name:        "Test Policy 6",
				Description: "Sixth test policy",
			},
			{
				Id:          "policy7",
				Name:        "Test Policy 7",
				Description: "Seventh test policy",
			},
		},
		Total: 20, // Total count of all policies
		Limit: 2,
		Skip:  5,
	}

	// Setup mock
	store.On("GetAll", ctx, query).Return(expectedPolicies, nil)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedPolicies, result)
	assert.Len(t, result.Policies, 2)
	assert.Equal(t, 20, result.Total)
	assert.Equal(t, int64(2), result.Limit)
	assert.Equal(t, int64(5), result.Skip)
	store.AssertExpectations(t)
}

func TestService_GetAll_WithQuery(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	ctx := context.Background()
	query := sdk.PolicyQuery{
		Query: "specific search term",
		Limit: 50,
		Skip:  0,
	}

	expectedPolicies := &sdk.PolicyList{
		Policies: []sdk.Policy{
			{
				Id:          "policy1",
				Name:        "Specific Policy",
				Description: "Policy matching the search term",
			},
			{
				Id:          "policy2",
				Name:        "Another Specific Policy",
				Description: "Another policy matching the search term",
			},
		},
		Total: 2,
		Limit: 50,
		Skip:  0,
	}

	// Setup mock
	store.On("GetAll", ctx, query).Return(expectedPolicies, nil)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedPolicies, result)
	assert.Len(t, result.Policies, 2)
	store.AssertExpectations(t)
}

func TestService_GetAll_ContextCancellation(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	query := sdk.PolicyQuery{
		Query: "test",
		Limit: 10,
		Skip:  0,
	}

	expectedError := context.Canceled

	// Setup mock
	store.On("GetAll", ctx, query).Return(nil, expectedError)

	// Execute
	result, err := svc.GetAll(ctx, query)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	store.AssertExpectations(t)
}

func TestService_GetAll_BusinessLogic(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	t.Run("service_delegates_to_store", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{Query: "test"}

		store.On("GetAll", ctx, query).Return(&sdk.PolicyList{}, nil).Once()

		_, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		store.AssertExpectations(t)
	})

	t.Run("service_passes_exact_parameters", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{
			Query: "exact-query",
			Limit: 42,
			Skip:  7,
		}

		// Verify exact parameters are passed through
		store.On("GetAll",
			mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
			mock.MatchedBy(func(q sdk.PolicyQuery) bool {
				return q.Query == "exact-query" &&
					q.Limit == 42 &&
					q.Skip == 7
			}),
		).Return(&sdk.PolicyList{}, nil).Once()

		_, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		store.AssertExpectations(t)
	})

	t.Run("service_returns_exact_store_result", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{Query: "test"}

		expectedResult := &sdk.PolicyList{
			Policies: []sdk.Policy{{Id: "test-policy"}},
			Total:    1,
		}

		store.On("GetAll", ctx, query).Return(expectedResult, nil).Once()

		result, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		store.AssertExpectations(t)
	})

	t.Run("service_returns_exact_store_error", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{Query: "test"}

		expectedError := errors.New("specific store error")

		store.On("GetAll", ctx, query).Return(nil, expectedError).Once()

		result, err := svc.GetAll(ctx, query)

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		store.AssertExpectations(t)
	})
}

func TestService_GetAll_EdgeCases(t *testing.T) {
	store := &MockStore{}
	svc := NewService(store)

	t.Run("empty_query", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{
			Query: "",
			Limit: 10,
			Skip:  0,
		}

		expectedResult := &sdk.PolicyList{
			Policies: []sdk.Policy{},
			Total:    0,
			Limit:    10,
			Skip:     0,
		}

		store.On("GetAll", ctx, query).Return(expectedResult, nil).Once()

		result, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		store.AssertExpectations(t)
	})

	t.Run("zero_limit", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{
			Query: "test",
			Limit: 0,
			Skip:  0,
		}

		expectedResult := &sdk.PolicyList{
			Policies: []sdk.Policy{},
			Total:    0,
			Limit:    0,
			Skip:     0,
		}

		store.On("GetAll", ctx, query).Return(expectedResult, nil).Once()

		result, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		store.AssertExpectations(t)
	})

	t.Run("large_skip_value", func(t *testing.T) {
		ctx := context.Background()
		query := sdk.PolicyQuery{
			Query: "test",
			Limit: 10,
			Skip:  1000,
		}

		expectedResult := &sdk.PolicyList{
			Policies: []sdk.Policy{},
			Total:    50,
			Limit:    10,
			Skip:     1000,
		}

		store.On("GetAll", ctx, query).Return(expectedResult, nil).Once()

		result, err := svc.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.Empty(t, result.Policies) // No results beyond total
		store.AssertExpectations(t)
	})
}
