package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
)

// CheckAndAddDefaultProject ensures that a default project exists in the system.
// This function is called during system initialization to create a default project
// if one doesn't already exist. The default project provides a basic organizational
// unit for initial system setup and testing.
//
// The function performs the following operations:
// 1. Checks if a project named "Default Project" already exists
// 2. If not found, creates a new default project with system tags
// 3. Logs the creation process for audit purposes
//
// Default project characteristics:
// - Name: "Default Project"
// - Description: System-generated description
// - Tags: ["default", "system"]
// - Created by: "system"
//
// Parameters:
//   - svc: Project service for project operations
//
// Returns:
//   - error: Error if project lookup or creation fails
func CheckAndAddDefaultProject(svc project.Service) error {
	_, err := svc.GetByName(context.Background(), "Default Project")
	if err != nil && err != sdk.ErrProjectNotFound {
		return fmt.Errorf("error fetching default project: %w", err)
	}
	if err == nil {
		return nil // Default project already exists
	}
	log.Info("Default project not found, creating it...")

	t := time.Now()
	defaultProject := &sdk.Project{
		Name:        "Default Project",
		Description: "This is the default project created by go-iam.",
		CreatedAt:   &t,
		CreatedBy:   "system",
		UpdatedAt:   &t,
		UpdatedBy:   "system",
		Tags:        []string{"default", "system"},
	}

	err = svc.Create(context.Background(), defaultProject)
	if err != nil {
		return fmt.Errorf("error creating default project: %w", err)
	}
	log.Info("Default project created successfully.")
	return nil
}
