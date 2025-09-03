package resource

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestStore_Search_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	query := sdk.ResourceQuery{
		ProjectIds: []string{"project1"},
		Skip:       0,
		Limit:      10,
	}

	md := models.GetResourceModel()
	expectedCond := bson.D{{Key: md.EnabledKey, Value: true}, {Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: query.ProjectIds}}}}

	// Mock CountDocuments success
	mockDB.On("CountDocuments", ctx, md, expectedCond, mock.Anything).Return(int64(2), nil)

	cursor, _ := mongo.NewCursorFromDocuments([]interface{}{
		&models.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		},
	}, nil, nil)

	// Mock Find success
	mockDB.On("Find", ctx, md, expectedCond, mock.Anything).Return(cursor, nil)

	result, err := store.Search(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Resources, 1)
	mockDB.AssertExpectations(t)
}

func TestStore_Search_Success_WithFilters(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	query := sdk.ResourceQuery{
		ProjectIds:  []string{"project1"},
		Skip:        0,
		Limit:       10,
		Name:        "user",
		Description: "management",
		Key:         "users",
	}

	md := models.GetResourceModel()
	expectedCond := bson.D{
		{Key: md.EnabledKey, Value: true},
		{Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: query.ProjectIds}}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: md.NameKey, Value: primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", query.Name), Options: "i"}}},
			bson.D{{Key: md.DescriptionKey, Value: primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", query.Description), Options: "i"}}},
			bson.D{{Key: md.KeyKey, Value: primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", query.Key), Options: "i"}}},
		}},
	}

	// Mock CountDocuments success
	mockDB.On("CountDocuments", ctx, md, expectedCond, mock.Anything).Return(int64(2), nil)

	cursor, _ := mongo.NewCursorFromDocuments([]interface{}{
		&models.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		},
	}, nil, nil)

	// Mock Find success
	mockDB.On("Find", ctx, md, expectedCond, mock.Anything).Return(cursor, nil)

	result, err := store.Search(ctx, query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Resources, 1)
	mockDB.AssertExpectations(t)
}

func TestStore_Search_CountError(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	query := sdk.ResourceQuery{
		ProjectIds: []string{"project1"},
		Skip:       0,
		Limit:      10,
	}

	md := models.GetResourceModel()
	expectedCond := bson.D{{Key: md.EnabledKey, Value: true}, {Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: query.ProjectIds}}}}

	// Mock CountDocuments error
	mockDB.On("CountDocuments", ctx, md, expectedCond, mock.Anything).Return(int64(0), errors.New("count error"))

	result, err := store.Search(ctx, query)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error counting resources")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestStore_Search_FindError(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	query := sdk.ResourceQuery{
		ProjectIds: []string{"project1"},
		Skip:       0,
		Limit:      10,
	}

	md := models.GetResourceModel()
	expectedCond := bson.D{{Key: md.EnabledKey, Value: true}, {Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: query.ProjectIds}}}}

	// Mock CountDocuments success
	mockDB.On("CountDocuments", ctx, md, expectedCond, mock.Anything).Return(int64(2), nil)

	// Mock Find error
	mockDB.On("Find", ctx, md, expectedCond, mock.Anything).Return(nil, errors.New("find error"))

	result, err := store.Search(ctx, query)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error finding resources")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestStore_Get_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetResourceModel()

	record := &models.Resource{
		ID:          "resource1",
		Key:         "users",
		Name:        "Users Resource",
		Description: "Resource for user management",
		ProjectId:   "project1",
		Enabled:     true,
		CreatedBy:   "user1",
	}

	mockSingleResult := mongo.NewSingleResultFromDocument(record, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: "resource1"}, {Key: md.EnabledKey, Value: true}}, mock.Anything).Return(mockSingleResult)

	store := NewStore(mockDB)
	result, err := store.Get(ctx, "resource1")

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "resource1", result.ID)

}

func TestStore_Get_Error(t *testing.T) {
	mockDB := test.SetupMockDB()

	ctx := context.Background()
	md := models.GetResourceModel()

	mockSingleResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: "resource1"}, {Key: md.EnabledKey, Value: true}}, mock.Anything).Return(mockSingleResult, mongo.ErrNoDocuments)

	store := NewStore(mockDB)
	result, err := store.Get(ctx, "resource1")

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestStore_Get_FilterConstruction(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resourceID := "resource1"

	md := models.GetResourceModel()
	expectedFilter := bson.D{{Key: md.IdKey, Value: resourceID}, {Key: md.EnabledKey, Value: true}}

	// Create a mock SingleResult
	mockSingleResult := &mongo.SingleResult{}

	mockDB.On("FindOne", ctx, md, expectedFilter, mock.Anything).Return(mockSingleResult)

	// Since we can't easily mock SingleResult.Decode(), this test verifies the correct calls
	result, err := store.Get(ctx, resourceID)

	// This will fail due to mocking limitations, but verifies the correct calls
	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestStore_Create_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resource := &sdk.Resource{
		Name:        "Test Resource",
		Description: "Test description",
		Key:         "test-key",
		ProjectId:   "project1",
		Enabled:     true,
		CreatedBy:   "user1",
	}

	md := models.GetResourceModel()

	expectedResult := &mongo.InsertOneResult{
		InsertedID: "generated-id",
	}

	mockDB.On("InsertOne", ctx, md, mock.MatchedBy(func(doc interface{}) bool {
		modelDoc, ok := doc.(*models.Resource)
		if !ok {
			return false
		}
		return modelDoc.Name == resource.Name &&
			modelDoc.Description == resource.Description &&
			modelDoc.Key == resource.Key &&
			modelDoc.ProjectId == resource.ProjectId &&
			modelDoc.CreatedBy == resource.CreatedBy &&
			modelDoc.ID != "" && // Should have generated ID
			modelDoc.CreatedAt != nil // Should have timestamp
	}), mock.Anything).Return(expectedResult, nil)

	result, err := store.Create(ctx, resource)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, resource.ID)      // Should be set by Create method
	assert.NotNil(t, resource.CreatedAt) // Should be set by Create method
	mockDB.AssertExpectations(t)
}

func TestStore_Create_Error(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resource := &sdk.Resource{
		Name:        "Test Resource",
		Description: "Test description",
		Key:         "test-key",
		ProjectId:   "project1",
		Enabled:     true,
		CreatedBy:   "user1",
	}

	md := models.GetResourceModel()

	mockDB.On("InsertOne", ctx, md, mock.Anything, mock.Anything).Return((*mongo.InsertOneResult)(nil), errors.New("insert error"))

	result, err := store.Create(ctx, resource)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating resource")
	assert.Empty(t, result)
	mockDB.AssertExpectations(t)
}

func TestStore_Update_Success(t *testing.T) {
	mockDB := test.SetupMockDB()

	ctx := context.Background()
	resourceID := "resource1"

	now := time.Now()
	resource := &sdk.Resource{
		ID:        resourceID,
		Name:      "Updated Resource",
		CreatedAt: &now,
		CreatedBy: "original-user",
	}

	md := models.GetResourceModel()

	// Mock the Get call first (this will fail but we test the flow)
	expectedGetFilter := bson.D{{Key: md.IdKey, Value: resourceID}, {Key: md.EnabledKey, Value: true}}
	mockSingleResult := mongo.NewSingleResultFromDocument(resource, nil, nil)
	mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult, nil)
	mockDB.On("UpdateOne", ctx, md, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)

	store := NewStore(mockDB)
	err := store.Update(ctx, resource)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestStore_Update_Error(t *testing.T) {
	mockDB := test.SetupMockDB()

	ctx := context.Background()
	resourceID := "resource1"

	now := time.Now()
	resource := &sdk.Resource{
		ID:        resourceID,
		Name:      "Updated Resource",
		CreatedAt: &now,
		CreatedBy: "original-user",
	}

	md := models.GetResourceModel()

	// Mock the Get call first (this will fail but we test the flow)
	expectedGetFilter := bson.D{{Key: md.IdKey, Value: resourceID}, {Key: md.EnabledKey, Value: true}}
	mockSingleResult := mongo.NewSingleResultFromDocument(resource, nil, nil)
	mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult, nil)
	mockDB.On("UpdateOne", ctx, md, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, errors.New("error updating resource"))

	store := NewStore(mockDB)
	err := store.Update(ctx, resource)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestStore_Update_ResourceNotFound(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resource := &sdk.Resource{
		ID: "", // Empty ID should trigger error
	}

	err := store.Update(ctx, resource)

	assert.Error(t, err)
	assert.Equal(t, ErrResourceNotFound, err)
}

func TestStore_Update_UpdateCall(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resourceID := "resource1"

	now := time.Now()
	resource := &sdk.Resource{
		ID:        resourceID,
		Name:      "Updated Resource",
		CreatedAt: &now,
		CreatedBy: "original-user",
	}

	md := models.GetResourceModel()

	// Mock the Get call first (this will fail but we test the flow)
	expectedGetFilter := bson.D{{Key: md.IdKey, Value: resourceID}, {Key: md.EnabledKey, Value: true}}
	mockSingleResult := &mongo.SingleResult{}
	mockDB.On("FindOne", ctx, md, expectedGetFilter, mock.Anything).Return(mockSingleResult)

	err := store.Update(ctx, resource)

	// This will fail due to the Get() call mocking limitations
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestStore_Delete_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resourceID := "resource1"

	md := models.GetResourceModel()
	expectedFilter := bson.D{{Key: md.IdKey, Value: resourceID}}
	expectedUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: md.EnabledKey, Value: false}}}}

	expectedResult := &mongo.UpdateResult{
		ModifiedCount: 1,
	}

	mockDB.On("UpdateOne", ctx, md, expectedFilter, expectedUpdate, mock.Anything).Return(expectedResult, nil)

	err := store.Delete(ctx, resourceID)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestStore_Delete_Error(t *testing.T) {
	mockDB := test.SetupMockDB()
	store := NewStore(mockDB)

	ctx := context.Background()
	resourceID := "resource1"

	md := models.GetResourceModel()
	expectedFilter := bson.D{{Key: md.IdKey, Value: resourceID}}
	expectedUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: md.EnabledKey, Value: false}}}}

	mockDB.On("UpdateOne", ctx, md, expectedFilter, expectedUpdate, mock.Anything).Return((*mongo.UpdateResult)(nil), errors.New("delete error"))

	err := store.Delete(ctx, resourceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting resource")
	mockDB.AssertExpectations(t)
}

func TestNewStore(t *testing.T) {
	mockDB := test.SetupMockDB()

	store := NewStore(mockDB)

	assert.NotNil(t, store)
}

// Helper function tests
func TestFromModelToSdk(t *testing.T) {
	now := time.Now()
	model := &models.Resource{
		ID:          "resource1",
		Name:        "Test Resource",
		Description: "Test description",
		Key:         "test-key",
		ProjectId:   "project1",
		Enabled:     true,
		CreatedAt:   &now,
		CreatedBy:   "user1",
		UpdatedAt:   &now,
		UpdatedBy:   "user1",
	}

	result := fromModelToSdk(model)

	assert.Equal(t, model.ID, result.ID)
	assert.Equal(t, model.Name, result.Name)
	assert.Equal(t, model.Description, result.Description)
	assert.Equal(t, model.Key, result.Key)
	assert.Equal(t, model.ProjectId, result.ProjectId)
	assert.Equal(t, model.Enabled, result.Enabled)
	assert.Equal(t, model.CreatedAt, result.CreatedAt)
	assert.Equal(t, model.CreatedBy, result.CreatedBy)
	assert.Equal(t, model.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, model.UpdatedBy, result.UpdatedBy)
}

func TestFromSdkToModel(t *testing.T) {
	now := time.Now()
	sdkResource := sdk.Resource{
		ID:          "resource1",
		Name:        "Test Resource",
		Description: "Test description",
		Key:         "test-key",
		ProjectId:   "project1",
		Enabled:     true,
		CreatedAt:   &now,
		CreatedBy:   "user1",
		UpdatedAt:   &now,
		UpdatedBy:   "user1",
	}

	result := fromSdkToModel(sdkResource)

	assert.Equal(t, sdkResource.ID, result.ID)
	assert.Equal(t, sdkResource.Name, result.Name)
	assert.Equal(t, sdkResource.Description, result.Description)
	assert.Equal(t, sdkResource.Key, result.Key)
	assert.Equal(t, sdkResource.ProjectId, result.ProjectId)
	assert.Equal(t, sdkResource.Enabled, result.Enabled)
	assert.Equal(t, sdkResource.CreatedAt, result.CreatedAt)
	assert.Equal(t, sdkResource.CreatedBy, result.CreatedBy)
	assert.Equal(t, sdkResource.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, sdkResource.UpdatedBy, result.UpdatedBy)
}

func TestFromModelListToSdk(t *testing.T) {
	now := time.Now()
	models := []models.Resource{
		{
			ID:          "resource1",
			Name:        "Test Resource 1",
			Description: "Test description 1",
			Key:         "test-key-1",
			ProjectId:   "project1",
			Enabled:     true,
			CreatedAt:   &now,
			CreatedBy:   "user1",
		},
		{
			ID:          "resource2",
			Name:        "Test Resource 2",
			Description: "Test description 2",
			Key:         "test-key-2",
			ProjectId:   "project2",
			Enabled:     true,
			CreatedAt:   &now,
			CreatedBy:   "user2",
		},
	}

	result := fromModelListToSdk(models)

	assert.Len(t, result, 2)
	assert.Equal(t, models[0].ID, result[0].ID)
	assert.Equal(t, models[0].Name, result[0].Name)
	assert.Equal(t, models[1].ID, result[1].ID)
	assert.Equal(t, models[1].Name, result[1].Name)
}

func TestStore_Search_FilterConstruction(t *testing.T) {
	tests := []struct {
		name          string
		query         sdk.ResourceQuery
		expectedCalls func(*test.MockDB, models.ResourceModel)
	}{
		{
			name: "query with name filter",
			query: sdk.ResourceQuery{
				Name:       "test",
				ProjectIds: []string{"project1"},
				Skip:       0,
				Limit:      10,
			},
			expectedCalls: func(mockDB *test.MockDB, md models.ResourceModel) {
				// When name is provided, it should use $or filter
				mockDB.On("CountDocuments", mock.Anything, md, mock.MatchedBy(func(filter interface{}) bool {
					filterDoc, ok := filter.(bson.D)
					if !ok {
						return false
					}
					// Should contain $or key when name filter is provided
					for _, elem := range filterDoc {
						if elem.Key == "$or" {
							return true
						}
					}
					return false
				}), mock.Anything).Return(int64(0), errors.New("test error"))
			},
		},
		{
			name: "query without filters",
			query: sdk.ResourceQuery{
				ProjectIds: []string{"project1"},
				Skip:       0,
				Limit:      10,
			},
			expectedCalls: func(mockDB *test.MockDB, md models.ResourceModel) {
				// When no search filters, should use enabled and project filters only
				expectedCond := bson.D{{Key: md.EnabledKey, Value: true}, {Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: []string{"project1"}}}}}
				mockDB.On("CountDocuments", mock.Anything, md, expectedCond, mock.Anything).Return(int64(0), errors.New("test error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := test.SetupMockDB()
			store := NewStore(mockDB)

			md := models.GetResourceModel()
			tt.expectedCalls(mockDB, md)

			_, err := store.Search(context.Background(), tt.query)

			assert.Error(t, err) // Expected since we're mocking an error
			mockDB.AssertExpectations(t)
		})
	}
}
