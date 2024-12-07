package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/routes"
)

func main() {
	app := fiber.New()

	cnf := setupServer(app)

	routes.RegisterRoutes(app)

	for _, route := range app.GetRoutes() {
		if route.Method == "OPTIONS" || route.Method == "HEAD" || route.Method == "TRACE" || route.Method == "CONNECT" {
			continue
		}
		if route.Path == "/" {
			continue
		}
		log.Infof("%s %s", route.Method, route.Path)
	}

	err := app.Listen(":" + cnf.Server.Port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func setupServer(app *fiber.App) *config.AppConfig {
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
	app.Use((*prv).Handle)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
	}))

	return cnf
}
