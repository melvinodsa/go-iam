package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// CreateRoute registers the routes for user management
func CreateRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Post(routePath, Create)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPost,
		Name:        "Create User",
		Description: "Create a new user",
		RequestBody: &docs.ApiRequestBody{
			Description: "User data",
			Content:     new(sdk.User),
		},
		Response: &docs.ApiResponse{
			Description: "User created successfully",
			Content:     new(sdk.UserResponse),
		},
		Tags: routeTags,
	})
}

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

	return c.Status(http.StatusCreated).JSON(sdk.UserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    payload,
	})
}

// GetByIdRoute registers the route to get a user by ID
func GetByIdRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Get(routePath, GetById)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get User",
		Description: "Get a user by ID",
		Response: &docs.ApiResponse{
			Description: "User fetched successfully",
			Content:     new(sdk.UserResponse),
		},
		// Parameters for the user ID in the path
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the user",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

// Get user by ID
func GetById(c *fiber.Ctx) error {
	log.Debug("received get user request")
	id := c.Params("id")

	pr := providers.GetProviders(c)
	ds, err := pr.S.User.GetById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), sdk.ErrUserNotFound.Error()) {
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

// GetAllRoute registers the route to get all users
func GetAllRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Get(routePath, GetAll)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get All Users",
		Description: "Get all users",
		Response: &docs.ApiResponse{
			Description: "Users fetched successfully",
			Content:     new(sdk.UserListResponse),
		},
		// Parameters for pagination
		Parameters: []docs.ApiParameter{
			{
				Name:        "query",
				In:          "query",
				Description: "Search query for filtering users",
				Required:    false,
			},
			{
				Name:        "skip",
				In:          "query",
				Description: "Number of users to skip for pagination. Default is 0",
				Required:    false,
			},
			{
				Name:        "limit",
				In:          "query",
				Description: "Maximum number of users to return. Default is 10",
				Required:    false,
			},
		},
		Tags: routeTags,
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

// UpdateRoute registers the route for updating a user
func UpdateRoute(router fiber.Router, basePath string) {
	routePath := "/:id"
	path := basePath + routePath
	router.Put(routePath, Update)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update User",
		Description: "Update a user by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "User data",
			Content:     new(sdk.User),
		},
		Response: &docs.ApiResponse{
			Description: "User updated successfully",
			Content:     new(sdk.UserResponse),
		},
		// Parameters for the user ID in the path
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the user",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

// Update user
func Update(c *fiber.Ctx) error {
	log.Debug("received update user request")
	id := c.Params("id")

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

// UpdateRolesRoute registers the route for updating user roles
func UpdateRolesRoute(router fiber.Router, basePath string) {
	routePath := "/:id/roles"
	path := basePath + routePath
	router.Put(routePath, UpdateRoles)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update User Roles",
		Description: "Update roles for a user by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "User roles update data",
			Content:     new(sdk.UserRoleUpdate),
		},
		Response: &docs.ApiResponse{
			Description: "User roles updated successfully",
			Content:     new(sdk.UserResponse),
		},
		// Parameters for the user ID in the path
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the user",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func UpdateRoles(c *fiber.Ctx) error {
	log.Debug("received update user roles request")
	id := c.Params("id")

	payload := new(sdk.UserRoleUpdate)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update user roles request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}

	pr := providers.GetProviders(c)
	for _, roleId := range payload.ToBeRemoved {
		if err := pr.S.User.RemoveRoleFromUser(c.Context(), id, roleId); err != nil {
			if errors.Is(err, sdk.ErrRoleNotFound) {
				return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
					Success: false,
					Message: fmt.Sprintf("Role %s not found", roleId),
				})
			}
			message := fmt.Sprintf("failed to remove role %s from user %s. %v", roleId, id, err)
			log.Errorw("failed to remove role from user", "error", message)
			return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
				Success: false,
				Message: message,
			})
		}
	}

	for _, roleId := range payload.ToBeAdded {
		if err := pr.S.User.AddRoleToUser(c.Context(), id, roleId); err != nil {
			if errors.Is(err, sdk.ErrRoleNotFound) {
				return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
					Success: false,
					Message: fmt.Sprintf("Role %s not found", roleId),
				})
			}
			message := fmt.Sprintf("failed to add role %s to user %s. %v", roleId, id, err)
			log.Errorw("failed to add role to user", "error", message)
			return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
				Success: false,
				Message: message,
			})
		}
	}

	log.Debug("user roles updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User roles updated successfully",
	})
}

// UpdatePoliciesRoute registers the route for updating user policies
func UpdatePoliciesRoute(router fiber.Router, basePath string) {
	routePath := "/:id/policies"
	path := basePath + routePath
	router.Put(routePath, UpdatePolicies)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Update User Policies",
		Description: "Update policies for a user by ID",
		RequestBody: &docs.ApiRequestBody{
			Description: "User policies update data",
			Content:     new(sdk.UserPolicyUpdate),
		},
		Response: &docs.ApiResponse{
			Description: "User policies updated successfully",
			Content:     new(sdk.UserResponse),
		},
		// Parameters for the user ID in the path
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the user",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func UpdatePolicies(c *fiber.Ctx) error {
	log.Debug("received update user policies request")
	id := c.Params("id")

	payload := new(sdk.UserPolicyUpdate)
	if err := c.BodyParser(payload); err != nil {
		log.Errorw("invalid update user roles request", "error", err)
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}

	pr := providers.GetProviders(c)
	err := pr.S.User.RemovePolicyFromUser(c.Context(), id, payload.ToBeRemoved)
	if err != nil {
		if errors.Is(err, sdk.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
				Success: false,
				Message: fmt.Sprintf("User %s not found", id),
			})
		}
		message := fmt.Sprintf("failed to remove policies from user %s. %v", id, err)
		log.Errorw("failed to remove policy from user", "error", message)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	err = pr.S.User.AddPolicyToUser(c.Context(), id, payload.ToBeAdded)
	if err != nil {
		if errors.Is(err, sdk.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
				Success: false,
				Message: fmt.Sprintf("User %s not found", id),
			})
		}

		message := fmt.Sprintf("failed to add policies to user %s. %v", id, err)
		log.Errorw("failed to add policy to user", "error", message)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user policies updated successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User policies updated successfully",
	})
}

// TransferOwnershipRoute registers the route for transferring user ownership
func TransferOwnershipRoute(router fiber.Router, basePath string) {
	routePath := "/:id/transfer-ownership/:oldId"
	path := basePath + routePath
	router.Put(routePath, TransferOwnership)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodPut,
		Name:        "Transfer User Ownership",
		Description: "Transfer ownership of a user",
		Response: &docs.ApiResponse{
			Description: "User ownership transferred successfully",
			Content:     new(sdk.UserResponse),
		},
		// Parameters for the user ID and old ID in the path
		Parameters: []docs.ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "The ID of the user",
				Required:    true,
			},
			{
				Name:        "oldId",
				In:          "path",
				Description: "The ID of the old user",
				Required:    true,
			},
		},
		Tags: routeTags,
	})
}

func TransferOwnership(c *fiber.Ctx) error {
	log.Debug("received transfer user ownership request")
	id := c.Params("id")
	oldId := c.Params("oldId")

	pr := providers.GetProviders(c)
	err := pr.S.User.TransferOwnership(c.Context(), id, oldId)
	if err != nil {
		if errors.Is(err, sdk.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(sdk.UserResponse{
				Success: false,
				Message: fmt.Sprintf("User %s not found", id),
			})
		}

		message := fmt.Sprintf("failed to transfer ownership from user %s to user %s. %v", oldId, id, err)
		log.Errorw("failed to transfer ownership", "error", message)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user ownership transferred successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User ownership transferred successfully",
	})
}
