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

// MigrationFunc represents a function that performs a migration
type MigrationFunc func(ctx context.Context, db DB) error

// MigrationInfo contains metadata about a migration
type MigrationInfo struct {
	Version     string
	Name        string
	Description string
	Up          MigrationFunc
	Down        MigrationFunc
}

var registeredMigrations []MigrationInfo

// RegisterMigration registers a new migration
func RegisterMigration(migration MigrationInfo) {
	registeredMigrations = append(registeredMigrations, migration)
}

// CheckAndRunMigrations checks if migrations need to be run and executes them
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

// IsMigrationApplied checks if a specific migration has been applied
func IsMigrationApplied(ctx context.Context, db DB, version string) (bool, error) {
	migrationModel := models.GetMigrationModel()
	filter := bson.M{migrationModel.VersionKey: version}

	count, err := db.CountDocuments(ctx, migrationModel, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check migration %s: %w", version, err)
	}

	return count > 0, nil
}

func GetMigrations() []MigrationInfo {
	return registeredMigrations
}

func ResetMigrations() {
	registeredMigrations = []MigrationInfo{}
}
