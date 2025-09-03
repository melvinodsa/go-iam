package test

import (
	"context"

	"github.com/melvinodsa/go-iam/db"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test helper functions
func SetupMockDB() *MockDB {
	return new(MockDB)
}

// MockCollection implements db.DbCollection interface for testing
type mockCollection struct {
	name   string
	dbName string
}

func NewMockCollection(name, dbName string) *mockCollection {
	return &mockCollection{
		name:   name,
		dbName: dbName,
	}
}

func (m mockCollection) Name() string {
	return m.name
}

func (m mockCollection) DbName() string {
	return m.dbName
}

// MockDB is a mock implementation of the DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) FindOne(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockDB) Find(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) InsertOne(ctx context.Context, col db.DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, col, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockDB) UpdateOne(ctx context.Context, col db.DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, col, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockDB) DeleteOne(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockDB) Aggregate(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, col, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockDB) CountDocuments(ctx context.Context, col db.DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, col, filter, opts)
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
