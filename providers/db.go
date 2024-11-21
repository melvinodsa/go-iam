package providers

import (
	"fmt"

	"github.com/melvinodsa/go-iam/api-server/config"
	"github.com/melvinodsa/go-iam/api-server/db"
)

func NewDBConnection(cnf config.AppConfig) (db.DB, error) {
	conn, err := db.NewMongoConnection(cnf.DB.Host())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	return conn, nil
}
