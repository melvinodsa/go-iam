package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// import fiber logger
	"github.com/gofiber/fiber/v2/log"
)

type DB interface {
	DbQuerier
	DbClient
}

type DbCollection interface {
	Name() string
	DbName() string
}

type DbQuerier interface {
	FindOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	InsertOne(ctx context.Context, col DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type DbClient interface {
	SetDbInContext(ctx context.Context) context.Context
	Disconnect(ctx context.Context) error
}

type MongoConnection struct {
	client *mongo.Client
}

type dbCtxKey struct{}

func (m MongoConnection) SetDbInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbCtxKey{}, &m)
}

func GetDbFromContext(ctx context.Context) DB {
	vl := ctx.Value(dbCtxKey{})
	db, ok := vl.(*MongoConnection)
	if !ok {
		log.Fatal("db not found in context")
	}
	return db
}

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

func (m *MongoConnection) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoConnection) FindOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.client.Database(col.DbName()).Collection(col.Name()).FindOne(ctx, filter, opts...)
}

func (m *MongoConnection) Find(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).Find(ctx, filter, opts...)
}

func (m *MongoConnection) Aggregate(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).Aggregate(ctx, filter, opts...)
}

func (m *MongoConnection) InsertOne(ctx context.Context, col DbCollection, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).InsertOne(ctx, document, opts...)
}

func (m *MongoConnection) UpdateOne(ctx context.Context, col DbCollection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).UpdateOne(ctx, filter, update, opts...)
}

func (m *MongoConnection) DeleteOne(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).DeleteOne(ctx, filter, opts...)
}

func (m *MongoConnection) CountDocuments(ctx context.Context, col DbCollection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return m.client.Database(col.DbName()).Collection(col.Name()).CountDocuments(ctx, filter, opts...)
}
