package role

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

// CreateRoute registers the routes for the role
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Post(routePath, Create)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create Role",
		Description: "Create a new role",
		RequestBody: &docs.ApiRequestBody{
			Description: "Role data",
			Content:     new(sdk.Role),
		},
		Response: &docs.ApiResponse{
			Description: "Role created successfully",
			Content:     new(sdk.RoleResponse),
		},
		Tags: routeTags,
	})
}

// Create handles the creation of a new role
func Create(c *fiber.Ctx) error {
	log.Debug("received create role request")
	payload := new(sdk.Role)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}
	log.Debug("parsed create role request")

	pr := providers.GetProviders(c)
	err := pr.S.Role.Create(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to create role. %w", err).Error()
		log.Errorw("failed to create role", "error", err)
		return c.Status(status).JSON(sdk.RoleResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("role created successfully")

	return c.Status(http.StatusCreated).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role created successfully",
		Data:    payload,
	})
}

// SearchRoute registers the route for searching roles
func SearchRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Get(routePath, Search)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Search Roles",
		Description: "Search for roles",
		Response: &docs.ApiResponse{
			Description: "Roles fetched successfully",
			Content:     new(sdk.RoleListResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "query",
				In:          "query",
				Description: "Search query for roles",
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
}

// Search searches for roles based on the given criteria
func Search(c *fiber.Ctx) error {
	log.Debug("received search role request")

	// Parse search criteria from query parameters
	query := sdk.RoleQuery{
		SearchQuery: c.Query("query"),
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
	ds, err := pr.S.Role.GetAll(c.Context(), query)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to search roles. %w", err).Error()
		log.Error("failed to search roles", "error", err)
		return c.Status(status).JSON(sdk.RoleListResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("roles searched successfully")
	return c.Status(http.StatusOK).JSON(sdk.RoleListResponse{
		Success: true,
		Message: "Roles searched successfully",
		Data:    ds,
	})
}

// GetRoute registers the route for fetching a specific role
func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Get(routePath, Get)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Role",
		Description: "Get a role by ID",
		Response: &docs.ApiResponse{
			Description: "Role fetched successfully",
			Content:     new(sdk.RoleResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the role",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

// Get retrieves a specific role by ID
func Get(c *fiber.Ctx) error {
	log.Debug("received get role request")
	id := c.Params("id")

	pr := providers.GetProviders(c)
	ds, err := pr.S.Role.GetById(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get role. %w", err).Error()
		if errors.Is(err, sdk.ErrRoleNotFound) {
			status = http.StatusNotFound
			message = "role not found"
		}
		log.Error("failed to get role", "error", message)
		return c.Status(status).JSON(sdk.RoleResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("role fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role fetched successfully",
		Data:    ds,
	})
}

// UpdateRoute registers the route for updating a role
func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Put(routePath, Update)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update Role",
		Description: "Update an existing role",
		RequestBody: &docs.ApiRequestBody{
			Description: "Role data",
			Content:     new(sdk.Role),
		},
		Response: &docs.ApiResponse{
			Description: "Role updated successfully",
			Content:     new(sdk.RoleResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the role",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

// Update modifies an existing role
func Update(c *fiber.Ctx) error {
	log.Debug("received update role request")
	id := c.Params("id")

	payload := new(sdk.Role)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update role request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}

	payload.Id = id
	pr := providers.GetProviders(c)
	err := pr.S.Role.Update(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to update role. %w", err).Error()
		if errors.Is(err, sdk.ErrRoleNotFound) {
			status = http.StatusNotFound
			message = "role not found"
		}
		log.Error("failed to update role", "error", err)
		return c.Status(status).JSON(sdk.RoleResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("role updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role updated successfully",
		Data:    payload,
	})
}
