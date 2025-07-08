package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

func GetProjects(ctx context.Context) []string {
	return ctx.Value("projects").([]string)
}
