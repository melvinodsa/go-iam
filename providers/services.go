package providers

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/project"
)

type Service struct {
	Projects      project.Service
	Clients       client.Service
	AuthProviders authprovider.Service
}

func NewServices(db db.DB, cache *Cache, enc encrypt.Service) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)
	csvc := client.NewService(cstr)
	apStr := authprovider.NewStore(enc, db)
	apSvc := authprovider.NewService(apStr)
	return &Service{Projects: psvc, Clients: csvc, AuthProviders: apSvc}
}
