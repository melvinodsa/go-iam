package role

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewStore(t *testing.T) {
	mockDB := test.SetupMockDB()

	store := NewStore(mockDB)

	assert.NotNil(t, store)
	assert.Implements(t, (*Store)(nil), store)
}

func TestStore_Create(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Resources: map[string]sdk.Resources{
				"users": {
					Id:   "resource1",
					Key:  "users",
					Name: "Users Resource",
				},
			},
			Enabled: true,
		}

		mockDB.On("InsertOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("models.Role"), mock.Anything).Return(&mongo.InsertOneResult{}, nil)

		err := store.Create(ctx, role)

		assert.NoError(t, err)
		assert.NotEmpty(t, role.Id)      // ID should be generated
		assert.NotNil(t, role.CreatedAt) // CreatedAt should be set
		assert.NotNil(t, role.UpdatedAt) // UpdatedAt should be set
		mockDB.AssertExpectations(t)
	})

	t.Run("nil_role", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		err := store.Create(ctx, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role cannot be nil")
		mockDB.AssertNotCalled(t, "InsertOne")
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockDB.On("InsertOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("models.Role"), mock.Anything).Return(&mongo.InsertOneResult{}, errors.New("database error"))

		err := store.Create(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create role")
		assert.Contains(t, err.Error(), "database error")
		mockDB.AssertExpectations(t)
	})
}

func TestStore_Update(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Updated Role",
			Description: "An updated role",
			ProjectId:   "project1",
			Resources: map[string]sdk.Resources{
				"users": {
					Id:   "resource1",
					Key:  "users",
					Name: "Users Resource",
				},
			},
			Enabled: true,
		}

		mockDB.On("UpdateOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(&mongo.UpdateResult{ModifiedCount: 1}, nil)

		err := store.Update(ctx, role)

		assert.NoError(t, err)
		assert.NotNil(t, role.UpdatedAt) // UpdatedAt should be set
		mockDB.AssertExpectations(t)
	})

	t.Run("nil_role", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		err := store.Update(ctx, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role ID is required")
		mockDB.AssertNotCalled(t, "UpdateOne")
	})

	t.Run("empty_role_id", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		err := store.Update(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role ID is required")
		mockDB.AssertNotCalled(t, "UpdateOne")
	})

	t.Run("role_not_found", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "nonexistent",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockDB.On("UpdateOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(&mongo.UpdateResult{ModifiedCount: 0}, nil)

		err := store.Update(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockDB.On("UpdateOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(&mongo.UpdateResult{}, errors.New("database error"))

		err := store.Update(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update role")
		assert.Contains(t, err.Error(), "database error")
		mockDB.AssertExpectations(t)
	})
}

func TestStore_GetById(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		createdAt := time.Date(2025, time.September, 1, 8, 15, 42, 133000000, time.UTC)
		updatedAt := time.Date(2025, time.September, 1, 9, 15, 42, 133000000, time.UTC)

		// Create a single result from the role document
		roleDoc := bson.D{
			{Key: "id", Value: "role1"},
			{Key: "name", Value: "Test Role"},
			{Key: "description", Value: "A test role"},
			{Key: "project_id", Value: "project1"},
			{Key: "resources", Value: map[string]models.Resources{
				"users": {
					Id:   "resource1",
					Key:  "users",
					Name: "Users Resource",
				},
			}},
			{Key: "enabled", Value: true},
			{Key: "created_at", Value: createdAt},
			{Key: "created_by", Value: "user1"},
			{Key: "updated_at", Value: updatedAt},
			{Key: "updated_by", Value: "user2"},
		}
		mockResult := mongo.NewSingleResultFromDocument(roleDoc, nil, nil)

		mockDB.On("FindOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(mockResult)

		result, err := store.GetById(ctx, "role1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "role1", result.Id)
		assert.Equal(t, "Test Role", result.Name)
		assert.Equal(t, "A test role", result.Description)
		assert.Equal(t, "project1", result.ProjectId)
		assert.True(t, result.Enabled)
		assert.Equal(t, 1, len(result.Resources))
		assert.Contains(t, result.Resources, "users")
		assert.Equal(t, "Users Resource", result.Resources["users"].Name)
		assert.Equal(t, &createdAt, result.CreatedAt)
		assert.Equal(t, "user1", result.CreatedBy)
		assert.Equal(t, &updatedAt, result.UpdatedAt)
		assert.Equal(t, "user2", result.UpdatedBy)

		mockDB.AssertExpectations(t)
	})

	t.Run("empty_id", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		result, err := store.GetById(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "role ID cannot be empty")
		mockDB.AssertNotCalled(t, "FindOne")
	})

	t.Run("role_not_found", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		mockResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
		mockDB.On("FindOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(mockResult)

		result, err := store.GetById(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "role with ID nonexistent not found")

		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		mockResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("database error"), nil)
		mockDB.On("FindOne", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(mockResult)

		result, err := store.GetById(ctx, "role1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find role")
		assert.Contains(t, err.Error(), "database error")

		mockDB.AssertExpectations(t)
	})
}

func TestStore_GetAll(t *testing.T) {
	t.Run("successful_get_all_with_search", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		query := sdk.RoleQuery{
			ProjectIds:  []string{"project1", "project2"},
			SearchQuery: "test",
			Skip:        0,
			Limit:       10,
		}

		// Create cursor from role documents
		roleDocuments := []interface{}{
			bson.D{
				{Key: "id", Value: "role1"},
				{Key: "name", Value: "Test Role 1"},
				{Key: "description", Value: "A test role"},
				{Key: "project_id", Value: "project1"},
				{Key: "enabled", Value: true},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
			bson.D{
				{Key: "id", Value: "role2"},
				{Key: "name", Value: "Test Role 2"},
				{Key: "description", Value: "Another test role"},
				{Key: "project_id", Value: "project2"},
				{Key: "enabled", Value: true},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
		}
		cursor, _ := mongo.NewCursorFromDocuments(roleDocuments, nil, nil)

		mockDB.On("CountDocuments", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(int64(2), nil)
		mockDB.On("Find", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(cursor, nil)

		result, err := store.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, int64(0), result.Skip)
		assert.Equal(t, int64(10), result.Limit)
		assert.Equal(t, 2, len(result.Roles))
		assert.Equal(t, "role1", result.Roles[0].Id)
		assert.Equal(t, "Test Role 1", result.Roles[0].Name)
		assert.Equal(t, "role2", result.Roles[1].Id)
		assert.Equal(t, "Test Role 2", result.Roles[1].Name)

		mockDB.AssertExpectations(t)
	})

	t.Run("successful_get_all_without_search", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		query := sdk.RoleQuery{
			ProjectIds:  []string{"project1"},
			SearchQuery: "",
			Skip:        0,
			Limit:       5,
		}

		// Create cursor from role documents
		roleDocuments := []interface{}{
			bson.D{
				{Key: "id", Value: "role1"},
				{Key: "name", Value: "Role 1"},
				{Key: "description", Value: "First role"},
				{Key: "project_id", Value: "project1"},
				{Key: "enabled", Value: true},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
		}
		cursor, _ := mongo.NewCursorFromDocuments(roleDocuments, nil, nil)

		mockDB.On("CountDocuments", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(int64(1), nil)
		mockDB.On("Find", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(cursor, nil)

		result, err := store.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Total)
		assert.Equal(t, int64(0), result.Skip)
		assert.Equal(t, int64(5), result.Limit)
		assert.Equal(t, 1, len(result.Roles))
		assert.Equal(t, "role1", result.Roles[0].Id)

		mockDB.AssertExpectations(t)
	})

	t.Run("empty_result", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		query := sdk.RoleQuery{
			ProjectIds:  []string{"project1"},
			SearchQuery: "nonexistent",
			Skip:        0,
			Limit:       10,
		}

		// Create empty cursor
		cursor, _ := mongo.NewCursorFromDocuments([]interface{}{}, nil, nil)

		mockDB.On("CountDocuments", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(int64(0), nil)
		mockDB.On("Find", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(cursor, nil)

		result, err := store.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Total)
		assert.Equal(t, 0, len(result.Roles))

		mockDB.AssertExpectations(t)
	})

	t.Run("count_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		query := sdk.RoleQuery{
			ProjectIds:  []string{"project1"},
			SearchQuery: "test",
			Skip:        0,
			Limit:       10,
		}

		mockDB.On("CountDocuments", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(int64(0), errors.New("count error"))

		result, err := store.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error counting roles")
		assert.Contains(t, err.Error(), "count error")

		mockDB.AssertExpectations(t)
	})

	t.Run("find_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		query := sdk.RoleQuery{
			ProjectIds:  []string{"project1"},
			SearchQuery: "test",
			Skip:        0,
			Limit:       10,
		}

		mockDB.On("CountDocuments", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(int64(2), nil)
		mockDB.On("Find", ctx, mock.AnythingOfType("models.RoleModel"), mock.AnythingOfType("primitive.D"), mock.Anything).Return(&mongo.Cursor{}, errors.New("find error"))

		result, err := store.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to fetch roles")
		assert.Contains(t, err.Error(), "find error")

		mockDB.AssertExpectations(t)
	})
}
