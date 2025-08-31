package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongoConnection_SetDbInContext(t *testing.T) {
	conn := &db.MongoConnection{}
	ctx := context.Background()

	resultCtx := conn.SetDbInContext(ctx)

	// Verify that the context contains the connection
	retrievedDB := db.GetDbFromContext(resultCtx)
	assert.NotNil(t, retrievedDB)
	assert.Equal(t, conn, retrievedDB)
}

func TestGetDbFromContext_Success(t *testing.T) {
	conn := &db.MongoConnection{}
	ctx := context.WithValue(context.Background(), db.DbCtxKey{}, conn)

	db := db.GetDbFromContext(ctx)

	assert.NotNil(t, db)
	assert.Equal(t, conn, db)
}

func TestGetDbFromContext_Panic(t *testing.T) {
	// This test checks that GetDbFromContext panics when no DB is in context
	// We need to capture the panic since log.Fatal will exit the program
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to log.Fatal
			assert.Contains(t, r.(string), "db not found in context")
		}
	}()

	ctx := context.Background()
	db.GetDbFromContext(ctx)

	t.Error("Expected panic but didn't get one")
}

func TestMongoConnection_Disconnect(t *testing.T) {
	// This is an integration test that would require a real MongoDB connection
	// For unit testing, we'll test the interface compliance
	var _ db.DB = (*db.MongoConnection)(nil)
	var _ db.DbClient = (*db.MongoConnection)(nil)
	var _ db.DbQuerier = (*db.MongoConnection)(nil)
}

func TestMockCollection(t *testing.T) {
	col := test.NewMockCollection("test_collection", "test_db")

	assert.Equal(t, "test_collection", col.Name())
	assert.Equal(t, "test_db", col.DbName())
}

func TestMockDB_InsertOne(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	col := test.NewMockCollection("test", "testdb")
	doc := bson.M{"name": "test"}

	expectedResult := &mongo.InsertOneResult{
		InsertedID: "test-id",
	}

	mockDB.On("InsertOne", ctx, col, doc, mock.AnythingOfType("[]*options.InsertOneOptions")).Return(expectedResult, nil)

	result, err := mockDB.InsertOne(ctx, col, doc)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestMockDB_UpdateOne(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	col := test.NewMockCollection("test", "testdb")
	filter := bson.M{"_id": "test-id"}
	update := bson.M{"$set": bson.M{"name": "updated"}}

	expectedResult := &mongo.UpdateResult{
		MatchedCount:  1,
		ModifiedCount: 1,
	}

	mockDB.On("UpdateOne", ctx, col, filter, update, mock.AnythingOfType("[]*options.UpdateOptions")).Return(expectedResult, nil)

	result, err := mockDB.UpdateOne(ctx, col, filter, update)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestMockDB_DeleteOne(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	col := test.NewMockCollection("test", "testdb")
	filter := bson.M{"_id": "test-id"}

	expectedResult := &mongo.DeleteResult{
		DeletedCount: 1,
	}

	mockDB.On("DeleteOne", ctx, col, filter, mock.AnythingOfType("[]*options.DeleteOptions")).Return(expectedResult, nil)

	result, err := mockDB.DeleteOne(ctx, col, filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestMockDB_CountDocuments(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	col := test.NewMockCollection("test", "testdb")
	filter := bson.M{"status": "active"}

	expectedCount := int64(5)

	mockDB.On("CountDocuments", ctx, col, filter, mock.AnythingOfType("[]*options.CountOptions")).Return(expectedCount, nil)

	count, err := mockDB.CountDocuments(ctx, col, filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockDB.AssertExpectations(t)
}

func TestMockDB_UpdateMany(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	col := test.NewMockCollection("test", "testdb")
	filter := bson.M{"status": "pending"}
	update := bson.M{"$set": bson.M{"status": "processed"}}

	expectedResult := &mongo.UpdateResult{
		MatchedCount:  3,
		ModifiedCount: 3,
	}

	mockDB.On("UpdateMany", ctx, col, filter, update, mock.AnythingOfType("[]*options.UpdateOptions")).Return(expectedResult, nil)

	result, err := mockDB.UpdateMany(ctx, col, filter, update)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

type mockDbContextKey struct{}

func TestMockDB_SetDbInContext(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx := context.Background()
	expectedCtx := context.WithValue(ctx, mockDbContextKey{}, mockDB)

	mockDB.On("SetDbInContext", ctx).Return(expectedCtx)

	result := mockDB.SetDbInContext(ctx)

	assert.Equal(t, expectedCtx, result)
	mockDB.AssertExpectations(t)
}

func TestMockDB_Disconnect(t *testing.T) {
	mockDB := new(test.MockDB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockDB.On("Disconnect", ctx).Return(nil)

	err := mockDB.Disconnect(ctx)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkSetDbInContext(b *testing.B) {
	conn := &db.MongoConnection{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.SetDbInContext(ctx)
	}
}

func BenchmarkGetDbFromContext(b *testing.B) {
	conn := &db.MongoConnection{}
	ctx := conn.SetDbInContext(context.Background())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.GetDbFromContext(ctx)
	}
}
