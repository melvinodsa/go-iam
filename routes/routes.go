package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/routes/auth"
	"github.com/melvinodsa/go-iam/routes/authprovider"
	"github.com/melvinodsa/go-iam/routes/client"
	"github.com/melvinodsa/go-iam/routes/me"
	"github.com/melvinodsa/go-iam/routes/project"
	"github.com/melvinodsa/go-iam/routes/resource"
)

func RegisterRoutes(app *fiber.App) {
	project.RegisterRoutes(app.Group("/project"))
	client.RegisterRoutes(app.Group("/client"))
	authprovider.RegisterRoutes(app.Group("/authprovider"))
	auth.RegisterRoutes(app.Group("/auth"))
	me.RegisterRoutes(app.Group("/me"))
	resource.RegisterRoutes(app.Group("/resource"))
}
