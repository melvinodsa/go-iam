// Package models provides database model definitions and access patterns
// for the Go IAM system. All models implement the DbCollection interface
// and provide BSON field mappings for MongoDB operations.
package models

// iam is an embedded struct that provides common database configuration
// for all Go IAM models. It implements the DbName() method required
// by the DbCollection interface.
type iam struct{}

// DbName returns the MongoDB database name used by all Go IAM models.
// This implements the DbCollection interface requirement.
func (i iam) DbName() string {
	return "iam"
}
