package test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// MockMongoClient is a mock implementation of mongo.Client
type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	args := m.Called(name, opts)
	return args.Get(0).(*mongo.Database)
}

func (m *MockMongoClient) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(mongo.ListDatabasesResult), args.Error(1)
}

func (m *MockMongoClient) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	args := m.Called(ctx, rp)
	return args.Error(0)
}

func (m *MockMongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	args := m.Called(opts)
	return args.Get(0).(mongo.Session), args.Error(1)
}

func (m *MockMongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockMongoClient) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	args := m.Called(ctx, opts, fn)
	return args.Error(0)
}

func (m *MockMongoClient) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := m.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}

func (m *MockMongoClient) NumberSessionsInProgress() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMongoClient) EndSessions(ctx context.Context, sessions ...mongo.Session) error {
	args := m.Called(ctx, sessions)
	return args.Error(0)
}

func (m *MockMongoClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockMongoDatabase is a mock implementation of mongo.Database
type MockMongoDatabase struct {
	mock.Mock
}

func (m *MockMongoDatabase) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMongoDatabase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	args := m.Called(name, opts)
	return args.Get(0).(*mongo.Collection)
}

func (m *MockMongoDatabase) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoDatabase) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	args := m.Called(ctx, runCommand, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoDatabase) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, runCommand, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoDatabase) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoDatabase) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMongoDatabase) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoDatabase) ReadConcern() *readconcern.ReadConcern {
	args := m.Called()
	return args.Get(0).(*readconcern.ReadConcern)
}

func (m *MockMongoDatabase) ReadPreference() *readpref.ReadPref {
	args := m.Called()
	return args.Get(0).(*readpref.ReadPref)
}

func (m *MockMongoDatabase) WriteConcern() *writeconcern.WriteConcern {
	args := m.Called()
	return args.Get(0).(*writeconcern.WriteConcern)
}

func (m *MockMongoDatabase) Client() *mongo.Client {
	args := m.Called()
	return args.Get(0).(*mongo.Client)
}

// MockMongoCollection is a mock implementation of mongo.Collection
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMongoCollection) Database() *mongo.Database {
	args := m.Called()
	return args.Get(0).(*mongo.Database)
}

func (m *MockMongoCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := m.Called(ctx, models, opts)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	args := m.Called(ctx, documents, opts)
	return args.Get(0).(*mongo.InsertManyResult), args.Error(1)
}

func (m *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockMongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, replacement, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, replacement, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMongoCollection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMongoCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := m.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockMongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoCollection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := m.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}

func (m *MockMongoCollection) Indexes() mongo.IndexView {
	args := m.Called()
	return args.Get(0).(mongo.IndexView)
}

func (m *MockMongoCollection) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoCollection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	args := m.Called(opts)
	return args.Get(0).(*mongo.Collection), args.Error(1)
}

func (m *MockMongoCollection) ReadConcern() *readconcern.ReadConcern {
	args := m.Called()
	return args.Get(0).(*readconcern.ReadConcern)
}

func (m *MockMongoCollection) ReadPreference() *readpref.ReadPref {
	args := m.Called()
	return args.Get(0).(*readpref.ReadPref)
}

func (m *MockMongoCollection) WriteConcern() *writeconcern.WriteConcern {
	args := m.Called()
	return args.Get(0).(*writeconcern.WriteConcern)
}

// MockMongoSingleResult is a mock implementation of mongo.SingleResult
type MockMongoSingleResult struct {
	mock.Mock
}

func (m *MockMongoSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockMongoSingleResult) DecodeBytes() (bson.Raw, error) {
	args := m.Called()
	return args.Get(0).(bson.Raw), args.Error(1)
}

func (m *MockMongoSingleResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockMongoCursor is a mock implementation of mongo.Cursor
type MockMongoCursor struct {
	mock.Mock
}

func (m *MockMongoCursor) ID() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockMongoCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockMongoCursor) TryNext(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockMongoCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockMongoCursor) DecodeBytes() (bson.Raw, error) {
	args := m.Called()
	return args.Get(0).(bson.Raw), args.Error(1)
}

func (m *MockMongoCursor) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMongoCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

func (m *MockMongoCursor) RemainingBatchLength() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMongoCursor) BatchSize() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockMongoCursor) SetBatchSize(size int32) *mongo.Cursor {
	args := m.Called(size)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMaxTime(d time.Duration) *mongo.Cursor {
	args := m.Called(d)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetComment(comment interface{}) *mongo.Cursor {
	args := m.Called(comment)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetHint(hint interface{}) *mongo.Cursor {
	args := m.Called(hint)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetLimit(limit int64) *mongo.Cursor {
	args := m.Called(limit)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMaxAwaitTime(d time.Duration) *mongo.Cursor {
	args := m.Called(d)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMin(bson.Raw) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMax(bson.Raw) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetReturnKey(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetShowRecordID(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetSnapshot(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetSort(interface{}) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetNoCursorTimeout(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetOplogReplay(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetPartial(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetTailable(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetAllowDiskUse(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetAllowPartialResults(bool) *mongo.Cursor {
	args := m.Called()
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetBatchSizeFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetCursorType(ct int) *mongo.Cursor {
	args := m.Called(ct)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMaxTimeFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMinFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetMaxFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetReturnKeyFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetShowRecordIDFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetSnapshotFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetSortFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetNoCursorTimeoutFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetOplogReplayFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetPartialFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetTailableFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetAllowDiskUseFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}

func (m *MockMongoCursor) SetAllowPartialResultsFromCursor(cursor *mongo.Cursor) *mongo.Cursor {
	args := m.Called(cursor)
	return args.Get(0).(*mongo.Cursor)
}
