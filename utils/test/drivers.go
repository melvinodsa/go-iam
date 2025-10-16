package test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupMockMongoDriver creates a complete mock MongoDB driver setup
func SetupMockMongoDriver() (*MockMongoClient, *MockMongoDatabase, *MockMongoCollection) {
	mockClient := &MockMongoClient{}
	mockDatabase := &MockMongoDatabase{}
	mockCollection := &MockMongoCollection{}

	// Set up the chain: Client -> Database -> Collection
	mockClient.On("Database", mock.AnythingOfType("string"), mock.Anything).Return(mockDatabase)
	mockDatabase.On("Collection", mock.AnythingOfType("string"), mock.Anything).Return(mockCollection)

	return mockClient, mockDatabase, mockCollection
}

// SetupMockRedisDriver creates a complete mock Redis driver setup
func SetupMockRedisDriver() *MockRedisClient {
	return &MockRedisClient{}
}

// SetupMockRedisWithCommands creates a Redis mock with common command mocks
func SetupMockRedisWithCommands() *MockRedisClient {
	mockClient := &MockRedisClient{}

	// Mock common Redis commands
	mockClient.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(&MockRedisStatusCmd{})
	mockClient.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(&MockRedisStringCmd{})
	mockClient.On("Del", mock.Anything, mock.AnythingOfType("[]string")).Return(&MockRedisIntCmd{})
	mockClient.On("Exists", mock.Anything, mock.AnythingOfType("[]string")).Return(&MockRedisIntCmd{})
	mockClient.On("Expire", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(&MockRedisBoolCmd{})
	mockClient.On("Keys", mock.Anything, mock.AnythingOfType("string")).Return(&MockRedisStringSliceCmd{})
	mockClient.On("Ping", mock.Anything).Return(&MockRedisStatusCmd{})
	mockClient.On("Close").Return(nil)

	return mockClient
}

// SetupMockMongoWithOperations creates a MongoDB mock with common operation mocks
func SetupMockMongoWithOperations() (*MockMongoClient, *MockMongoDatabase, *MockMongoCollection) {
	mockClient, mockDatabase, mockCollection := SetupMockMongoDriver()

	// Mock common MongoDB operations
	mockCollection.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil)
	mockCollection.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(&MockMongoSingleResult{})
	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(&MockMongoCursor{}, nil)
	mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)
	mockCollection.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	mockCollection.On("CountDocuments", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
	mockCollection.On("Aggregate", mock.Anything, mock.Anything, mock.Anything).Return(&MockMongoCursor{}, nil)

	return mockClient, mockDatabase, mockCollection
}

// MockMongoSingleResultWithData creates a SingleResult mock that returns specific data
func MockMongoSingleResultWithData(data interface{}) *MockMongoSingleResult {
	mockResult := &MockMongoSingleResult{}
	mockResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		// This would need to be implemented based on the actual data structure
		// For now, we just return no error
	}).Return(nil)
	mockResult.On("Err").Return(nil)
	return mockResult
}

// MockMongoSingleResultWithError creates a SingleResult mock that returns an error
func MockMongoSingleResultWithError(err error) *MockMongoSingleResult {
	mockResult := &MockMongoSingleResult{}
	mockResult.On("Decode", mock.Anything).Return(err)
	mockResult.On("Err").Return(err)
	return mockResult
}

// MockMongoCursorWithData creates a Cursor mock that returns specific data
func MockMongoCursorWithData(data []interface{}) *MockMongoCursor {
	mockCursor := &MockMongoCursor{}

	// Mock cursor iteration
	callCount := 0
	mockCursor.On("Next", mock.Anything).Run(func(args mock.Arguments) {
		callCount++
	}).Return(func(ctx context.Context) bool {
		return callCount <= len(data)
	})

	mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		// This would need to be implemented based on the actual data structure
		// For now, we just return no error
	}).Return(nil)

	mockCursor.On("Err").Return(nil)
	mockCursor.On("Close", mock.Anything).Return(nil)

	return mockCursor
}

// MockRedisStatusCmdWithValue creates a StatusCmd mock that returns a specific value
func MockRedisStatusCmdWithValue(value string) *MockRedisStatusCmd {
	mockCmd := &MockRedisStatusCmd{}
	mockCmd.On("Val").Return(value)
	mockCmd.On("Result").Return(value, nil)
	mockCmd.On("Err").Return(nil)
	return mockCmd
}

// MockRedisStringCmdWithValue creates a StringCmd mock that returns a specific value
func MockRedisStringCmdWithValue(value string) *MockRedisStringCmd {
	mockCmd := &MockRedisStringCmd{}
	mockCmd.On("Val").Return(value)
	mockCmd.On("Result").Return(value, nil)
	mockCmd.On("Err").Return(nil)
	return mockCmd
}

// MockRedisIntCmdWithValue creates an IntCmd mock that returns a specific value
func MockRedisIntCmdWithValue(value int64) *MockRedisIntCmd {
	mockCmd := &MockRedisIntCmd{}
	mockCmd.On("Val").Return(value)
	mockCmd.On("Result").Return(value, nil)
	mockCmd.On("Err").Return(nil)
	return mockCmd
}

// MockRedisBoolCmdWithValue creates a BoolCmd mock that returns a specific value
func MockRedisBoolCmdWithValue(value bool) *MockRedisBoolCmd {
	mockCmd := &MockRedisBoolCmd{}
	mockCmd.On("Val").Return(value)
	mockCmd.On("Result").Return(value, nil)
	mockCmd.On("Err").Return(nil)
	return mockCmd
}

// MockRedisStringSliceCmdWithValue creates a StringSliceCmd mock that returns specific values
func MockRedisStringSliceCmdWithValue(values []string) *MockRedisStringSliceCmd {
	mockCmd := &MockRedisStringSliceCmd{}
	mockCmd.On("Val").Return(values)
	mockCmd.On("Result").Return(values, nil)
	mockCmd.On("Err").Return(nil)
	return mockCmd
}

// MockRedisCmdWithError creates a Cmd mock that returns an error
func MockRedisCmdWithError(err error) *MockRedisCmd {
	mockCmd := &MockRedisCmd{}
	mockCmd.On("Err").Return(err)
	return mockCmd
}

// MockMongoResultWithError creates a MongoDB result mock that returns an error
func MockMongoResultWithError(err error) *MockMongoSingleResult {
	mockResult := &MockMongoSingleResult{}
	mockResult.On("Decode", mock.Anything).Return(err)
	mockResult.On("Err").Return(err)
	return mockResult
}

// MockMongoCursorWithError creates a MongoDB cursor mock that returns an error
func MockMongoCursorWithError(err error) *MockMongoCursor {
	mockCursor := &MockMongoCursor{}
	mockCursor.On("Next", mock.Anything).Return(false)
	mockCursor.On("Err").Return(err)
	mockCursor.On("Close", mock.Anything).Return(err)
	return mockCursor
}

// Helper function to create a mock context
func MockContext() context.Context {
	return context.Background()
}

// Helper function to create a mock context with timeout
func MockContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// Helper function to create a mock collection name
func MockCollectionName() string {
	return "test_collection"
}

// Helper function to create a mock database name
func MockDatabaseName() string {
	return "test_database"
}

// Helper function to create a mock Redis key
func MockRedisKey() string {
	return "test_key"
}

// Helper function to create a mock Redis value
func MockRedisValue() string {
	return "test_value"
}

// Helper function to create a mock TTL duration
func MockTTL() time.Duration {
	return 5 * time.Minute
}

// Helper function to create a mock BSON document
func MockBSONDocument() bson.M {
	return bson.M{
		"id":    "test_id",
		"name":  "test_name",
		"value": "test_value",
	}
}

// Helper function to create a mock BSON filter
func MockBSONFilter() bson.M {
	return bson.M{
		"id": "test_id",
	}
}

// Helper function to create a mock BSON update
func MockBSONUpdate() bson.M {
	return bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
}
