package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MigrationFunc represents a function that performs a database migration.
// It receives a context and database interface to execute migration operations.
type MigrationFunc func(ctx context.Context, db DB) error

// MigrationInfo contains metadata and implementation for a database migration.
// Each migration should have a unique version identifier and both up/down functions.
type MigrationInfo struct {
	Version     string        // Unique version identifier for the migration
	Name        string        // Human-readable name of the migration
	Description string        // Detailed description of what the migration does
	Up          MigrationFunc // Function to apply the migration
	Down        MigrationFunc // Function to rollback the migration
}

var registeredMigrations []MigrationInfo

// RegisterMigration registers a new migration to be tracked and executed.
// Migrations should be registered in the order they should be applied.
//
// Parameters:
//   - migration: MigrationInfo containing all migration details
func RegisterMigration(migration MigrationInfo) {
	registeredMigrations = append(registeredMigrations, migration)
}

// CheckAndRunMigrations checks if migrations need to be run and executes them.
// This function should be called during application startup to ensure the database
// schema is up to date. It applies migrations in the order they were registered.
//
// Parameters:
//   - ctx: Context for the migration operations
//   - db: Database interface to execute migrations against
//
// Returns an error if any migration fails to apply or if there are database issues.
// Successfully applied migrations are recorded in the migrations collection.
func CheckAndRunMigrations(ctx context.Context, db DB) error {
	migrationModel := models.GetMigrationModel()

	log.Info("Checking database migrations...")

	for _, migration := range registeredMigrations {
		// Check if migration has already been applied
		filter := bson.M{migrationModel.VersionKey: migration.Version}
		result := db.FindOne(ctx, migrationModel, filter)

		var existingMigration models.Migration
		err := result.Decode(&existingMigration)

		if errors.Is(err, mongo.ErrNoDocuments) {
			// Migration not found, need to apply it
			log.Infof("Applying migration %s: %s", migration.Version, migration.Name)

			if err := migration.Up(ctx, db); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}

			// Record the migration as applied
			now := time.Now()
			migrationRecord := models.Migration{
				Id:          migration.Version,
				Version:     migration.Version,
				Name:        migration.Name,
				Description: migration.Description,
				AppliedAt:   &now,
				CreatedAt:   &now,
				CreatedBy:   "system",
				UpdatedAt:   &now,
				UpdatedBy:   "system",
			}

			if _, err := db.InsertOne(ctx, migrationModel, migrationRecord); err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
			}

			log.Infof("Successfully applied migration %s", migration.Version)
		} else if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.Version, err)
		} else {
			log.Debugf("Migration %s already applied", migration.Version)
		}
	}

	log.Info("All migrations checked and applied successfully")
	return nil
}

// IsMigrationApplied checks if a specific migration has been applied to the database.
// This function can be used to conditionally run migrations or check database state.
//
// Parameters:
//   - ctx: Context for the database operation
//   - db: Database interface to query
//   - version: Version identifier of the migration to check
//
// Returns:
//   - bool: true if the migration has been applied, false otherwise
//   - error: An error if the database query fails
func IsMigrationApplied(ctx context.Context, db DB, version string) (bool, error) {
	migrationModel := models.GetMigrationModel()
	filter := bson.M{migrationModel.VersionKey: version}

	count, err := db.CountDocuments(ctx, migrationModel, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check migration %s: %w", version, err)
	}

	return count > 0, nil
}

// GetMigrations returns a copy of all registered migrations.
// This function is useful for debugging, testing, or displaying migration status.
//
// Returns a slice containing all registered MigrationInfo structs.
func GetMigrations() []MigrationInfo {
	return registeredMigrations
}

// ResetMigrations clears all registered migrations.
// This function is primarily intended for testing purposes to reset migration state.
// Use with caution in production environments.
func ResetMigrations() {
	registeredMigrations = []MigrationInfo{}
}
