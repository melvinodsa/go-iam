package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/middlewares/auth"
	"github.com/melvinodsa/go-iam/middlewares/projects"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	goaiamclient "github.com/melvinodsa/go-iam/utils/goiamclient"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

func InjectTestProviders(cnf config.AppConfig, d db.DB) (*providers.Provider, error) {

	var cS cache.Service = cache.NewMockService()

	enc, err := encrypt.NewService(cnf.Encrypter.Key())
	if err != nil {
		return nil, fmt.Errorf("error creating encrypter: %w", err)
	}

	jwtSvc := jwt.NewService(cnf.Jwt.Secret())

	svcs := providers.NewServices(d, cS, enc, jwtSvc, cnf.Server.TokenCacheTTLInMinutes, cnf.Server.AuthProviderRefetchIntervalInMinutes)
	pm := projects.NewMiddlewares(svcs.Projects)
	am, err := auth.NewMiddlewares(svcs.Auth, svcs.Clients)
	if err != nil {
		return nil, err
	}
	authClient, err := goaiamclient.GetGoIamClient(svcs.Clients)
	if err != nil {
		return nil, err
	}

	pvd := &providers.Provider{
		S:          svcs,
		D:          d,
		C:          cS,
		PM:         pm,
		AM:         am,
		AuthClient: authClient,
	}

	// subscribe to client events for checking auth client
	svcs.Clients.Subscribe(goiamuniverse.EventClientCreated, pvd)
	svcs.Clients.Subscribe(goiamuniverse.EventClientUpdated, pvd)
	svcs.Clients.Subscribe(goiamuniverse.EventClientCreated, svcs.Auth)
	svcs.Clients.Subscribe(goiamuniverse.EventClientUpdated, svcs.Auth)

	// creating default project if it doesn't exist
	err = providers.CheckAndAddDefaultProject(svcs.Projects)
	if err != nil {
		log.Errorw("error checking and adding default project", "error", err)
		return nil, fmt.Errorf("error checking and adding default project: %w", err)
	}

	return pvd, nil
}

func SetupTestServer(app *fiber.App, db db.DB) *providers.Provider {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)
	prv, err := InjectTestProviders(*cnf, db)
	if err != nil {
		log.Fatalf("error injecting providers %s", err)
	}
	app.Use((*cnf).Handle)
	app.Use(providers.Handle(prv))
	app.Use(cors.New())

	app.Use(prv.PM.Projects)

	return prv
}
