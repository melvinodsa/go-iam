package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/routes"
)

func SetupServer(app *fiber.App) *config.AppConfig {
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)
	prv, err := providers.InjectDefaultProviders(*cnf)
	if err != nil {
		log.Fatalf("error injecting providers %s", err)
	}
	app.Use((*cnf).Handle)
	app.Use(providers.Handle(prv))
	app.Use(cors.New())

	app.Use(prv.PM.Projects)
	routes.RegisterRoutes(app, prv)

	return cnf
}
