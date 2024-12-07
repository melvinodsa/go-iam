package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/routes/client"
	"github.com/melvinodsa/go-iam/routes/project"
)

func RegisterRoutes(app *fiber.App) {
	project.RegisterRoutes(app.Group("/project"))
	client.RegisterRoutes(app.Group("/client"))
}
