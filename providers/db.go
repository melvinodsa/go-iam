package providers

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"

	_ "github.com/melvinodsa/go-iam/db/migrations"
)

// NewDBConnection establishes a MongoDB connection and applies database migrations.
// This function initializes the database connection using the provided configuration
// and ensures the database schema is up to date by running all pending migrations.
//
// The function performs the following operations:
// 1. Creates MongoDB connection using the configured host
// 2. Applies all registered migrations to ensure schema consistency
// 3. Returns the database interface for use by services
//
// Parameters:
//   - cnf: Application configuration containing database settings
//
// Returns:
//   - db.DB: Database interface for service operations
//   - error: Error if connection fails or migrations cannot be applied
func NewDBConnection(cnf config.AppConfig) (db.DB, error) {
	conn, err := db.NewMongoConnection(cnf.DB.Host())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	// apply migrations
	if err := db.CheckAndRunMigrations(context.Background(), conn); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}
	return conn, nil
}
