package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestCollection implements DbCollection for testing
type TestCollection struct {
	name   string
	dbName string
}

func (t TestCollection) Name() string   { return t.name }
func (t TestCollection) DbName() string { return t.dbName }

func TestNewMongoConnection(t *testing.T) {
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}

	t.Run("successful connection", func(t *testing.T) {
		conn, err := db.NewMongoConnection(mongoURL)
		assert.NoError(t, err)
		assert.NotNil(t, conn)

		// Test disconnect
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = conn.Disconnect(ctx)
		assert.NoError(t, err)
	})

	t.Run("invalid connection string", func(t *testing.T) {
		conn, err := db.NewMongoConnection("invalid://connection")
		assert.Error(t, err)
		assert.Nil(t, conn)
	})
}

func TestMongoConnection_DatabaseOperations(t *testing.T) {
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}

	conn, err := db.NewMongoConnection(mongoURL)
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn.Disconnect(ctx)
	}()

	ctx := context.Background()
	testCol := TestCollection{name: "test_collection", dbName: "test_db"}

	t.Run("InsertOne and FindOne", func(t *testing.T) {
		doc := bson.M{"name": "test_user", "email": "test@example.com", "created_at": time.Now()}

		// Insert document
		result, err := conn.InsertOne(ctx, testCol, doc)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.InsertedID)

		// Find the inserted document
		var foundDoc bson.M
		err = conn.FindOne(ctx, testCol, bson.M{"name": "test_user"}).Decode(&foundDoc)
		assert.NoError(t, err)
		assert.Equal(t, "test_user", foundDoc["name"])
		assert.Equal(t, "test@example.com", foundDoc["email"])
	})

	t.Run("UpdateOne", func(t *testing.T) {
		// Insert a document first
		doc := bson.M{"name": "update_test", "status": "pending"}
		insertResult, err := conn.InsertOne(ctx, testCol, doc)
		require.NoError(t, err)

		// Update the document
		filter := bson.M{"_id": insertResult.InsertedID}
		update := bson.M{"$set": bson.M{"status": "completed"}}
		updateResult, err := conn.UpdateOne(ctx, testCol, filter, update)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), updateResult.MatchedCount)
		assert.Equal(t, int64(1), updateResult.ModifiedCount)

		// Verify the update
		var updatedDoc bson.M
		err = conn.FindOne(ctx, testCol, filter).Decode(&updatedDoc)
		assert.NoError(t, err)
		assert.Equal(t, "completed", updatedDoc["status"])
	})

	t.Run("Find multiple documents", func(t *testing.T) {
		// Insert multiple documents
		docs := []interface{}{
			bson.M{"category": "test", "value": 1},
			bson.M{"category": "test", "value": 2},
			bson.M{"category": "test", "value": 3},
		}

		for _, doc := range docs {
			_, err := conn.InsertOne(ctx, testCol, doc)
			require.NoError(t, err)
		}

		// Find all documents with category "test"
		cursor, err := conn.Find(ctx, testCol, bson.M{"category": "test"})
		assert.NoError(t, err)
		assert.NotNil(t, cursor)

		var results []bson.M
		err = cursor.All(ctx, &results)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 3)
	})

	t.Run("CountDocuments", func(t *testing.T) {
		count, err := conn.CountDocuments(ctx, testCol, bson.M{"category": "test"})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})

	t.Run("UpdateMany", func(t *testing.T) {
		// Update all documents with category "test"
		filter := bson.M{"category": "test"}
		update := bson.M{"$set": bson.M{"updated": true}}
		result, err := conn.UpdateMany(ctx, testCol, filter, update)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, result.MatchedCount, int64(3))
		assert.GreaterOrEqual(t, result.ModifiedCount, int64(3))
	})

	t.Run("DeleteOne", func(t *testing.T) {
		// Insert a document to delete
		doc := bson.M{"name": "delete_test", "temp": true}
		insertResult, err := conn.InsertOne(ctx, testCol, doc)
		require.NoError(t, err)

		// Delete the document
		filter := bson.M{"_id": insertResult.InsertedID}
		deleteResult, err := conn.DeleteOne(ctx, testCol, filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteResult.DeletedCount)

		// Verify deletion
		var deletedDoc bson.M
		err = conn.FindOne(ctx, testCol, filter).Decode(&deletedDoc)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("Aggregate", func(t *testing.T) {
		// Simple aggregation pipeline
		pipeline := []bson.M{
			{"$match": bson.M{"category": "test"}},
			{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
		}

		cursor, err := conn.Aggregate(ctx, testCol, pipeline)
		assert.NoError(t, err)
		assert.NotNil(t, cursor)

		var results []bson.M
		err = cursor.All(ctx, &results)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
	})

	t.Run("BulkWrite", func(t *testing.T) {
		models := []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.M{"bulk": "insert1"}),
			mongo.NewInsertOneModel().SetDocument(bson.M{"bulk": "insert2"}),
		}

		result, err := conn.BulkWrite(ctx, testCol, models)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), result.InsertedCount)
	})
}