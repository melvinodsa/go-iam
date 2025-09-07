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
	"github.com/melvinodsa/go-iam/utils/docs"
)

// CreateRoute registers the routes for the resource
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create Resource",
		Description: "Create a new resource",
		RequestBody: &docs.ApiRequestBody{
			Description: "Resource data",
			Content:     new(sdk.Resource),
		},
		Response: &docs.ApiResponse{
			Description: "Resource created successfully",
			Content:     new(sdk.ResourceResponse),
		},
		Tags: routeTags,
	})
	router.Post(routePath, Create)
}

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

	return c.Status(http.StatusCreated).JSON(sdk.ResourceResponse{
		Success: true,
		Message: "Resource created successfully",
		Data:    payload,
	})
}

// GetRoute registers the route for getting a resource
func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Resource",
		Description: "Get a resource by ID",
		Response: &docs.ApiResponse{
			Description: "Resource fetched successfully",
			Content:     new(sdk.ResourceResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the resource",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
	router.Get(routePath, Get)
}

func Get(c *fiber.Ctx) error {
	log.Debug("received get resource request")
	id := c.Params("id")

	pr := providers.GetProviders(c)
	ds, err := pr.S.Resources.Get(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get resource. %w", err).Error()
		if errors.Is(err, sdk.ErrResourceNotFound) {
			status = http.StatusNotFound
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

// SearchRoute registers the route for searching resources
func SearchRoute(router fiber.Router, basePath string) {
	routePath := "/search"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Search Resources",
		Description: "Search for resources",
		Response: &docs.ApiResponse{
			Description: "Resources fetched successfully",
			Content:     new(sdk.ResourcesResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "name",
				In:          "query",
				Description: "Name of the resource to search for",
				Required:    false,
			},
			{
				Name:        "description",
				In:          "query",
				Description: "Description of the resource to search for",
				Required:    false,
			},
			{
				Name:        "key",
				In:          "query",
				Description: "Key of the resource to search for",
				Required:    false,
			},
			{
				Name:        "skip",
				In:          "query",
				Description: "Number of records to skip for pagination. Default is 0",
				Required:    false,
			},
			{
				Name:        "limit",
				In:          "query",
				Description: "Maximum number of records to return. Default is 10",
				Required:    false,
			},
		},
		Tags: routeTags,
	})
	router.Get(routePath, Search)
}

// Search searches for resources based on the given criteria
func Search(c *fiber.Ctx) error {
	log.Debug("received search resources request")

	// Parse search criteria from query parameters
	query := sdk.ResourceQuery{
		Name:        c.Query("name"),
		Description: c.Query("description"),
		Key:         c.Query("key"),
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

// UpdateRoute registers the route for updating a resource
func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update Resource",
		Description: "Update an existing resource",
		Response: &docs.ApiResponse{
			Description: "Resource updated successfully",
			Content:     new(sdk.ResourceResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the resource",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
	router.Put(routePath, Update)
}

// Update modifies an existing resource
func Update(c *fiber.Ctx) error {
	log.Debug("received update resource request")
	id := c.Params("id")

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
		if errors.Is(err, sdk.ErrResourceNotFound) {
			status = http.StatusNotFound
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

// DeleteRoute registers the route for deleting a resource
func DeleteRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodDelete,
		Name:        "Delete Resource",
		Description: "Delete a resource by ID",
		Response: &docs.ApiResponse{
			Description: "Resource deleted successfully",
			Content:     new(sdk.ResourceResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the resource",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
	router.Delete(routePath, Delete)
}

func Delete(c *fiber.Ctx) error {
	log.Debug("received delete resource request")
	id := c.Params("id")

	pr := providers.GetProviders(c)
	err := pr.S.Resources.Delete(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to delete resource. %w", err).Error()
		if errors.Is(err, sdk.ErrResourceNotFound) {
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
