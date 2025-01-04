package providers

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/auth"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/services/project"
)

type Service struct {
	Projects      project.Service
	Clients       client.Service
	AuthProviders authprovider.Service
	Auth          auth.Service
}

func NewServices(db db.DB, cache *cache.Service, enc encrypt.Service, jwtSvc jwt.Service) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)
	csvc := client.NewService(cstr, psvc)
	apStr := authprovider.NewStore(enc, db)
	apSvc := authprovider.NewService(apStr, psvc)
	authSvc := auth.NewService(apSvc, csvc, *cache, jwtSvc, enc)
	return &Service{Projects: psvc, Clients: csvc, AuthProviders: apSvc, Auth: authSvc}
}
