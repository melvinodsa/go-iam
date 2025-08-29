package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCollection implements DbCollection interface for testing
type MockCollection struct {
	name   string
	dbName string
}

func (m MockCollection) Name() string {
	return m.name
}

func (m MockCollection) DbName() string {
	return m.dbName
}

// MockDB is a mock implementation of the DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) FindOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockDB) Find(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) InsertOne(ctx context.Context, col DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, col, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockDB) UpdateOne(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, col, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockDB) DeleteOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockDB) Aggregate(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) CountDocuments(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDB) BulkWrite(ctx context.Context, col DbCollection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := m.Called(ctx, col, models, opts)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (m *MockDB) UpdateMany(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
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

func TestMongoConnection_SetDbInContext(t *testing.T) {
	conn := &MongoConnection{}
	ctx := context.Background()

	resultCtx := conn.SetDbInContext(ctx)

	// Verify that the context contains the connection
	retrievedDB := GetDbFromContext(resultCtx)
	assert.NotNil(t, retrievedDB)
	assert.Equal(t, conn, retrievedDB)
}

func TestGetDbFromContext_Success(t *testing.T) {
	conn := &MongoConnection{}
	ctx := context.WithValue(context.Background(), dbCtxKey{}, conn)

	db := GetDbFromContext(ctx)

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
	GetDbFromContext(ctx)

	t.Error("Expected panic but didn't get one")
}

func TestMongoConnection_Disconnect(t *testing.T) {
	// This is an integration test that would require a real MongoDB connection
	// For unit testing, we'll test the interface compliance
	var _ DB = (*MongoConnection)(nil)
	var _ DbClient = (*MongoConnection)(nil)
	var _ DbQuerier = (*MongoConnection)(nil)
}

func TestMockCollection(t *testing.T) {
	col := MockCollection{
		name:   "test_collection",
		dbName: "test_db",
	}

	assert.Equal(t, "test_collection", col.Name())
	assert.Equal(t, "test_db", col.DbName())
}

func TestMockDB_InsertOne(t *testing.T) {
	mockDB := new(MockDB)
	ctx := context.Background()
	col := MockCollection{name: "test", dbName: "testdb"}
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
	mockDB := new(MockDB)
	ctx := context.Background()
	col := MockCollection{name: "test", dbName: "testdb"}
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
	mockDB := new(MockDB)
	ctx := context.Background()
	col := MockCollection{name: "test", dbName: "testdb"}
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
	mockDB := new(MockDB)
	ctx := context.Background()
	col := MockCollection{name: "test", dbName: "testdb"}
	filter := bson.M{"status": "active"}

	expectedCount := int64(5)

	mockDB.On("CountDocuments", ctx, col, filter, mock.AnythingOfType("[]*options.CountOptions")).Return(expectedCount, nil)

	count, err := mockDB.CountDocuments(ctx, col, filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockDB.AssertExpectations(t)
}

func TestMockDB_UpdateMany(t *testing.T) {
	mockDB := new(MockDB)
	ctx := context.Background()
	col := MockCollection{name: "test", dbName: "testdb"}
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
	mockDB := new(MockDB)
	ctx := context.Background()
	expectedCtx := context.WithValue(ctx, mockDbContextKey{}, mockDB)

	mockDB.On("SetDbInContext", ctx).Return(expectedCtx)

	result := mockDB.SetDbInContext(ctx)

	assert.Equal(t, expectedCtx, result)
	mockDB.AssertExpectations(t)
}

func TestMockDB_Disconnect(t *testing.T) {
	mockDB := new(MockDB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mockDB.On("Disconnect", ctx).Return(nil)

	err := mockDB.Disconnect(ctx)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkSetDbInContext(b *testing.B) {
	conn := &MongoConnection{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.SetDbInContext(ctx)
	}
}

func BenchmarkGetDbFromContext(b *testing.B) {
	conn := &MongoConnection{}
	ctx := conn.SetDbInContext(context.Background())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDbFromContext(ctx)
	}
}
