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
)

func Create(c *fiber.Ctx) error {
	log.Debug("received create client request")
	payload := new(sdk.Client)
	if err := c.BodyParser(payload); err != nil {
		return sdk.ClientBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}
	log.Debug("parsed create client request")
	pr := providers.GetProviders(c)
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

// Get client
func Get(c *fiber.Ctx) error {
	log.Debug("received get client request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get client request. client id not found")
		return sdk.ClientBadRequest("Invalid request. Client id is required", c)
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.Clients.Get(c.Context(), id)
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

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get clients request")
	pr := providers.GetProviders(c)
	ds, err := pr.S.Clients.GetAll(c.Context())
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

// Update client
func Update(c *fiber.Ctx) error {
	log.Debug("received update client request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update client request. client id not found")
		return sdk.ClientBadRequest("Invalid request. Client id is required", c)
	}
	payload := new(sdk.Client)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update client request", "error", err)
		return sdk.ClientBadRequest(fmt.Errorf("invalid request. %w", err).Error(), c)
	}

	payload.Id = id
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
