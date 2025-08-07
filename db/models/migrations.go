package models

import "time"

type Migration struct {
	Id          string     `bson:"id"`
	Version     string     `bson:"version"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	AppliedAt   *time.Time `bson:"applied_at"`
	Checksum    string     `bson:"checksum"`
	CreatedAt   *time.Time `bson:"created_at"`
	CreatedBy   string     `bson:"created_by"`
	UpdatedAt   *time.Time `bson:"updated_at"`
	UpdatedBy   string     `bson:"updated_by"`
}

type MigrationModel struct {
	iam
	IdKey          string
	VersionKey     string
	NameKey        string
	DescriptionKey string
	AppliedAtKey   string
	ChecksumKey    string
}

func (m MigrationModel) Name() string {
	return "migrations"
}

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
