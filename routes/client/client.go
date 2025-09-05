package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// CreateRoute registers the routes for the client
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Post(routePath, Create)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create Client",
		Description: "Create a new client",
		RequestBody: &docs.ApiRequestBody{
			Description: "Client data",
			Content:     new(sdk.Client),
		},
		Response: &docs.ApiResponse{
			Description: "Client created successfully",
			Content:     new(sdk.ClientResponse),
		},
		Tags: routeTags,
	})
}

func Create(c *fiber.Ctx) error {
	log.Debug("received create client request")
	payload := new(sdk.Client)
	if err := c.BodyParser(payload); err != nil {
		return sdk.ClientBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}
	log.Debug("parsed create client request")
	pr := providers.GetProviders(c)

	if payload.ServiceAccountEmail == "" && payload.DefaultAuthProviderId == "" {
		message := "either service account email or default auth provider id must be provided"
		return sdk.ClientBadRequest(message, c)
	}

	if payload.ServiceAccountEmail == "" && payload.DefaultAuthProviderId == "" {
		message := "either service account email or default auth provider id must be provided"
		return sdk.ClientBadRequest(message, c)
	}
	err := pr.S.Clients.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Errorf("failed to create client. %w", err).Error()
		log.Errorw("failed to create client", "error", err)
		return sdk.ClientInternalServerError(message, c)
	}
	log.Debug("client created successfully")

	return c.Status(http.StatusOK).JSON(sdk.ClientResponse{
		Success: true,
		Message: "Client created successfully",
		Data:    payload,
	})
}

// GetRoute registers the route for getting a client
func GetRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Get(routePath, Get)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Client",
		Description: "Get a client by ID",
		Response: &docs.ApiResponse{
			Description: "Client fetched successfully",
			Content:     new(sdk.ClientResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the client",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

// Get client
func Get(c *fiber.Ctx) error {
	log.Debug("received get client request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get client request. client id not found")
		return sdk.ClientBadRequest("Invalid request. Client id is required", c)
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.Clients.Get(c.Context(), id, false)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			return sdk.ClientBadRequest("Client not found", c)
		}
		message := fmt.Errorf("failed to get client. %w", err).Error()
		log.Error("failed to get client", "error", message)
		return sdk.ClientInternalServerError(message, c)
	}

	log.Debug("client fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.ClientResponse{
		Success: true,
		Message: "Client fetched successfully",
		Data:    ds,
	})
}

// FetchAllRoute registers the route for fetching all clients
func FetchAllRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Get(routePath, FetchAll)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Fetch All Clients",
		Description: "Fetch all clients",
		Response: &docs.ApiResponse{
			Description: "Clients fetched successfully",
			Content:     new(sdk.ClientsResponse),
		},
		Tags: routeTags,
	})
}

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get clients request")
	pr := providers.GetProviders(c)
	ds, err := pr.S.Clients.GetAll(c.Context(), sdk.ClientQueryParams{})
	if err != nil {
		message := fmt.Errorf("failed to get clients. %w", err).Error()
		log.Error("failed to get clients", "error", err)
		return sdk.ClientsInternalServerError(message, c)
	}

	log.Debug("clients fetched successfully")

	return c.Status(http.StatusOK).JSON(sdk.ClientsResponse{
		Success: true,
		Message: "Clients fetched successfully",
		Data:    ds,
	})
}

// UpdateRoute registers the route for updating a client
func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Put(routePath, Update)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update Client",
		Description: "Update a client by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "Client data",
			Content:     new(sdk.Client),
		},
		Response: &docs.ApiResponse{
			Description: "Client updated successfully",
			Content:     new(sdk.ClientResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the client",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func Update(c *fiber.Ctx) error {
	log.Debug("received update client request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update client request. client id not found")
		return sdk.ClientBadRequest("Invalid request. Client id is required", c)
	}
	// Use ClientUpdateRequest to handle email
	payload := new(sdk.Client)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update client request", "error", err)
		return sdk.ClientBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}

	payload.Id = id
	if payload.ServiceAccountEmail == "" && payload.DefaultAuthProviderId == "" {
		message := "either service account email or default auth provider id must be provided"
		return sdk.ClientBadRequest(message, c)
	}
	pr := providers.GetProviders(c)

	err := pr.S.Clients.Update(c.Context(), payload)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			return sdk.ClientBadRequest("Client not found", c)
		}
		message := fmt.Errorf("failed to update client. %w", err).Error()
		log.Error("failed to update client", "error", err)
		return sdk.ClientInternalServerError(message, c)
	}

	log.Debug("client updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.ClientResponse{
		Success: true,
		Message: "Client updated successfully",
		Data:    payload,
	})
}

func RegenerateSecretRoute(router fiber.Router, basePath string) {
	routePath := "/:id/regenerate-secret"
	path := basePath + routePath
	router.Put(routePath, RegenerateSecret)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Regenerate Client Secret",
		Description: "Regenerate the secret for a client",
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the client",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func RegenerateSecret(c *fiber.Ctx) error {
	log.Debug("received regenerate client secret request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid regenerate client secret request. client id not found")
		return sdk.ClientBadRequest("Invalid request. Client id is required", c)
	}
	pr := providers.GetProviders(c)

	client, err := pr.S.Clients.RegenerateSecret(c.Context(), id)
	if err != nil {
		message := fmt.Errorf("failed to regenerate client secret. %w", err).Error()
		log.Error("failed to regenerate client secret", "error", err)
		return sdk.ClientInternalServerError(message, c)
	}

	log.Debug("client secret regenerated successfully")
	return c.Status(http.StatusOK).JSON(sdk.ClientResponse{
		Success: true,
		Message: "Client secret regenerated successfully",
		Data:    client,
	})
}
