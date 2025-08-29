package db

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Test helper functions
func setupMockDB() *MockDB {
	return new(MockDB)
}

func createTestMigration(version, name string, shouldFail bool) MigrationInfo {
	return MigrationInfo{
		Version:     version,
		Name:        name,
		Description: "Test migration " + version,
		Up: func(ctx context.Context, db DB) error {
			if shouldFail {
				return errors.New("migration failed")
			}
			return nil
		},
		Down: func(ctx context.Context, db DB) error {
			if shouldFail {
				return errors.New("rollback failed")
			}
			return nil
		},
	}
}

func TestRegisterMigration(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Reset migrations for test
	registeredMigrations = []MigrationInfo{}

	migration := createTestMigration("001", "test_migration", false)
	RegisterMigration(migration)

	assert.Len(t, registeredMigrations, 1)
	assert.Equal(t, "001", registeredMigrations[0].Version)
	assert.Equal(t, "test_migration", registeredMigrations[0].Name)
}

func TestCheckAndRunMigrations_NewMigration(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register test migration
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "test_migration", false),
	}

	// Mock FindOne to return ErrNoDocuments (migration not found)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult)

	// Mock InsertOne for recording migration
	insertResult := &mongo.InsertOneResult{InsertedID: "001"}
	mockDB.On("InsertOne", ctx, migrationModel, mock.AnythingOfType("models.Migration"), mock.Anything).Return(insertResult, nil)

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestCheckAndRunMigrations_ExistingMigration(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register test migration
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "test_migration", false),
	}

	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{}, nil, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult)

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	// InsertOne should not be called since migration already exists
	mockDB.AssertNotCalled(t, "InsertOne")
}

func TestCheckAndRunMigrations_MigrationUpFails(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register test migration that fails
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "failing_migration", true),
	}

	// Mock FindOne to return ErrNoDocuments (migration not found)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult)

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to apply migration 001")
	mockDB.AssertExpectations(t)
	// InsertOne should not be called since migration failed
	mockDB.AssertNotCalled(t, "InsertOne")
}

func TestCheckAndRunMigrations_FindOneFails(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register test migration
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "test_migration", false),
	}

	// Mock FindOne to return a different error
	mockResult := mongo.NewSingleResultFromDocument(bson.D{}, errors.New("database error"), nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult)

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check migration 001: database error")
	mockDB.AssertExpectations(t)
}

func TestCheckAndRunMigrations_InsertOneFails(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register test migration
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "test_migration", false),
	}

	// Mock FindOne to return ErrNoDocuments (migration not found)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult)

	// Mock InsertOne to fail
	mockDB.On("InsertOne", ctx, migrationModel, mock.AnythingOfType("models.Migration"), mock.Anything).Return((*mongo.InsertOneResult)(nil), errors.New("insert failed"))

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to record migration 001: insert failed")
	mockDB.AssertExpectations(t)
}

func TestCheckAndRunMigrations_MultipleMigrations(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test
	mockDB := setupMockDB()
	ctx := context.Background()
	migrationModel := models.GetMigrationModel()

	// Register multiple test migrations
	registeredMigrations = []MigrationInfo{
		createTestMigration("001", "first_migration", false),
		createTestMigration("002", "second_migration", false),
	}

	// Mock FindOne for first migration (not found)
	mockResult1 := mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "001"}, mock.Anything).Return(mockResult1)

	// Mock FindOne for second migration (already exists)
	mockResult2 := mongo.NewSingleResultFromDocument(bson.D{}, nil, nil)
	mockDB.On("FindOne", ctx, migrationModel, bson.M{migrationModel.VersionKey: "002"}, mock.Anything).Return(mockResult2)

	// Mock InsertOne for first migration only
	insertResult := &mongo.InsertOneResult{InsertedID: "001"}
	mockDB.On("InsertOne", ctx, migrationModel, mock.AnythingOfType("models.Migration"), mock.Anything).Return(insertResult, nil)

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestIsMigrationApplied_Exists(t *testing.T) {
	mockDB := setupMockDB()
	ctx := context.Background()
	version := "001"
	migrationModel := models.GetMigrationModel()

	// Mock CountDocuments to return 1 (migration exists)
	mockDB.On("CountDocuments", ctx, migrationModel, bson.M{migrationModel.VersionKey: version}, mock.Anything).Return(int64(1), nil)

	// Execute
	applied, err := IsMigrationApplied(ctx, mockDB, version)

	// Assert
	assert.NoError(t, err)
	assert.True(t, applied)
	mockDB.AssertExpectations(t)
}

func TestIsMigrationApplied_NotExists(t *testing.T) {
	mockDB := setupMockDB()
	ctx := context.Background()
	version := "001"
	migrationModel := models.GetMigrationModel()

	// Mock CountDocuments to return 0 (migration does not exist)
	mockDB.On("CountDocuments", ctx, migrationModel, bson.M{migrationModel.VersionKey: version}, mock.Anything).Return(int64(0), nil)

	// Execute
	applied, err := IsMigrationApplied(ctx, mockDB, version)

	// Assert
	assert.NoError(t, err)
	assert.False(t, applied)
	mockDB.AssertExpectations(t)
}

func TestIsMigrationApplied_Error(t *testing.T) {
	mockDB := setupMockDB()
	ctx := context.Background()
	version := "001"
	migrationModel := models.GetMigrationModel()

	// Mock CountDocuments to return error
	mockDB.On("CountDocuments", ctx, migrationModel, bson.M{migrationModel.VersionKey: version}, mock.Anything).Return(int64(0), errors.New("database error"))

	// Execute
	applied, err := IsMigrationApplied(ctx, mockDB, version)

	// Assert
	assert.Error(t, err)
	assert.False(t, applied)
	assert.Contains(t, err.Error(), "failed to check migration 001")
	mockDB.AssertExpectations(t)
}

func TestCheckAndRunMigrations_EmptyMigrations(t *testing.T) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	// Setup test with no migrations
	mockDB := setupMockDB()
	ctx := context.Background()
	registeredMigrations = []MigrationInfo{}

	// Execute
	err := CheckAndRunMigrations(ctx, mockDB)

	// Assert
	assert.NoError(t, err)
	// No database calls should be made
	mockDB.AssertNotCalled(t, "FindOne")
	mockDB.AssertNotCalled(t, "InsertOne")
}

func TestMigrationInfo_Structure(t *testing.T) {
	migration := MigrationInfo{
		Version:     "001",
		Name:        "test_migration",
		Description: "A test migration",
		Up: func(ctx context.Context, db DB) error {
			return nil
		},
		Down: func(ctx context.Context, db DB) error {
			return nil
		},
	}

	assert.Equal(t, "001", migration.Version)
	assert.Equal(t, "test_migration", migration.Name)
	assert.Equal(t, "A test migration", migration.Description)
	assert.NotNil(t, migration.Up)
	assert.NotNil(t, migration.Down)

	// Test that functions can be called
	err := migration.Up(context.Background(), nil)
	assert.NoError(t, err)

	err = migration.Down(context.Background(), nil)
	assert.NoError(t, err)
}

// Benchmark tests
func BenchmarkRegisterMigration(b *testing.B) {
	// Save original state
	originalMigrations := registeredMigrations
	defer func() {
		registeredMigrations = originalMigrations
	}()

	migration := createTestMigration("001", "benchmark_migration", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registeredMigrations = []MigrationInfo{} // Reset for each iteration
		RegisterMigration(migration)
	}
}

func BenchmarkIsMigrationApplied(b *testing.B) {
	mockDB := setupMockDB()
	ctx := context.Background()
	version := "001"
	migrationModel := models.GetMigrationModel()

	mockDB.On("CountDocuments", ctx, migrationModel, mock.Anything, mock.Anything).Return(int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsMigrationApplied(ctx, mockDB, version)
	}
}
