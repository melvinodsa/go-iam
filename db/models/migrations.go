package models

import "time"

// Migration represents a database migration record in the Go IAM system.
// Migrations track schema changes and data transformations applied to the database.
// This ensures database consistency across different environments and deployments.
type Migration struct {
	Id          string     `bson:"id"`          // Unique identifier for the migration
	Version     string     `bson:"version"`     // Version identifier of the migration
	Name        string     `bson:"name"`        // Human-readable name of the migration
	Description string     `bson:"description"` // Detailed description of what the migration does
	AppliedAt   *time.Time `bson:"applied_at"`  // Timestamp when the migration was applied
	Checksum    string     `bson:"checksum"`    // Checksum to verify migration integrity
	CreatedAt   *time.Time `bson:"created_at"`  // Timestamp when the migration record was created
	CreatedBy   string     `bson:"created_by"`  // User or system that created the migration record
	UpdatedAt   *time.Time `bson:"updated_at"`  // Timestamp when the migration record was last updated
	UpdatedBy   string     `bson:"updated_by"`  // User or system that last updated the migration record
}

// MigrationModel provides database access patterns and field mappings for Migration entities.
// It embeds the iam struct to inherit the database name and implements collection operations.
type MigrationModel struct {
	iam                   // Embedded struct providing DbName() method
	IdKey          string // BSON field key for migration ID
	VersionKey     string // BSON field key for migration version
	NameKey        string // BSON field key for migration name
	DescriptionKey string // BSON field key for migration description
	AppliedAtKey   string // BSON field key for application timestamp
	ChecksumKey    string // BSON field key for migration checksum
}

// Name returns the MongoDB collection name for migrations.
// This implements the DbCollection interface.
func (m MigrationModel) Name() string {
	return "migrations"
}

// GetMigrationModel returns a properly initialized MigrationModel with all field mappings.
// This function provides a singleton pattern for accessing migration model operations.
//
// Returns a MigrationModel instance with all BSON field keys mapped to their respective field names.
func GetMigrationModel() MigrationModel {
	return MigrationModel{
		IdKey:          "id",
		VersionKey:     "version",
		NameKey:        "name",
		DescriptionKey: "description",
		AppliedAtKey:   "applied_at",
		ChecksumKey:    "checksum",
	}
}
