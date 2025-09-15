// Package db provides database abstraction layer for MongoDB operations.
// It defines interfaces and implementations for database connectivity,
// querying, and transaction management used throughout the Go IAM application.
package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// import fiber logger
	"github.com/gofiber/fiber/v2/log"
)

// DB is the main database interface that combines querying and client operations.
// It provides a unified interface for all database operations in the application.
type DB interface {
	DbQuerier
	DbClient
}

// DbCollection represents a database collection with its metadata.
// Implementations should provide the collection name and database name.
type DbCollection interface {
	Name() string   // Returns the collection name
	DbName() string // Returns the database name
}

// DbQuerier defines the interface for database query operations.
// All MongoDB query operations are abstracted through this interface.
type DbQuerier interface {
	// FindOne executes a find operation that returns a single document
	FindOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult

	// Find executes a find operation that can return multiple documents
	Find(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)

	// InsertOne inserts a single document into the collection
	InsertOne(ctx context.Context, col DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// UpdateOne updates a single document in the collection
	UpdateOne(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	// DeleteOne deletes a single document from the collection
	DeleteOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// Aggregate executes an aggregation pipeline
	Aggregate(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)

	// CountDocuments counts the number of documents matching the filter
	CountDocuments(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error)

	// BulkWrite executes multiple write operations in bulk
	BulkWrite(ctx context.Context, col DbCollection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)

	// UpdateMany updates multiple documents in the collection
	UpdateMany(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// DbClient defines the interface for database client operations.
// It handles database connection lifecycle and context management.
// DbClient defines the interface for database client operations.
// It handles database connection lifecycle and context management.
type DbClient interface {
	// SetDbInContext stores the database connection in the given context
	SetDbInContext(ctx context.Context) context.Context

	// Disconnect closes the database connection
	Disconnect(ctx context.Context) error
}

// MongoConnection implements the DB interface for MongoDB.
// It wraps a MongoDB client and provides all database operations.
type MongoConnection struct {
	client *mongo.Client // MongoDB client instance
}

// DbCtxKey is the context key type used for storing database connections in context.
type DbCtxKey struct{}

// SetDbInContext stores this MongoConnection instance in the provided context.
// This allows handlers and services to access the database connection.
//
// Parameters:
//   - ctx: The context to store the database connection in
//
// Returns a new context with the database connection stored.
func (m MongoConnection) SetDbInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, DbCtxKey{}, &m)
}

// GetDbFromContext retrieves the database connection from the provided context.
// This function is used by handlers and services to access the database.
//
// Parameters:
//   - ctx: The context containing the database connection
//
// Returns the DB interface implementation stored in the context.
// Panics if no database connection is found in the context.
func GetDbFromContext(ctx context.Context) DB {
	vl := ctx.Value(DbCtxKey{})
	db, ok := vl.(*MongoConnection)
	if !ok {
		panic("db not found in context")
	}
	return db
}

// NewMongoConnection creates a new MongoDB connection using the provided connection URL.
// It establishes a connection to MongoDB and returns a MongoConnection instance.
//
// Parameters:
//   - url: MongoDB connection string (e.g., "mongodb://localhost:27017")
//
// Returns:
//   - *MongoConnection: A new MongoDB connection instance
//   - error: An error if the connection cannot be established
//
// The connection uses MongoDB Server API version 1 for compatibility.
func NewMongoConnection(url string) (*MongoConnection, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	log.Info("Connected to MongoDB")
	return &MongoConnection{client: client}, nil
}

// Disconnect closes the MongoDB connection.
// This should be called when shutting down the application to properly clean up resources.
//
// Parameters:
//   - ctx: Context for the disconnect operation
//
// Returns an error if the disconnect operation fails.
func (m *MongoConnection) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// FindOne executes a find operation that returns at most one document.
// This is a wrapper around MongoDB's FindOne operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to query
//   - filter: Query filter
//   - opts: Optional find options
//
// Returns a SingleResult that can be decoded into a document.
func (m *MongoConnection) FindOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.client.Database(col.DbName()).Collection(col.Name()).FindOne(ctx, filter, opts...)
}

// Find executes a find operation that can return multiple documents.
// This is a wrapper around MongoDB's Find operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to query
//   - filter: Query filter
//   - opts: Optional find options
//
// Returns a Cursor for iterating over results and an error if the operation fails.
func (m *MongoConnection) Find(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).Find(ctx, filter, opts...)
}

// Aggregate executes an aggregation pipeline operation.
// This is a wrapper around MongoDB's Aggregate operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to aggregate
//   - filter: Aggregation pipeline stages
//   - opts: Optional aggregation options
//
// Returns a Cursor for iterating over aggregation results and an error if the operation fails.
func (m *MongoConnection) Aggregate(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).Aggregate(ctx, filter, opts...)
}

// InsertOne inserts a single document into the collection.
// This is a wrapper around MongoDB's InsertOne operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to insert into
//   - document: Document to insert
//   - opts: Optional insert options
//
// Returns InsertOneResult containing the inserted document's ID and an error if the operation fails.
func (m *MongoConnection) InsertOne(ctx context.Context, col DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).InsertOne(ctx, document, opts...)
}

// UpdateOne updates a single document in the collection.
// This is a wrapper around MongoDB's UpdateOne operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to update
//   - filter: Query filter to match documents
//   - update: Update operations to apply
//   - opts: Optional update options
//
// Returns UpdateResult containing information about the update operation and an error if it fails.
func (m *MongoConnection) UpdateOne(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).UpdateOne(ctx, filter, update, opts...)
}

// DeleteOne deletes a single document from the collection.
// This is a wrapper around MongoDB's DeleteOne operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to delete from
//   - filter: Query filter to match the document to delete
//   - opts: Optional delete options
//
// Returns DeleteResult containing information about the delete operation and an error if it fails.
func (m *MongoConnection) DeleteOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).DeleteOne(ctx, filter, opts...)
}

// CountDocuments counts the number of documents in the collection that match the filter.
// This is a wrapper around MongoDB's CountDocuments operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to count documents in
//   - filter: Query filter to match documents
//   - opts: Optional count options
//
// Returns the number of matching documents and an error if the operation fails.
func (m *MongoConnection) CountDocuments(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).CountDocuments(ctx, filter, opts...)
}

// BulkWrite executes multiple write operations in a single bulk operation.
// This is a wrapper around MongoDB's BulkWrite operation for improved performance.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to perform bulk operations on
//   - models: Slice of write models (InsertOne, UpdateOne, DeleteOne, etc.)
//   - opts: Optional bulk write options
//
// Returns BulkWriteResult containing information about the bulk operation and an error if it fails.
func (m *MongoConnection) BulkWrite(ctx context.Context, col DbCollection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).BulkWrite(ctx, models, opts...)
}

// UpdateMany updates multiple documents in the collection that match the filter.
// This is a wrapper around MongoDB's UpdateMany operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - col: Collection to update
//   - filter: Query filter to match documents
//   - update: Update operations to apply to all matching documents
//   - opts: Optional update options
//
// Returns UpdateResult containing information about the update operation and an error if it fails.
func (m *MongoConnection) UpdateMany(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).UpdateMany(ctx, filter, update, opts...)
}
