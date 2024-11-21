package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/api-server/routes/projects"
)

func RegisterRoutes(app *fiber.App) {
	projects.RegisterRoutes(app.Group("/projects"))
}
