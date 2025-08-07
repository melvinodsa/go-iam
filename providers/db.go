package providers

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
)

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
