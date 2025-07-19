package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/routes/auth"
	"github.com/melvinodsa/go-iam/routes/authprovider"
	"github.com/melvinodsa/go-iam/routes/client"
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
	project.RegisterRoutes(ap.Group("/project"))
	client.RegisterRoutes(ap.Group("/client"))
	authprovider.RegisterRoutes(ap.Group("/authprovider"))
	auth.RegisterRoutes(ap.Group("/auth"))
	user.RegisterRoutes(ap.Group("/user"))
	resource.RegisterRoutes(ap.Group("/resource"))
	role.RegisterRoutes(ap.Group("/role"))
	policy.RegisterRoutes(ap.Group("/policy"))
	me.RegisterRoutes(app.Group("/me"))
}

func RegisterOpenRoutes(app *fiber.App, prv *providers.Provider) {
	me.RegisterOpenRoutes(app.Group("/me"), prv)
	auth.RegisterRoutes(app.Group("/auth"))
}
