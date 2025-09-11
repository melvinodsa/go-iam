package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
)

// MockDB implements the db.DB interface for testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) FindOne(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, col, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockDB) Find(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) InsertOne(ctx context.Context, col db.DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, col, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockDB) UpdateOne(ctx context.Context, col db.DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, col, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockDB) DeleteOne(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, col, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockDB) Aggregate(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) CountDocuments(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, col, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDB) BulkWrite(ctx context.Context, col db.DbCollection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := m.Called(ctx, col, models, opts)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (m *MockDB) UpdateMany(ctx context.Context, col db.DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, col, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockDB) SetDbInContext(ctx context.Context) context.Context {
	args := m.Called(ctx)
	return args.Get(0).(context.Context)
}

func (m *MockDB) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func createContextWithProjects() context.Context {
	metadata := sdk.Metadata{
		User: &sdk.User{
			Id:        "test-user",
			ProjectId: "project-123",
		},
		ProjectIds: []string{"project-123", "project-456"},
	}
	return middlewares.AddMetadata(context.Background(), metadata)
}

// TestNewStore tests the NewStore constructor
func TestNewStore(t *testing.T) {
	mockDB := &MockDB{}

	store := NewStore(mockDB)

	assert.NotNil(t, store)
}

// TestStoreCreate tests the Create method
func TestStoreCreate(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("success", func(t *testing.T) {
		user := &sdk.User{
			ProjectId: "project-123",
			Name:      "Test User",
			Email:     "test@example.com",
		}

		insertResult := &mongo.InsertOneResult{InsertedID: "user-123"}
		mockDB.On("InsertOne", ctx, mock.Anything, mock.Anything).Return(insertResult, nil)

		err := s.Create(ctx, user)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		user := &sdk.User{
			ProjectId: "project-123",
			Name:      "Test User",
		}

		mockDB.On("InsertOne", ctx, mock.Anything, mock.Anything).Return((*mongo.InsertOneResult)(nil), errors.New("database error"))

		err := s.Create(ctx, user)

		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})
}

// TestStoreUpdate tests the Update method
func TestStoreUpdate(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("empty_id_error", func(t *testing.T) {
		user := &sdk.User{
			ProjectId: "project-123",
			Name:      "Test User",
		}

		err := s.Update(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, ErrorUserNotFound, err)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		user := &sdk.User{
			Id:        "user-123",
			ProjectId: "project-123",
			Name:      "Test User",
		}

		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		err := s.Update(ctx, user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		mockDB.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		now := time.Now()
		user := &sdk.User{
			Id:        "user-123",
			ProjectId: "project-123",
			Name:      "Updated User",
		}

		existingUser := models.User{
			Id:        "user-123",
			ProjectId: "project-123",
			Name:      "Test User",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		userDoc, _ := bson.Marshal(existingUser)
		mockSingleResult := mongo.NewSingleResultFromDocument(userDoc, nil, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		updateResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}
		mockDB.On("UpdateOne", ctx, mock.Anything, mock.Anything, mock.Anything).Return(updateResult, nil)

		err := s.Update(ctx, user)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("update_database_error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		now := time.Now()
		user := &sdk.User{
			Id:        "user-123",
			ProjectId: "project-123",
			Name:      "Updated User",
		}

		// Mock GetById to return existing user successfully
		existingUser := models.User{
			Id:        "user-123",
			ProjectId: "project-123",
			Name:      "Test User",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		userDoc, _ := bson.Marshal(existingUser)
		mockSingleResult := mongo.NewSingleResultFromDocument(userDoc, nil, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		// Mock UpdateOne to return an error
		mockDB.On("UpdateOne", ctx, mock.Anything, mock.Anything, mock.Anything).Return((*mongo.UpdateResult)(nil), errors.New("database update error"))

		err := s.Update(ctx, user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating user")
		assert.Contains(t, err.Error(), "database update error")
		mockDB.AssertExpectations(t)
	})
}

// TestStoreGetById tests the GetById method
func TestStoreGetById(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("success", func(t *testing.T) {
		userID := "user-123"
		now := time.Now()
		user := models.User{
			Id:        userID,
			ProjectId: "project-123",
			Name:      "Test User",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		userDoc, _ := bson.Marshal(user)
		mockSingleResult := mongo.NewSingleResultFromDocument(userDoc, nil, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetById(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.Id)
		mockDB.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		userID := "nonexistent-id"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetById(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, ErrorUserNotFound, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		userID := "user-123"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("database connection error"), nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetById(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		assert.Contains(t, err.Error(), "database connection error")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

// TestStoreGetByEmail tests the GetByEmail method
func TestStoreGetByEmail(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		projectId := "project-123"
		now := time.Now()
		user := models.User{
			Id:        "user-123",
			ProjectId: projectId,
			Email:     email,
			Name:      "Test User",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		userDoc, _ := bson.Marshal(user)
		mockSingleResult := mongo.NewSingleResultFromDocument(userDoc, nil, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByEmail(ctx, email, projectId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, email, result.Email)
		mockDB.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		email := "nonexistent@example.com"
		projectId := "project-123"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByEmail(ctx, email, projectId)

		assert.Error(t, err)
		assert.Equal(t, ErrorUserNotFound, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		email := "test@example.com"
		projectId := "project-123"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("database connection error"), nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByEmail(ctx, email, projectId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		assert.Contains(t, err.Error(), "database connection error")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

// TestStoreGetByPhone tests the GetByPhone method
func TestStoreGetByPhone(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("success", func(t *testing.T) {
		phone := "+1234567890"
		projectId := "project-123"
		now := time.Now()
		user := models.User{
			Id:        "user-123",
			ProjectId: projectId,
			Phone:     phone,
			Name:      "Test User",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		userDoc, _ := bson.Marshal(user)
		mockSingleResult := mongo.NewSingleResultFromDocument(userDoc, nil, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByPhone(ctx, phone, projectId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, phone, result.Phone)
		mockDB.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		phone := "+9999999999"
		projectId := "project-123"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByPhone(ctx, phone, projectId)

		assert.Error(t, err)
		assert.Equal(t, ErrorUserNotFound, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		phone := "+1234567890"
		projectId := "project-123"
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("database connection error"), nil)
		mockDB.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

		result, err := s.GetByPhone(ctx, phone, projectId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		assert.Contains(t, err.Error(), "database connection error")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

// TestStoreGetAll tests the GetAll method
func TestStoreGetAll(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("count_error", func(t *testing.T) {
		query := sdk.UserQuery{
			Skip:  0,
			Limit: 10,
		}

		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(0), errors.New("database error"))

		result, err := s.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error counting users")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("find_error", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		query := sdk.UserQuery{
			Skip:  0,
			Limit: 10,
		}

		// Mock CountDocuments to succeed
		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(5), nil)

		// Mock Find to fail
		mockDB.On("Find", ctx, mock.Anything, mock.Anything, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("database find error"))

		result, err := s.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding all users")
		assert.Contains(t, err.Error(), "database find error")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("query_with_search_and_role", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		query := sdk.UserQuery{
			Skip:        0,
			Limit:       10,
			SearchQuery: "test",
			RoleId:      "admin-role",
		}

		// This test will fail at CountDocuments, but it will exercise the query building logic
		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(0), errors.New("expected error"))

		result, err := s.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error counting users")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty_search_query", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		query := sdk.UserQuery{
			Skip:        0,
			Limit:       10,
			SearchQuery: "", // Empty search query to test different path
			RoleId:      "",
		}

		// This will test the path where no search filter is added
		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(0), errors.New("expected error"))

		result, err := s.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error counting users")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	// This test focuses on exercising the query building logic without cursor operations
	t.Run("complex_query_building", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		query := sdk.UserQuery{
			Skip:        5,
			Limit:       20,
			SearchQuery: "john",
			RoleId:      "admin",
		}

		// Mock CountDocuments to succeed, then fail Find to avoid cursor issues
		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(25), nil)
		mockDB.On("Find", ctx, mock.Anything, mock.Anything, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("find error"))

		result, err := s.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding all users")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("success_with_cursor", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		query := sdk.UserQuery{
			Skip:  0,
			Limit: 10,
		}

		// Mock CountDocuments to succeed
		mockDB.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(2), nil)

		// Create mock users for the cursor
		now := time.Now()
		mockUsers := []interface{}{
			models.User{
				Id:        "user1",
				ProjectId: "project-123",
				Name:      "User One",
				Email:     "user1@example.com",
				Enabled:   true,
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			models.User{
				Id:        "user2",
				ProjectId: "project-123",
				Name:      "User Two",
				Email:     "user2@example.com",
				Enabled:   true,
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		}

		// Create a real cursor using NewCursorFromDocuments
		cursor, _ := mongo.NewCursorFromDocuments(mockUsers, nil, nil)
		mockDB.On("Find", ctx, mock.Anything, mock.Anything, mock.Anything).Return(cursor, nil)

		result, err := s.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Users, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, int64(0), result.Skip)
		assert.Equal(t, int64(10), result.Limit)
		mockDB.AssertExpectations(t)
	})
}

// TestStoreRemoveResourceFromAll tests the RemoveResourceFromAll method
func TestStoreRemoveResourceFromAll(t *testing.T) {
	ctx := createContextWithProjects()
	mockDB := &MockDB{}
	s := NewStore(mockDB)

	t.Run("update_many_error", func(t *testing.T) {
		resourceKey := "resource-123"

		mockDB.On("UpdateMany", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*mongo.UpdateResult)(nil), errors.New("database error"))

		err := s.RemoveResourceFromAll(ctx, resourceKey)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error removing resource from all users")
		mockDB.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockDB.ExpectedCalls = nil
		resourceKey := "resource-123"

		updateResult := &mongo.UpdateResult{MatchedCount: 5, ModifiedCount: 5}
		mockDB.On("UpdateMany", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(updateResult, nil)

		err := s.RemoveResourceFromAll(ctx, resourceKey)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
