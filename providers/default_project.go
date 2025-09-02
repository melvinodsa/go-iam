package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
)

func CheckAndAddDefaultProject(svc project.Service) error {
	_, err := svc.GetByName(context.Background(), "Default Project")
	if err != nil && err != project.ErrProjectNotFound {
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
