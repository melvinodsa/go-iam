package policy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/policybeta"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// CreateRoute registers the routes for the policy
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create Policy",
		Description: "Create a new policy",
		RequestBody: &docs.ApiRequestBody{
			Description: "Policy data",
			Content:     new(sdk.Policy),
		},
		Response: &docs.ApiResponse{
			Description: "Policy created successfully",
			Content:     new(sdk.PolicyResponse),
		},
		Tags: routeTags,
	})
}

func Create(c *fiber.Ctx) error {
	log.Debug("received create policy request")
	payload := new(sdk.Policy)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}
	log.Debug("parsed create policy request")
	pr := providers.GetProviders(c)
	err := pr.S.Policy.Create(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to create policy. %w", err).Error()
		log.Errorw("failed to create policy", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("policy created successfully")

	return c.Status(http.StatusOK).JSON(sdk.PolicyResponse{
		Success: true,
		Message: "Policy created successfully",
		Data:    payload,
	})
}

// Get retrieves a policy by its ID
func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Policy",
		Description: "Get a policy by ID",
		Response: &docs.ApiResponse{
			Description: "Policy fetched successfully",
			Content:     new(sdk.PolicyResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the policy",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func Get(c *fiber.Ctx) error {
	log.Debug("received get policy request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get policy request. policy id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: "Invalid request. Policy id is required",
		})
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.Policy.Get(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get policy. %w", err).Error()
		if errors.Is(err, policybeta.ErrPolicyNotFound) {
			status = http.StatusBadRequest
			message = "policy not found"
		}
		log.Error("failed to get policy", "error", message)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("policy fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.PolicyResponse{
		Success: true,
		Message: "Policy fetched successfully",
		Data:    ds,
	})
}

// FetchAll retrieves all policies
func FetchAllRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Fetch All Policies",
		Description: "Fetch all policies",
		Response: &docs.ApiResponse{
			Description: "Policies fetched successfully",
			Content:     new(sdk.PoliciesResponse),
		},
		Tags: routeTags,
	})
}

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get Policy request")
	pr := providers.GetProviders(c)
	ds, err := pr.S.Policy.GetAll(c.Context())
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get Policy. %w", err).Error()
		log.Error("failed to get Policy", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("Policy fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.PoliciesResponse{
		Success: true,
		Message: "Policy fetched successfully",
		Data:    ds,
	})
}

// UpdateRoute registers the route for updating a policy
func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update Policy",
		Description: "Update a policy by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "Updated policy data",
			Content:     new(sdk.Policy),
		},
		Response: &docs.ApiResponse{
			Description: "Policy updated successfully",
			Content:     new(sdk.PolicyResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the policy",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func Update(c *fiber.Ctx) error {
	log.Debug("received update policy request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update policy request. policy id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: "Invalid request. Policy id is required",
		})
	}
	payload := new(sdk.Policy)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update policy request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}

	payload.Id = id
	pr := providers.GetProviders(c)
	err := pr.S.Policy.Update(c.Context(), payload)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to update policy. %w", err).Error()
		if errors.Is(err, policybeta.ErrPolicyNotFound) {
			status = http.StatusBadRequest
			message = "policy not found"
		}
		log.Error("failed to update policy", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("policy updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.PolicyResponse{
		Success: true,
		Message: "Policy updated successfully",
		Data:    payload,
	})
}

// DeleteRoute registers the route for deleting a policy
func DeleteRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodDelete,
		Name:        "Delete Policy",
		Description: "Delete a policy by ID",
		Response: &docs.ApiResponse{
			Description: "Policy deleted successfully",
			Content:     new(sdk.PolicyResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the policy",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func Delete(c *fiber.Ctx) error {
	log.Debug("received delete policy request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid delete policy request. policy id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: "Invalid request. Policy id is required",
		})
	}

	pr := providers.GetProviders(c)
	err := pr.S.Policy.Delete(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to delete policy. %w", err).Error()
		if errors.Is(err, policybeta.ErrPolicyNotFound) {
			status = http.StatusBadRequest
			message = "policy not found"
		}
		log.Error("failed to delete policy", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("policy deleted successfully")
	return c.Status(http.StatusOK).JSON(sdk.PolicyResponse{
		Success: true,
		Message: "Policy deleted successfully",
	})
}

// GetPoliciesByRoleIdRoute registers the route to get policies by role ID
func GetPoliciesByRoleIdRoute(router fiber.Router, basePath string) {
	routePath := "/role/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Policies by Role ID",
		Description: "Fetch policies associated with a specific role ID",
		Response: &docs.ApiResponse{
			Description: "Policies fetched successfully",
			Content:     new(sdk.PoliciesResponse),
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

func GetPoliciesByRoleId(c *fiber.Ctx) error {
	log.Debug("received get policies by role id request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get policies by role id request. role id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: "Invalid request. Role id is required",
		})
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.Policy.GetPoliciesByRoleId(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get policies by role id. %w", err).Error()
		log.Error("failed to get policies by role id", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("policies fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.PoliciesResponse{
		Success: true,
		Message: "Policies fetched successfully",
		Data:    ds,
	})
}

// SyncResourcesRoute registers the route for syncing resources by policy ID
func SyncResourcesRoute(router fiber.Router, basePath string) {
	routePath := "/sync"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Sync Resources",
		Description: "Sync resources by policy ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "Sync resources request",
			Content:     new(syncResourcesRequest),
		},
		Response: &docs.ApiResponse{
			Description: "Resources synced successfully",
			Content:     new(sdk.PolicyResponse),
		},
		Tags: routeTags,
	})
}

type syncResourcesRequest struct {
	Policies   map[string]string `json:"policies"`
	ResourceId string            `json:"resourceId"`
	Name       string            `json:"name"`
}

func SyncResources(c *fiber.Ctx) error {
	log.Debug("received sync resources request")
	payload := new(syncResourcesRequest)
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.PolicyResponse{
			Success: false,
			Message: fmt.Errorf("invalid request. %w", err).Error(),
		})
	}

	pr := providers.GetProviders(c)
	err := pr.S.Policy.SyncResourcesbyPolicyId(c.Context(), payload.Policies, payload.ResourceId, payload.Name)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to sync resources. %w", err).Error()
		log.Errorw("failed to sync resources", "error", err)
		return c.Status(status).JSON(sdk.PolicyResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("resources synced successfully")

	return c.Status(http.StatusOK).JSON(sdk.PolicyResponse{
		Success: true,
		Message: "Resources synced successfully",
	})
}
