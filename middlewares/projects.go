package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/sdk"
)

func (m Middlewares) Projects(c *fiber.Ctx) error {
	p, err := m.projectSvc.GetAll(c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Context().SetUserValue("projects", p)
	return c.Next()
}

func GetProjects(ctx context.Context) []sdk.Project {
	return ctx.Value("projects").([]sdk.Project)
}
