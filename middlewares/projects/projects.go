// Package projects provides project-based access control middleware for the Go IAM system.
// It handles extraction of project information from request headers and stores
// it in the request context for authorization decisions.
package projects

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
)

// Middlewares provides project-based access control middleware functionality.
// It encapsulates the project service needed for project-related operations
// and validates project access permissions.
type Middlewares struct {
	projectSvc project.Service // Project service for project operations
}

// NewMiddlewares creates a new project middleware instance.
// It initializes the middleware with the project service required
// for project-related operations and validations.
//
// Parameters:
//   - projectSvc: Project service for project operations
//
// Returns:
//   - *Middlewares: Configured project middleware instance
func NewMiddlewares(projectSvc project.Service) *Middlewares {
	return &Middlewares{
		projectSvc: projectSvc,
	}
}

// Projects is a Fiber middleware that extracts project IDs from request headers.
// It looks for the "X-Project-Ids" header containing a comma-separated list
// of project IDs and stores them in the request context for authorization
// decisions by downstream handlers.
//
// Header format: X-Project-Ids: project1,project2,project3
//
// Usage:
//
//	app.Use(projectMiddleware.Projects)
//
// Parameters:
//   - c: Fiber context containing the HTTP request
//
// Returns:
//   - error: Always returns nil (continues to next middleware)
func (m Middlewares) Projects(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	projectIds := []string{}
	projectIdsCsv, ok := headers["X-Project-Ids"]
	if ok && len(projectIdsCsv) > 0 {
		projectIds = strings.Split(projectIdsCsv[0], ",")
	}
	c.Context().SetUserValue(sdk.ProjectsTypeVal, projectIds)
	return c.Next()
}
