package providers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
)

type Provider struct {
	S *Service
	D db.DB
	C cache.Service
	M *middlewares.Middlewares
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

	svcs := NewServices(d, cS, enc, jwtSvc)
	authEnabled := checkIfGoIamClientEnabled(svcs.Clients)
	mid := middlewares.NewMiddlewares(svcs.Projects, d, authEnabled)

	svcs.Clients.Subscribe(sdk.EventClientCreated, mid)
	svcs.Clients.Subscribe(sdk.EventClientUpdated, mid)

	return &Provider{
		S: svcs,
		D: d,
		C: cS,
		M: mid,
	}, nil
}

type keyType struct {
	key string
}

var providerKey = keyType{"providers"}

func (p Provider) Handle(c *fiber.Ctx) error {
	c.Locals(providerKey, p)
	return c.Next()
}

func GetProviders(c *fiber.Ctx) Provider {
	return c.Locals(providerKey).(Provider)
}

func checkIfGoIamClientEnabled(svc client.Service) bool {
	prvs, err := svc.GetGoIamClients(context.Background(), sdk.ClientQueryParams{
		GoIamClient: true,
	})
	if err != nil {
		log.Errorw("error getting go iam client", "error", err)
		return false
	}
	log.Warn("IAM running in insecure mode. Create a client for Go IAM to make the application secure")
	return len(prvs) > 0
}
