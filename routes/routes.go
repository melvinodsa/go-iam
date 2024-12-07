package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/api-server/routes/project"
)

func RegisterRoutes(app *fiber.App) {
	project.RegisterRoutes(app.Group("/project"))
}
