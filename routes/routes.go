package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/routes/auth"
	"github.com/melvinodsa/go-iam/routes/authprovider"
	"github.com/melvinodsa/go-iam/routes/client"
	"github.com/melvinodsa/go-iam/routes/health"
	"github.com/melvinodsa/go-iam/routes/me"
	"github.com/melvinodsa/go-iam/routes/policy"
	"github.com/melvinodsa/go-iam/routes/project"
	"github.com/melvinodsa/go-iam/routes/resource"
	"github.com/melvinodsa/go-iam/routes/role"
	"github.com/melvinodsa/go-iam/routes/user"
)

func RegisterRoutes(app *fiber.App, prv *providers.Provider) {
	RegisterOpenRoutes(app, prv)
	RegisterAuthRoutes(app, prv)
}

func RegisterAuthRoutes(app *fiber.App, prv *providers.Provider) {
	ap := app.Use(prv.AM.User)
	project.RegisterRoutes(ap, "/project")
	client.RegisterRoutes(ap, "/client")
	authprovider.RegisterRoutes(ap, "/authprovider")
	auth.RegisterRoutes(ap, "/auth")
	user.RegisterRoutes(ap, "/user")
	resource.RegisterRoutes(ap, "/resource")
	role.RegisterRoutes(ap, "/role")
	policy.RegisterRoutes(ap, "/policy")
	me.RegisterRoutes(app, "/me")
}

func RegisterOpenRoutes(app *fiber.App, prv *providers.Provider) {
	me.RegisterOpenRoutes(app, "/me", prv)
	auth.RegisterRoutes(app, "/auth")
	health.RegisterRoutes(app, "/health")

	app.Static("/docs", "./docs")
}
