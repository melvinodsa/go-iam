package authprovider

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// CreateRoute registers the routes for the authprovider
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create AuthProvider",
		Description: "Create a new authprovider",
		Tags:        routeTags,
		RequestBody: &docs.ApiRequestBody{
			Description: "AuthProvider data",
			Content:     new(sdk.AuthProvider),
		},
		Response: &docs.ApiResponse{
			Description: "AuthProvider created successfully",
			Content:     new(sdk.AuthProviderResponse),
		},
	})
	router.Post(routePath, Create)
}

func Create(c *fiber.Ctx) error {
	log.Debug("received create authprovider request")
	payload := new(sdk.AuthProvider)
	if err := c.BodyParser(payload); err != nil {
		return sdk.AuthProviderBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}
	log.Debug("parsed create authprovider request")
	pr := providers.GetProviders(c)
	err := pr.S.AuthProviders.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Errorf("failed to create authprovider. %w", err).Error()
		log.Errorw("failed to create authprovider", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	log.Debug("authprovider created successfully")

	return c.Status(http.StatusOK).JSON(sdk.AuthProviderResponse{
		Success: true,
		Message: "Authprovider created successfully",
		Data:    payload,
	})
}

func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get AuthProvider",
		Description: "Get an authprovider by ID",
		Tags:        routeTags,
		Response: &docs.ApiResponse{
			Description: "AuthProvider fetched successfully",
			Content:     new(sdk.AuthProviderResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the authprovider",
				Required:    true,
			},
		},
	})
	router.Get(routePath, Get)
}

func Get(c *fiber.Ctx) error {
	log.Debug("received get authprovider request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get authprovider request. authprovider id not found")
		return sdk.AuthProviderBadRequest("Invalid request. Authprovider id is required", c)
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.AuthProviders.Get(c.Context(), id, false)
	if err != nil {
		if errors.Is(err, authprovider.ErrAuthProviderNotFound) {
			return sdk.AuthProviderBadRequest("Auth Provider not found", c)
		}
		message := fmt.Errorf("failed to get authprovider. %w", err).Error()
		log.Errorw("failed to get authprovider", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}

	log.Debug("authprovider fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.AuthProviderResponse{
		Success: true,
		Message: "Authprovider fetched successfully",
		Data:    ds,
	})
}

func FetchAllRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Fetch All AuthProviders",
		Description: "Fetch all authproviders",
		Tags:        routeTags,
		Response: &docs.ApiResponse{
			Description: "AuthProviders fetched successfully",
			Content:     new(sdk.AuthProvidersResponse),
		},
	})
	router.Get(routePath, FetchAll)
}

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get authproviders request")
	pr := providers.GetProviders(c)
	ds, err := pr.S.AuthProviders.GetAll(c.Context(), sdk.AuthProviderQueryParams{})
	if err != nil {
		message := fmt.Errorf("failed to get authproviders. %w", err).Error()
		log.Errorw("failed to get authproviders", "error", message)
		return sdk.AuthProvidersInternalServerError(message, c)
	}

	log.Debug("authproviders fetched successfully")

	return c.Status(http.StatusOK).JSON(sdk.AuthProvidersResponse{
		Success: true,
		Message: "Authproviders fetched successfully",
		Data:    ds,
	})
}

func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update AuthProvider",
		Description: "Update an authprovider by ID",
		Tags:        routeTags,
		RequestBody: &docs.ApiRequestBody{
			Description: "AuthProvider data",
			Content:     new(sdk.AuthProvider),
		},
		Response: &docs.ApiResponse{
			Description: "AuthProvider updated successfully",
			Content:     new(sdk.AuthProviderResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the authprovider",
				Required:    true,
			},
		},
	})
	router.Put(routePath, Update)
}

func Update(c *fiber.Ctx) error {
	log.Debug("received update authprovider request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update authprovider request. authprovider id not found")
		return sdk.AuthProviderBadRequest("Invalid request. Authprovider id is required", c)
	}
	payload := new(sdk.AuthProvider)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update authprovider request", "error", err)
		return sdk.AuthProviderBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}

	payload.Id = id
	pr := providers.GetProviders(c)
	err := pr.S.AuthProviders.Update(c.Context(), payload)
	if err != nil {
		if errors.Is(err, authprovider.ErrAuthProviderNotFound) {
			return sdk.AuthProviderBadRequest("Auth Provider not found", c)
		}
		message := fmt.Errorf("failed to update authprovider. %w", err).Error()
		log.Errorw("failed to update authprovider", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}

	log.Debug("authprovider updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.AuthProviderResponse{
		Success: true,
		Message: "Authprovider updated successfully",
		Data:    payload,
	})
}

func EnableServiceAccountRoute(router fiber.Router, basePath string) {
	routePath := "/enable-service-account"
	path := basePath + routePath
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Enable Service Account",
		Description: "Enable service account authentication for a project",
		Tags:        routeTags,
		RequestBody: &docs.ApiRequestBody{
			Description: "Project ID",
			Content: map[string]interface{}{
				"project_id": "string",
			},
		},
		Response: &docs.ApiResponse{
			Description: "Service account enabled successfully",
			Content:     new(sdk.AuthProviderResponse),
		},
	})
	router.Post(routePath, EnableServiceAccount)
}

func EnableServiceAccount(c *fiber.Ctx) error {
	log.Debug("received enable service account request")
	
	payload := struct {
		ProjectId string `json:"project_id"`
	}{}
	
	if err := c.BodyParser(&payload); err != nil {
		return sdk.AuthProviderBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}
	
	if payload.ProjectId == "" {
		return sdk.AuthProviderBadRequest("project_id is required", c)
	}
	
	pr := providers.GetProviders(c)
	
	// Check if service account auth provider already exists for this project
	existingProviders, err := pr.S.AuthProviders.GetAll(c.Context(), sdk.AuthProviderQueryParams{
		ProjectIds: []string{payload.ProjectId},
	})
	if err != nil {
		return sdk.AuthProviderInternalServerError(fmt.Errorf("failed to check existing providers. %w", err).Error(), c)
	}
	
	for _, provider := range existingProviders {
		if provider.Provider == sdk.AuthProviderTypeGoIAMClient {
			return c.Status(http.StatusOK).JSON(sdk.AuthProviderResponse{
				Success: true,
				Message: "Service account already enabled for this project",
				Data:    &provider,
			})
		}
	}
	
	// Create GOIAM/CLIENT auth provider
	authProvider := &sdk.AuthProvider{
		Name:      "GOIAM/CLIENT",
		Icon:      "key",
		Provider:  sdk.AuthProviderTypeGoIAMClient,
		ProjectId: payload.ProjectId,
		Params:    []sdk.AuthProviderParam{},
		Enabled:   true,
	}
	
	err = pr.S.AuthProviders.Create(c.Context(), authProvider)
	if err != nil {
		message := fmt.Errorf("failed to enable service account. %w", err).Error()
		log.Errorw("failed to enable service account", "error", message)
		return sdk.AuthProviderInternalServerError(message, c)
	}
	
	log.Debug("service account enabled successfully")
	return c.Status(http.StatusOK).JSON(sdk.AuthProviderResponse{
		Success: true,
		Message: "Service account enabled successfully",
		Data:    authProvider,
	})
}