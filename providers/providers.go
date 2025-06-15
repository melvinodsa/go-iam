package providers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
)

type Provider struct {
	S *Service
	D db.DB
	C cache.Service
	M middlewares.Middlewares
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
	mid := middlewares.NewMiddlewares(svcs.Projects, d)

	return &Provider{
		S: svcs,
		D: d,
		C: cS,
		M: *mid,
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
