package test

import (
	"github.com/stretchr/testify/mock"
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
