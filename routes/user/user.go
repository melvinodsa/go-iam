package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
)

// Create user
func Create(c *fiber.Ctx) error {
	log.Debug("received create user request")
	payload := new(sdk.User)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed create user request")
	pr := providers.GetProviders(c)
	err := pr.S.User.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to create user. %v", err)
		log.Errorw("failed to create user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("user created successfully")

	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    payload,
	})
}

// Get user by ID
func GetById(c *fiber.Ctx) error {
	log.Debug("received get user request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid get user request. user id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: "Invalid request. User ID is required",
		})
	}
	pr := providers.GetProviders(c)
	ds, err := pr.S.User.GetById(c.Context(), id)
	if err != nil {
		if errors.Is(err, sdk.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
				Success: false,
				Message: "User not found",
			})
		}
		message := fmt.Sprintf("failed to get user. %v", err)
		log.Error("failed to get user", "error", message)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    ds,
	})
}

// Get all users
func GetAll(c *fiber.Ctx) error {
	log.Debug("received get users request")
	query := sdk.UserQuery{
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
	users, err := pr.S.User.GetAll(c.Context(), query)
	if err != nil {
		message := fmt.Sprintf("failed to get users. %v", err)
		log.Error("failed to get users", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("users fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserListResponse{
		Success: true,
		Message: "Users fetched successfully",
		Data:    users,
	})
}

// Update user
func Update(c *fiber.Ctx) error {
	log.Debug("received update user request")
	id := c.Params("id")
	if id == "" {
		log.Error("invalid update user request. user id not found")
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: "Invalid request. User ID is required",
		})
	}
	payload := new(sdk.User)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update user request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}

	payload.Id = id
	pr := providers.GetProviders(c)
	err := pr.S.User.Update(c.Context(), payload)
	if err != nil {
		if errors.Is(err, sdk.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
				Success: false,
				Message: "User not found",
			})
		}
		message := fmt.Sprintf("failed to update user. %v", err)
		log.Error("failed to update user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    payload,
	})
}
