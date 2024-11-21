package providers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/api-server/config"
	"github.com/melvinodsa/go-iam/api-server/db"
)

type Provider struct {
	S *Service
	D db.DB
	C Cache
}

func InjectDefaultProviders(cnf config.AppConfig) (*Provider, error) {
	d, err := NewDBConnection(cnf)
	if err != nil {
		return nil, err
	}
	c := NewCache(cnf)

	svcs := NewServices(d, c)

	return &Provider{
		S: svcs,
		D: d,
		C: *c,
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
