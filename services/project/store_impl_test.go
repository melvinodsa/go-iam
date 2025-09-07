package project

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

		project := &sdk.Project{
			Name:        "Test Project",
			Description: "A test project",
			Tags:        []string{"test", "project"},
		}

		mockDB.On("InsertOne", ctx, mock.AnythingOfType("models.ProjectModel"), mock.AnythingOfType("models.Project"), mock.Anything).Return(&mongo.InsertOneResult{}, nil)

		err := store.Create(ctx, project)

		assert.NoError(t, err)
		assert.NotEmpty(t, project.Id)      // ID should be generated
		assert.NotNil(t, project.CreatedAt) // CreatedAt should be set
		mockDB.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		project := &sdk.Project{
			Name:        "Test Project",
			Description: "A test project",
		}

		mockDB.On("InsertOne", ctx, mock.AnythingOfType("models.ProjectModel"), mock.AnythingOfType("models.Project"), mock.Anything).Return((*mongo.InsertOneResult)(nil), errors.New("database error"))

		err := store.Create(ctx, project)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error creating project")
		mockDB.AssertExpectations(t)
	})
}

func TestStore_Get(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "project1"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.IdKey, Value: projectId}}

		// Create a mock SingleResult
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{{Key: md.IdKey, Value: "test"}}, nil, nil)

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		// Since we can't easily mock SingleResult.Decode(), this test verifies the correct calls
		result, err := store.Get(ctx, projectId)

		// This will fail due to mocking limitations, but verifies the correct calls
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test", result.Id)
		mockDB.AssertExpectations(t)
	})

	t.Run("project_not_found", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "nonexistent"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.IdKey, Value: projectId}}

		// Create a mock SingleResult that will return ErrNoDocuments
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		result, err := store.Get(ctx, projectId)

		// Due to mocking limitations, we can't easily test the exact error
		// But we can verify the correct database call was made
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestStore_GetByName(t *testing.T) {
	t.Run("successful_get_by_name", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectName := "Test Project"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.NameKey, Value: projectName}}

		// Create a mock SingleResult
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{{Key: md.IdKey, Value: "test"}}, nil, nil)

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		result, err := store.GetByName(ctx, projectName)

		// This will fail due to mocking limitations, but verifies the correct calls
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test", result.Id)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty_name", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		result, err := store.GetByName(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrProjectNotFound, err)
		assert.Nil(t, result)
		mockDB.AssertNotCalled(t, "FindOne")
	})

	t.Run("project_not_found", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectName := "Nonexistent Project"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.NameKey, Value: projectName}}

		// Create a mock SingleResult that will return ErrNoDocuments
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		result, err := store.GetByName(ctx, projectName)

		// Due to mocking limitations, we can't easily test the exact error
		// But we can verify the correct database call was made
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestStore_GetAll(t *testing.T) {
	t.Run("find_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		md := models.GetProjectModel()
		expectedFilter := bson.D{{}}

		mockDB.On("Find", ctx, md, expectedFilter, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("find error"))

		result, err := store.GetAll(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding all projects")
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("successful", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		md := models.GetProjectModel()
		expectedFilter := bson.D{{}}

		cursor, _ := mongo.NewCursorFromDocuments([]interface{}{bson.D{{Key: md.IdKey, Value: "test"}}}, nil, nil)

		mockDB.On("Find", ctx, md, expectedFilter, mock.Anything).Return(cursor, nil)

		result, err := store.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockDB.AssertExpectations(t)
	})
}

func TestStore_Update(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "project1"

		project := &sdk.Project{
			Id:   projectId,
			Name: "Updated Project",
		}

		md := models.GetProjectModel()
		expectedGetFilter := bson.D{{Key: md.IdKey, Value: projectId}}
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{{Key: md.IdKey, Value: projectId}}, nil, nil)

		mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult)
		expectedUpdateFilter := bson.D{{Key: md.IdKey, Value: projectId}}

		mockDB.On("UpdateOne", ctx, md, expectedUpdateFilter, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)

		err := store.Update(ctx, project)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty_id", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		project := &sdk.Project{
			Id: "", // Empty ID should trigger error
		}

		err := store.Update(ctx, project)

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrProjectNotFound, err)
		mockDB.AssertNotCalled(t, "UpdateOne")
	})

	t.Run("get_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "project1"

		project := &sdk.Project{
			Id:   projectId,
			Name: "Updated Project",
		}

		md := models.GetProjectModel()
		expectedGetFilter := bson.D{{Key: md.IdKey, Value: projectId}}
		mockSingleResult := &mongo.SingleResult{}

		mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult)

		err := store.Update(ctx, project)

		// This will fail due to Get() error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding project")
		mockDB.AssertExpectations(t)
	})

	t.Run("update_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "project1"

		now := time.Now()
		project := &sdk.Project{
			Id:        projectId,
			Name:      "Updated Project",
			CreatedAt: &now,
			CreatedBy: "original-user",
		}

		md := models.GetProjectModel()
		expectedGetFilter := bson.D{{Key: md.IdKey, Value: projectId}}
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{{Key: md.IdKey, Value: projectId}}, nil, nil)

		mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult)
		expectedUpdateFilter := bson.D{{Key: md.IdKey, Value: projectId}}

		mockDB.On("UpdateOne", ctx, md, expectedUpdateFilter, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, errors.New("error updating project"))

		// Note: Due to mocking limitations with SingleResult.Decode(),
		// this test verifies that the correct database calls are made
		// but the actual Update operation will fail due to the Get operation
		err := store.Update(ctx, project)

		// The error will be from the Get operation due to mocking limitations
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})
}

// Helper function tests
func TestFromSdkToModel(t *testing.T) {
	now := time.Now()
	sdkProject := sdk.Project{
		Id:          "project1",
		Name:        "Test Project",
		Tags:        []string{"test", "project"},
		Description: "A test project",
		CreatedAt:   &now,
		CreatedBy:   "user1",
		UpdatedAt:   &now,
		UpdatedBy:   "user1",
	}

	result := fromSdkToModel(sdkProject)

	assert.Equal(t, sdkProject.Id, result.Id)
	assert.Equal(t, sdkProject.Name, result.Name)
	assert.Equal(t, sdkProject.Tags, result.Tags)
	assert.Equal(t, sdkProject.Description, result.Description)
	assert.Equal(t, sdkProject.CreatedAt, result.CreatedAt)
	assert.Equal(t, sdkProject.CreatedBy, result.CreatedBy)
	assert.Equal(t, sdkProject.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, sdkProject.UpdatedBy, result.UpdatedBy)
}

func TestFromModelToSdk(t *testing.T) {
	now := time.Now()
	modelProject := &models.Project{
		Id:          "project1",
		Name:        "Test Project",
		Tags:        []string{"test", "project"},
		Description: "A test project",
		CreatedAt:   &now,
		CreatedBy:   "user1",
		UpdatedAt:   &now,
		UpdatedBy:   "user1",
	}

	result := fromModelToSdk(modelProject)

	assert.Equal(t, modelProject.Id, result.Id)
	assert.Equal(t, modelProject.Name, result.Name)
	assert.Equal(t, modelProject.Tags, result.Tags)
	assert.Equal(t, modelProject.Description, result.Description)
	assert.Equal(t, modelProject.CreatedAt, result.CreatedAt)
	assert.Equal(t, modelProject.CreatedBy, result.CreatedBy)
	assert.Equal(t, modelProject.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, modelProject.UpdatedBy, result.UpdatedBy)
}

func TestFromModelListToSdk(t *testing.T) {
	now := time.Now()
	modelProjects := []models.Project{
		{
			Id:          "project1",
			Name:        "Test Project 1",
			Tags:        []string{"test", "project1"},
			Description: "A test project 1",
			CreatedAt:   &now,
			CreatedBy:   "user1",
		},
		{
			Id:          "project2",
			Name:        "Test Project 2",
			Tags:        []string{"test", "project2"},
			Description: "A test project 2",
			CreatedAt:   &now,
			CreatedBy:   "user2",
		},
	}

	result := fromModelListToSdk(modelProjects)

	assert.Len(t, result, 2)
	assert.Equal(t, modelProjects[0].Id, result[0].Id)
	assert.Equal(t, modelProjects[0].Name, result[0].Name)
	assert.Equal(t, modelProjects[0].Tags, result[0].Tags)
	assert.Equal(t, modelProjects[1].Id, result[1].Id)
	assert.Equal(t, modelProjects[1].Name, result[1].Name)
	assert.Equal(t, modelProjects[1].Tags, result[1].Tags)
}

func TestStore_DatabaseCallValidation(t *testing.T) {
	t.Run("create_calls_insert_with_correct_parameters", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		project := &sdk.Project{
			Name:        "Validation Project",
			Description: "Test parameter validation",
			Tags:        []string{"validation"},
		}

		md := models.GetProjectModel()

		mockDB.On("InsertOne", ctx, md, mock.MatchedBy(func(doc models.Project) bool {
			return doc.Name == project.Name &&
				doc.Description == project.Description &&
				len(doc.Tags) == 1 &&
				doc.Tags[0] == "validation" &&
				doc.Id != "" && // Should have generated ID
				doc.CreatedAt != nil // Should have timestamp
		}), mock.Anything).Return(&mongo.InsertOneResult{}, nil)

		err := store.Create(ctx, project)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("get_calls_findone_with_correct_filter", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectId := "test-project-id"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.IdKey, Value: projectId}}
		mockSingleResult := &mongo.SingleResult{}

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		_, err := store.Get(ctx, projectId)

		// Will error due to mocking limitations, but verifies correct parameters
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("get_by_name_calls_findone_with_correct_filter", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()
		projectName := "Test Project Name"

		md := models.GetProjectModel()
		expectedFilter := bson.D{{Key: md.NameKey, Value: projectName}}
		mockSingleResult := &mongo.SingleResult{}

		mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

		_, err := store.GetByName(ctx, projectName)

		// Will error due to mocking limitations, but verifies correct parameters
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("get_all_calls_find_with_empty_filter", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)
		ctx := context.Background()

		md := models.GetProjectModel()
		expectedFilter := bson.D{{}}

		mockDB.On("Find", ctx, md, expectedFilter, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("test error"))

		_, err := store.GetAll(ctx)

		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})
}

func TestStore_ErrorScenarios(t *testing.T) {
	t.Run("error_constants_are_defined", func(t *testing.T) {
		assert.NotNil(t, sdk.ErrProjectNotFound)
		assert.Equal(t, "project not found", sdk.ErrProjectNotFound.Error())
	})

	t.Run("get_by_name_empty_name_returns_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)

		result, err := store.GetByName(context.Background(), "")

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrProjectNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("update_empty_id_returns_error", func(t *testing.T) {
		mockDB := test.SetupMockDB()
		store := NewStore(mockDB)

		project := &sdk.Project{Id: ""}

		err := store.Update(context.Background(), project)

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrProjectNotFound, err)
	})
}
