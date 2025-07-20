package providers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/middlewares/auth"
	"github.com/melvinodsa/go-iam/middlewares/projects"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/utils"
	goaiamclient "github.com/melvinodsa/go-iam/utils/goiamclient"
)

type Provider struct {
	S          *Service
	D          db.DB
	C          cache.Service
	PM         *projects.Middlewares
	AM         *auth.Middlewares
	AuthClient *sdk.Client
}

func InjectDefaultProviders(cnf config.AppConfig) (*Provider, error) {
	d, err := NewDBConnection(cnf)
	if err != nil {
		return nil, err
	}

	var cS cache.Service = cache.NewMockService()

	if cnf.Server.EnableRedis {
		cS = cache.NewRedisService(cnf.Redis.Host, cnf.Redis.Password)
	}

	enc, err := encrypt.NewService(cnf.Encrypter.Key())
	if err != nil {
		return nil, fmt.Errorf("error creating encrypter: %w", err)
	}

	jwtSvc := jwt.NewService(cnf.Jwt.Secret())

	svcs := NewServices(d, cS, enc, jwtSvc, cnf.Server.TokenCacheTTLInMinutes, cnf.Server.AuthProviderRefetchIntervalInMinutes)
	pm := projects.NewMiddlewares(svcs.Projects)
	am := auth.NewMiddlewares(svcs.Auth, svcs.Clients)

	pvd := &Provider{
		S:          svcs,
		D:          d,
		C:          cS,
		PM:         pm,
		AM:         am,
		AuthClient: goaiamclient.GetGoIamClient(svcs.Clients),
	}

	// subscribe to client events for checking auth client
	svcs.Clients.Subscribe(sdk.EventClientCreated, pvd)
	svcs.Clients.Subscribe(sdk.EventClientUpdated, pvd)

	// creating default project if it doesn't exist
	err = checkAndAddDefaultProject(svcs.Projects)
	if err != nil {
		log.Errorw("error checking and adding default project", "error", err)
		return nil, fmt.Errorf("error checking and adding default project: %w", err)
	}

	return pvd, nil
}

type keyType struct {
	key string
}

var providerKey = keyType{"providers"}

func Handle(p *Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(providerKey, p)
		return c.Next()
	}
}

func GetProviders(c *fiber.Ctx) *Provider {
	return c.Locals(providerKey).(*Provider)
}

func (p *Provider) HandleEvent(e utils.Event[sdk.Client]) {
	if e.Name() != sdk.EventClientCreated && e.Name() != sdk.EventClientUpdated {
		return
	}
	if !e.Payload().GoIamClient {
		return
	}
	p.AuthClient = goaiamclient.GetGoIamClient(p.S.Clients)
	p.AM.AuthClient = p.AuthClient
}
