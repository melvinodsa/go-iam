package providers

import (
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/services/auth"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/services/policy/system"
	"github.com/melvinodsa/go-iam/services/policybeta"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/services/resource"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type Service struct {
	Projects      project.Service
	Clients       client.Service
	AuthProviders authprovider.Service
	Auth          auth.Service
	Resources     resource.Service
	User          user.Service
	Role          role.Service
	Policy        policybeta.Service
}

func NewServices(db db.DB, cache cache.Service, enc encrypt.Service, jwtSvc jwt.Service, tokenTTL int64, refetchTTL int64) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)
	csvc := client.NewService(cstr, psvc)
	rstr := resource.NewStore(db)
	rsvc := resource.NewService(rstr)
	roleStr := role.NewStore(db)
	roleSvc := role.NewService(roleStr)
	userStr := user.NewStore(db)
	userSvc := user.NewService(userStr, roleSvc)

	// subscribing to role updates
	roleSvc.Subscribe(goiamuniverse.EventRoleUpdated, userSvc)
	// subscribing to resource create updates
	rsvc.Subscribe(goiamuniverse.EventResourceCreated, system.NewAccessToCreatedResource(userSvc))
	rsvc.Subscribe(goiamuniverse.EventResourceCreated, system.NewAddResourcesToUser(userSvc))
	rsvc.Subscribe(goiamuniverse.EventResourceCreated, system.NewAddResourcesToRole(userSvc, roleSvc))
	// adding default policies to a user when gets created
	userSvc.Subscribe(goiamuniverse.EventUserCreated, system.NewDefaultPoliciesOnUser(userSvc))

	apStr := authprovider.NewStore(enc, db)
	apSvc := authprovider.NewService(apStr, psvc)
	authSvc := auth.NewService(apSvc, csvc, cache, jwtSvc, enc, userSvc, tokenTTL, refetchTTL)
	polstr := policybeta.NewStore(db, rstr)
	polSvc := policybeta.NewService(polstr)

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
