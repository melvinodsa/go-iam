package policy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/policy"
)

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
		if errors.Is(err, policy.ErrPolicyNotFound) {
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
		if errors.Is(err, policy.ErrPolicyNotFound) {
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
		if errors.Is(err, policy.ErrPolicyNotFound) {
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
