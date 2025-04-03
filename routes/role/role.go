package role

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
)

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

	return c.Status(http.StatusOK).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role created successfully",
		Data:    payload,
	})
}

// Get retrieves a specific role by ID
func Get(c *fiber.Ctx) error {
	log.Debug("received get role request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get role request. role id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. Role id is required",
		})
	}

	pr := providers.GetProviders(c)
	ds, err := pr.S.Role.GetById(c.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get role. %w", err).Error()
		if errors.Is(err, sdk.ErrRoleNotFound) {
			status = http.StatusBadRequest
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

// Update modifies an existing role
func Update(c *fiber.Ctx) error {
	log.Debug("received update role request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update role request. role id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. Role id is required",
		})
	}

	// print request body
	fmt.Println(string(c.Body()))

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
			status = http.StatusBadRequest
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

func AddRoleToUser(c *fiber.Ctx) error {
	userid := c.Params("userid")
	if userid == "" {
		log.Error("invalid add role to user request. user id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. User id is required",
		})
	}

	roleid := c.Params("roleid")
	if roleid == "" {
		log.Error("invalid add role to user request. role id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. Role id is required",
		})
	}

	pr := providers.GetProviders(c)

	err := pr.S.Role.AddRoleToUser(c.Context(), userid, roleid)
	if err != nil {
		log.Error("error adding role to user", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.RoleResponse{
			Success: false,
			Message: fmt.Errorf("error adding role to user. %w", err).Error(),
		})
	}
	log.Debug("role added to user successfully")
	return c.Status(http.StatusOK).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role added to user successfully",
	})
}

func RemoveRoleFromUser(c *fiber.Ctx) error {
	userid := c.Params("userid")
	if userid == "" {
		log.Error("invalid remove role from user request. user id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. User id is required",
		})
	}

	roleid := c.Params("roleid")
	if roleid == "" {
		log.Error("invalid remove role from user request. role id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.RoleResponse{
			Success: false,
			Message: "Invalid request. Role id is required",
		})
	}

	pr := providers.GetProviders(c)

	err := pr.S.Role.RemoveRoleFromUser(c.Context(), userid, roleid)
	if err != nil {
		log.Error("error removing role from user", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.RoleResponse{
			Success: false,
			Message: fmt.Errorf("error removing role from user. %w", err).Error(),
		})
	}
	log.Debug("role removed from user successfully")
	return c.Status(http.StatusOK).JSON(sdk.RoleResponse{
		Success: true,
		Message: "Role removed from user successfully",
	})
}
