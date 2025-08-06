package policy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/docs"
)

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
		Parameters: []docs.ApiParameter{
			{
				Name:        "skip",
				In:          "query",
				Description: "Number of records to skip for pagination. Default is 0",
				Required:    false,
			},
			{
				Name:        "limit",
				In:          "query",
				Description: "Maximum number of records to return. Default is 10",
				Required:    false,
			},
		},
		Tags: routeTags,
	})
	router.Get(routePath, FetchAll)
}

func FetchAll(c *fiber.Ctx) error {
	log.Debug("received get policies request")
	pr := providers.GetProviders(c)

	query := sdk.PolicyQuery{
		Skip:  0,  // Default value
		Limit: 10, // Default value
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
	ds, err := pr.S.Policy.GetAll(c.Context(), query)
	if err != nil {
		status := http.StatusInternalServerError
		message := fmt.Errorf("failed to get Policy. %w", err).Error()
		log.Error("failed to get Policy", "error", err)
		return c.Status(status).JSON(sdk.PoliciesResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("Policies fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.PoliciesResponse{
		Success: true,
		Message: "Policies fetched successfully",
		Data:    *ds,
	})
}
