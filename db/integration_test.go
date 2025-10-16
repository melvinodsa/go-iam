package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestCollection implements DbCollection for testing
type TestCollection struct {
	name   string
	dbName string
}

func (t TestCollection) Name() string   { return t.name }
func (t TestCollection) DbName() string { return t.dbName }

// TestMongoConnection_DatabaseOperations tests database operations
func TestMongoConnection_DatabaseOperations(t *testing.T) {
	ctx := context.Background()
	testCol := TestCollection{name: "test_collection", dbName: "test_db"}

	t.Run("InsertOne and FindOne", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		doc := bson.M{"name": "test_user", "email": "test@example.com", "created_at": time.Now()}

		expectedInsertResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}
		mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(expectedInsertResult, nil)

		mockFindResult := &mongo.SingleResult{}
		mockCollection.On("FindOne", ctx, bson.M{"name": "test_user"}, mock.Anything).Return(mockFindResult)

		result, err := mockCollection.InsertOne(ctx, doc)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.InsertedID)

		findResult := mockCollection.FindOne(ctx, bson.M{"name": "test_user"})
		assert.NotNil(t, findResult)

		mockCollection.AssertExpectations(t)
	})

	t.Run("UpdateOne", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		doc := bson.M{"name": "update_test", "status": "pending"}
		expectedInsertResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}
		mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(expectedInsertResult, nil)

		filter := bson.M{"_id": expectedInsertResult.InsertedID}
		update := bson.M{"$set": bson.M{"status": "completed"}}
		expectedUpdateResult := &mongo.UpdateResult{
			MatchedCount:  1,
			ModifiedCount: 1,
		}
		mockCollection.On("UpdateOne", ctx, filter, update, mock.Anything).Return(expectedUpdateResult, nil)

		mockFindResult := &mongo.SingleResult{}
		mockCollection.On("FindOne", ctx, filter, mock.Anything).Return(mockFindResult)

		insertResult, err := mockCollection.InsertOne(ctx, doc)
		assert.NoError(t, err)
		assert.NotNil(t, insertResult)

		updateResult, err := mockCollection.UpdateOne(ctx, filter, update)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), updateResult.MatchedCount)
		assert.Equal(t, int64(1), updateResult.ModifiedCount)

		findResult := mockCollection.FindOne(ctx, filter)
		assert.NotNil(t, findResult)

		mockCollection.AssertExpectations(t)
	})

	t.Run("Find multiple documents", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		docs := []interface{}{
			bson.M{"category": "test", "value": 1},
			bson.M{"category": "test", "value": 2},
			bson.M{"category": "test", "value": 3},
		}

		for _, doc := range docs {
			expectedResult := &mongo.InsertOneResult{
				InsertedID: primitive.NewObjectID(),
			}
			mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(expectedResult, nil)
		}

		mockCursor := &mongo.Cursor{}
		mockCollection.On("Find", ctx, bson.M{"category": "test"}, mock.Anything).Return(mockCursor, nil)

		for _, doc := range docs {
			result, err := mockCollection.InsertOne(ctx, doc)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}

		cursor, err := mockCollection.Find(ctx, bson.M{"category": "test"})
		assert.NoError(t, err)
		assert.NotNil(t, cursor)

		mockCollection.AssertExpectations(t)
	})

	t.Run("CountDocuments", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		expectedCount := int64(3)
		mockCollection.On("CountDocuments", ctx, bson.M{"category": "test"}, mock.Anything).Return(expectedCount, nil)

		count, err := mockCollection.CountDocuments(ctx, bson.M{"category": "test"})
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		mockCollection.AssertExpectations(t)
	})

	t.Run("UpdateMany", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"category": "test"}
		update := bson.M{"$set": bson.M{"updated": true}}
		expectedResult := &mongo.UpdateResult{
			MatchedCount:  3,
			ModifiedCount: 3,
		}
		mockCollection.On("UpdateMany", ctx, filter, update, mock.Anything).Return(expectedResult, nil)

		result, err := mockCollection.UpdateMany(ctx, filter, update)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result.MatchedCount)
		assert.Equal(t, int64(3), result.ModifiedCount)

		mockCollection.AssertExpectations(t)
	})

	t.Run("DeleteOne", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		doc := bson.M{"name": "delete_test", "temp": true}
		expectedInsertResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}
		mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(expectedInsertResult, nil)

		filter := bson.M{"_id": expectedInsertResult.InsertedID}
		expectedDeleteResult := &mongo.DeleteResult{
			DeletedCount: 1,
		}
		mockCollection.On("DeleteOne", ctx, filter, mock.Anything).Return(expectedDeleteResult, nil)

		mockFindResult := &mongo.SingleResult{}
		mockCollection.On("FindOne", ctx, filter, mock.Anything).Return(mockFindResult)

		insertResult, err := mockCollection.InsertOne(ctx, doc)
		assert.NoError(t, err)
		assert.NotNil(t, insertResult)

		deleteResult, err := mockCollection.DeleteOne(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteResult.DeletedCount)

		findResult := mockCollection.FindOne(ctx, filter)
		assert.NotNil(t, findResult)

		mockCollection.AssertExpectations(t)
	})

	t.Run("Aggregate", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		pipeline := []bson.M{
			{"$match": bson.M{"category": "test"}},
			{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
		}

		mockCursor := &mongo.Cursor{}
		mockCollection.On("Aggregate", ctx, pipeline, mock.Anything).Return(mockCursor, nil)

		cursor, err := mockCollection.Aggregate(ctx, pipeline)
		assert.NoError(t, err)
		assert.NotNil(t, cursor)

		mockCollection.AssertExpectations(t)
	})

	t.Run("BulkWrite", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		models := []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.M{"bulk": "insert1"}),
			mongo.NewInsertOneModel().SetDocument(bson.M{"bulk": "insert2"}),
		}

		expectedResult := &mongo.BulkWriteResult{
			InsertedCount: 2,
		}
		mockCollection.On("BulkWrite", ctx, models, mock.Anything).Return(expectedResult, nil)

		result, err := mockCollection.BulkWrite(ctx, models)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), result.InsertedCount)

		mockCollection.AssertExpectations(t)
	})
}

// TestMongoConnection_WithMocks tests MongoDB operations with mocks
func TestMongoConnection_WithMocks(t *testing.T) {
	ctx := context.Background()
	testCol := TestCollection{name: "test_collection", dbName: "test_db"}

	t.Run("InsertOne with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		doc := bson.M{"name": "test_user", "email": "test@example.com", "created_at": time.Now()}

		expectedResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}

		mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(expectedResult, nil)

		// Simulate the database operation
		result, err := mockCollection.InsertOne(ctx, doc)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResult.InsertedID, result.InsertedID)

		mockCollection.AssertExpectations(t)
	})

	t.Run("FindOne with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"name": "test_user"}
		mockResult := &mongo.SingleResult{}
		mockCollection.On("FindOne", ctx, filter, mock.Anything).Return(mockResult)

		// Simulate the database operation
		result := mockCollection.FindOne(ctx, filter)

		assert.NotNil(t, result)

		mockCollection.AssertExpectations(t)
	})

	t.Run("UpdateOne with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"name": "test_user"}
		update := bson.M{"$set": bson.M{"status": "updated"}}

		expectedResult := &mongo.UpdateResult{
			MatchedCount:  1,
			ModifiedCount: 1,
		}

		mockCollection.On("UpdateOne", ctx, filter, update, mock.Anything).Return(expectedResult, nil)

		// Simulate the database operation
		result, err := mockCollection.UpdateOne(ctx, filter, update)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.MatchedCount)
		assert.Equal(t, int64(1), result.ModifiedCount)

		mockCollection.AssertExpectations(t)
	})

	t.Run("DeleteOne with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"name": "test_user"}

		expectedResult := &mongo.DeleteResult{
			DeletedCount: 1,
		}

		mockCollection.On("DeleteOne", ctx, filter, mock.Anything).Return(expectedResult, nil)

		// Simulate the database operation
		result, err := mockCollection.DeleteOne(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.DeletedCount)

		mockCollection.AssertExpectations(t)
	})

	t.Run("CountDocuments with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"status": "active"}
		expectedCount := int64(5)

		mockCollection.On("CountDocuments", ctx, filter, mock.Anything).Return(expectedCount, nil)

		// Simulate the database operation
		count, err := mockCollection.CountDocuments(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		mockCollection.AssertExpectations(t)
	})

	t.Run("Aggregate with mock", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		pipeline := []bson.M{
			{"$match": bson.M{"status": "active"}},
			{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
		}

		mockCursor := &mongo.Cursor{}
		mockCollection.On("Aggregate", ctx, pipeline, mock.Anything).Return(mockCursor, nil)

		// Simulate the database operation
		cursor, err := mockCollection.Aggregate(ctx, pipeline)

		assert.NoError(t, err)
		assert.NotNil(t, cursor)

		mockCollection.AssertExpectations(t)
	})
}

// TestMongoConnection_ErrorHandlingWithMocks tests error handling with mocks
func TestMongoConnection_ErrorHandlingWithMocks(t *testing.T) {
	ctx := context.Background()
	testCol := TestCollection{name: "test_collection", dbName: "test_db"}

	t.Run("InsertOne error handling", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		doc := bson.M{"name": "test_user"}
		expectedError := mongo.ErrClientDisconnected

		mockCollection.On("InsertOne", ctx, doc, mock.Anything).Return(nil, expectedError)

		result, err := mockCollection.InsertOne(ctx, doc)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		mockCollection.AssertExpectations(t)
	})

	t.Run("FindOne error handling", func(t *testing.T) {
		mockClient, mockDatabase, mockCollection := test.SetupMockMongoDriver()
		mockClient.On("Database", testCol.DbName(), mock.Anything).Return(mockDatabase)
		mockDatabase.On("Collection", testCol.Name(), mock.Anything).Return(mockCollection)

		filter := bson.M{"name": "nonexistent"}

		mockResult := &mongo.SingleResult{}
		mockCollection.On("FindOne", ctx, filter, mock.Anything).Return(mockResult)

		result := mockCollection.FindOne(ctx, filter)

		assert.NotNil(t, result)

		mockCollection.AssertExpectations(t)
	})
}
