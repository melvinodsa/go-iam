package providers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/encrypt"
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

	enc, err := encrypt.NewService(cnf.Encrypter.Key())
	if err != nil {
		return nil, fmt.Errorf("error creating encrypter: %w", err)
	}

	svcs := NewServices(d, c, enc)

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
