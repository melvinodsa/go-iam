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
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/utils/goiamclient"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/mock"
)

func GetServices(cnf config.AppConfig, cS cache.Service, d db.DB) (*providers.Service, error) {

	enc, err := encrypt.NewService(cnf.Encrypter.Key())
	if err != nil {
		return nil, fmt.Errorf("error creating encrypter: %w", err)
	}

	jwtSvc := jwt.NewService(cnf.Jwt.Secret())

	svcs := providers.NewServices(d, cS, enc, jwtSvc, cnf.Server.TokenCacheTTLInMinutes, cnf.Server.AuthProviderRefetchIntervalInMinutes)

	mockClientSvc := services.MockClientService{}
	mockProjectSvc := services.MockProjectService{}
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
	mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()
	mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{}, nil).Once()

	svcs.Clients = &mockClientSvc
	svcs.Projects = &mockProjectSvc

	return svcs, nil
}

func InjectTestProviders(svcs *providers.Service, cS cache.Service, d db.DB) (*providers.Provider, error) {

	pm := projects.NewMiddlewares(svcs.Projects)
	am, err := auth.NewMiddlewares(svcs.Auth, svcs.Clients)
	if err != nil {
		return nil, err
	}
	authClient, err := goiamclient.GetGoIamClient(svcs.Clients)
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

func SetupTestServer(app *fiber.App, cnf *config.AppConfig, svcs *providers.Service, cS cache.Service, db db.DB) *providers.Provider {
	prv, err := InjectTestProviders(svcs, cS, db)
	if err != nil {
		log.Fatalf("error injecting providers %s", err)
	}
	app.Use((*cnf).Handle)
	app.Use(providers.Handle(prv))
	app.Use(cors.New())

	app.Use(prv.PM.Projects)

	return prv
}
