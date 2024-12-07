package project

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/api-server/providers"
	"github.com/melvinodsa/go-iam/api-server/sdk"
	"github.com/melvinodsa/go-iam/api-server/services/project"
)

func Create(c *fiber.Ctx) error {
	log.Debug("received create project request")
	payload := new(sdk.Project)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.ProjectResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}
	log.Debug("parsed create project request")
	pr := providers.GetProviders(c)
	err := pr.S.Projects.Create(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to create project. %w", err).Error()
		log.Errorw("failed to create project", "error", err)
		return c.Status(status).JSON(sdk.ProjectResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("project created successfully")

	return c.Status(http.StatusOK).JSON(sdk.ProjectResponse{
		Success: true,
		Message: "Project created successfully",
		Data:    payload,
	})
}

// Get project
func Get(c *fiber.Ctx) error {
	log.Debug("received get project request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get project request. project id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.ProjectResponse{
			Success: false,
			Message: "Invalid request. Project id is required",
		})
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.Projects.Get(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get project. %w", err).Error()
		if errors.Is(err, project.ErrProjectNotFound) {
			status = http.StatusBadRequest
			message = "project not found"
		}
		log.Error("failed to get project", "error", message)
		return c.Status(status).JSON(sdk.ProjectResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("project fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.ProjectResponse{
		Success: true,
		Message: "Project fetched successfully",
		Data:    ds,
	})
}

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get projects request")
	pr := providers.GetProviders(c)
	ds, err := pr.S.Projects.GetAll(c.Context())
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get projects. %w", err).Error()
		log.Error("failed to get projects", "error", err)
		return c.Status(status).JSON(sdk.ProjectsResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("projects fetched successfully")

	return c.Status(http.StatusOK).JSON(sdk.ProjectsResponse{
		Success: true,
		Message: "Projects fetched successfully",
		Data:    ds,
	})
}

// Update project
func Update(c *fiber.Ctx) error {
	log.Debug("received update project request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update project request. project id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.ProjectResponse{
			Success: false,
			Message: "Invalid request. Project id is required",
		})
	}
	payload := new(sdk.Project)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update project request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.ProjectResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}

	payload.Id = id
	pr := providers.GetProviders(c)
	err := pr.S.Projects.Update(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to update project. %w", err).Error()
		if errors.Is(err, project.ErrProjectNotFound) {
			status = http.StatusBadRequest
			message = "project not found"
		}
		log.Error("failed to update project", "error", err)
		return c.Status(status).JSON(sdk.ProjectResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("project updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.ProjectResponse{
		Success: true,
		Message: "Project updated successfully",
		Data:    payload,
	})
}
