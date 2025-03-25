package resource

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/resource"
)

// Create handles the creation of a new resource
func Create(c *fiber.Ctx) error {
	log.Debug("received create resource request")
	payload := new(sdk.Resource)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.ResourceResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}
	log.Debug("parsed create resource request")

	pr := providers.GetProviders(c)
	err := pr.S.Resources.Create(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to create resource. %w", err).Error()
		log.Errorw("failed to create resource", "error", err)
		return c.Status(status).JSON(sdk.ResourceResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("resource created successfully")

	return c.Status(http.StatusOK).JSON(sdk.ResourceResponse{
		Success: true,
		Message: "Resource created successfully",
		Data:    payload,
	})
}

func Get(c *fiber.Ctx) error {
	log.Debug("received get resource request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get resource request. resource id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.ResourceResponse{
			Success: false,
			Message: "Invalid request. Resource id is required",
		})
	}

	pr := providers.GetProviders(c)
	ds, err := pr.S.Resources.Get(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get resource. %w", err).Error()
		if errors.Is(err, resource.ErrResourceNotFound) {
			status = http.StatusBadRequest
			message = "resource not found"
		}
		log.Error("failed to get resource", "error", message)
		return c.Status(status).JSON(sdk.ResourceResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("resource fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.ResourceResponse{
		Success: true,
		Message: "Resource fetched successfully",
		Data:    ds,
	})
}

// Search searches for resources based on the given criteria
func Search(c *fiber.Ctx) error {
	log.Debug("received search resources request")

	// Parse search criteria from query parameters
	query := sdk.ResourceQuery{
		Name:        c.Query("name"),
		Description: c.Query("description"),
		Skip:        0,  // Default value
		Limit:       10, // Default value
	}

	// Parse pagination parameters if provided
	if skip := c.Query("skip"); skip != "" {
		if val, err := strconv.ParseInt(skip, 10, 64); err == nil {
			query.Skip = val
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.ParseInt(limit, 10, 64); err == nil {
			query.Limit = val
		}
	}

	pr := providers.GetProviders(c)
	ds, err := pr.S.Resources.Search(c.Context(), query)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to search resources. %w", err).Error()
		log.Error("failed to search resources", "error", err)
		return c.Status(status).JSON(sdk.ResourcesResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("resources searched successfully")
	return c.Status(http.StatusOK).JSON(sdk.ResourcesResponse{
		Success: true,
		Message: "Resources searched successfully",
		Data:    ds,
	})
}

// Update modifies an existing resource
func Update(c *fiber.Ctx) error {
	log.Debug("received update resource request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update resource request. resource id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.ResourceResponse{
			Success: false,
			Message: "Invalid request. Resource id is required",
		})
	}

	payload := new(sdk.Resource)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update resource request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.ResourceResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}

	payload.ID = id
	pr := providers.GetProviders(c)
	err := pr.S.Resources.Update(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to update resource. %w", err).Error()
		if errors.Is(err, resource.ErrResourceNotFound) {
			status = http.StatusBadRequest
			message = "resource not found"
		}
		log.Error("failed to update resource", "error", err)
		return c.Status(status).JSON(sdk.ResourceResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("resource updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.ResourceResponse{
		Success: true,
		Message: "Resource updated successfully",
		Data:    payload,
	})
}

func Delete(c *fiber.Ctx) error {
	log.Debug("received delete resource request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid delete resource request. resource id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.ResourceResponse{
			Success: false,
			Message: "Invalid request. Resource id is required",
		})
	}

	pr := providers.GetProviders(c)
	err := pr.S.Resources.Delete(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to delete resource. %w", err).Error()
		if errors.Is(err, resource.ErrResourceNotFound) {
			status = http.StatusNotFound
			message = "resource not found"
		}
		log.Error("failed to delete resource", "error", err)
		return c.Status(status).JSON(sdk.ResourceResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("resource deleted successfully")
	return c.Status(http.StatusOK).JSON(sdk.ResourceResponse{
		Success: true,
		Message: "Resource deleted successfully",
	})
}
