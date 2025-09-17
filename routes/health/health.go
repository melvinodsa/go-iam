package health

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/utils/docs"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Data      HealthData `json:"data"`
	Timestamp string     `json:"timestamp"`
}

// HealthData contains the health check information
type HealthData struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Uptime     string            `json:"uptime"`
	Components map[string]string `json:"components"`
}

var startTime = time.Now()

func HealthRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Get(routePath, Health)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Health Check",
		Description: "Get application health status",
		Response: &docs.ApiResponse{
			Description: "Health status retrieved successfully",
			Content:     new(HealthResponse),
		},
		Tags:                 routeTags,
		ProjectIDNotRequired: true,
	})
}

func Health(c *fiber.Ctx) error {
	log.Debug("received health check request")

	pr := providers.GetProviders(c)
	uptime := time.Since(startTime)

	// Check component health
	components := make(map[string]string)

	// Check database connection
	if pr.D != nil {
		// Try to ping the database using a simple operation
		ctx := c.Context()
		_, err := pr.S.Projects.GetAll(ctx)
		if err != nil {
			components["database"] = "unhealthy"
			log.Warnw("database health check failed", "error", err)
		} else {
			components["database"] = "healthy"
		}
	} else {
		components["database"] = "unavailable"
	}

	// Check cache connection
	if pr.C != nil {
		// For now, just check if cache service exists
		// You could add a ping method to cache service for more thorough testing
		components["cache"] = "healthy"
	} else {
		components["cache"] = "unavailable"
	}

	// Determine overall status
	status := "healthy"
	for _, componentStatus := range components {
		if componentStatus == "unhealthy" {
			status = "unhealthy"
			break
		} else if componentStatus == "unavailable" && status == "healthy" {
			status = "degraded"
		}
	}

	healthData := HealthData{
		Status:     status,
		Version:    "1.0.0", // You can make this configurable
		Uptime:     uptime.String(),
		Components: components,
	}

	response := HealthResponse{
		Success:   true,
		Message:   "Health check completed successfully",
		Data:      healthData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Return different HTTP status codes based on health status
	statusCode := http.StatusOK
	switch status {
	case "healthy":
		statusCode = http.StatusOK
	case "unhealthy":
		statusCode = http.StatusServiceUnavailable
	case "degraded":
		statusCode = http.StatusPartialContent
	}

	log.Debug("health check completed successfully")
	return c.Status(statusCode).JSON(response)
}
