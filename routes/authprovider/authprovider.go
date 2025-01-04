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
)

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
