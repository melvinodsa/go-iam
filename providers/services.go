package providers

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/auth"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/services/policy"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/services/resource"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/services/user"
)

type Service struct {
	Projects      project.Service
	Clients       client.Service
	AuthProviders authprovider.Service
	Auth          auth.Service
	Resources     resource.Service
	User          user.Service
	Role          role.Service
	Policy        policy.Service
}

func NewServices(db db.DB, cache cache.Service, enc encrypt.Service, jwtSvc jwt.Service) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)
	csvc := client.NewService(cstr, psvc)
	userStr := user.NewStore(db)
	userSvc := user.NewService(userStr)
	rstr := resource.NewStore(db)
	rsvc := resource.NewService(rstr)

	apStr := authprovider.NewStore(enc, db)
	apSvc := authprovider.NewService(apStr, psvc)
	ustr := user.NewStore(db)
	usvc := user.NewService(ustr)
	authSvc := auth.NewService(apSvc, csvc, cache, jwtSvc, enc, usvc)
	polstr := policy.NewStore(db, rstr)
	polSvc := policy.NewService(polstr)
	roleStr := role.NewStore(db)
	roleSvc := role.NewService(roleStr, polSvc)

	return &Service{
		Projects:      psvc,
		Clients:       csvc,
		AuthProviders: apSvc,
		Auth:          authSvc,
		User:          userSvc,
		Resources:     rsvc,
		Role:          roleSvc,
		Policy:        polSvc,
	}
}
