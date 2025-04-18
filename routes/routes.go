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
	project.RegisterRoutes(app.Group("/project"))
	client.RegisterRoutes(app.Group("/client"))
	authprovider.RegisterRoutes(app.Group("/authprovider"))
	auth.RegisterRoutes(app.Group("/auth"))
	me.RegisterRoutes(app.Group("/me"))
	user.RegisterRoutes(app.Group("/user"))
	resource.RegisterRoutes(app.Group("/resource"))
	role.RegisterRoutes(app.Group("/role"))
	policy.RegisterRoutes(app.Group("/policy"))
}
