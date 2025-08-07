package projects

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/services/project"
)

type Middlewares struct {
	projectSvc project.Service
}

func NewMiddlewares(projectSvc project.Service) *Middlewares {
	return &Middlewares{
		projectSvc: projectSvc,
	}
}

func (m Middlewares) Projects(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	projectIds := []string{}
	projectIdsCsv, ok := headers["X-Project-Ids"]
	if ok && len(projectIdsCsv) > 0 {
		projectIds = strings.Split(projectIdsCsv[0], ",")
	}
	c.Context().SetUserValue("projects", projectIds)
	return c.Next()
}
