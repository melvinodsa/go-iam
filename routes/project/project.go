package project

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils/docs"
)

func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create Project",
		Description: "Create a new project",
		RequestBody: &docs.ApiRequestBody{
			Description: "Project data",
			Content:     new(sdk.Project),
		},
		Response: &docs.ApiResponse{
			Description: "Project created successfully",
			Content:     new(sdk.ProjectResponse),
		},
		ProjectIDNotRequired: true,
		Tags:                 routeTags,
	})
	router.Post(routePath, Create)
}

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

func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Project",
		Description: "Get a project by ID",
		Response: &docs.ApiResponse{
			Description: "Project fetched successfully",
			Content:     new(sdk.ProjectResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the project",
				Required:    true,
			},
		},
		ProjectIDNotRequired: true,
		Tags:                 routeTags,
	})
	router.Get(routePath, Get)
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

func FetchAllRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Fetch All Projects",
		Description: "Fetch all projects",
		Response: &docs.ApiResponse{
			Description: "Projects fetched successfully",
			Content:     new(sdk.ProjectsResponse),
		},
		ProjectIDNotRequired: true,
		Tags:                 routeTags,
	})
	router.Get(routePath, FetchAll)
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

func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update Project",
		Description: "Update a project by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "Project data",
			Content:     new(sdk.Project),
		},
		Response: &docs.ApiResponse{
			Description: "Project updated successfully",
			Content:     new(sdk.ProjectResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the project",
				Required:    true,
			},
		},
		ProjectIDNotRequired: true,
		Tags:                 routeTags,
	})
	router.Put(routePath, Update)
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
