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
	"github.com/melvinodsa/go-iam/services/policy/system"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/services/resource"
	"github.com/melvinodsa/go-iam/services/role"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

// Service is a container that holds all business logic services for the Go IAM system.
// It provides centralized access to all domain services and manages their dependencies.
// Services are organized by domain and provide the core functionality for IAM operations.
type Service struct {
	Projects      project.Service      // Project management service
	Clients       client.Service       // OAuth2/OIDC client management service
	AuthProviders authprovider.Service // Authentication provider management service
	Auth          auth.Service         // Authentication and token validation service
	Resources     resource.Service     // Resource management service
	User          user.Service         // User management and authorization service
	Role          role.Service         // Role-based access control service
	Policy        policy.Service       // Policy management service
}

// NewServices creates and configures all business logic services with their dependencies.
// This function initializes services in the correct order to satisfy dependencies and
// sets up event subscriptions for cross-service communication.
//
// Service initialization order:
// 1. Core services (project, resource, role, user)
// 2. Authentication services (auth provider, client, auth)
// 3. Policy services
// 4. Event subscriptions for reactive updates
//
// Parameters:
//   - db: Database connection interface
//   - cache: Cache service for performance optimization
//   - enc: Encryption service for sensitive data
//   - jwtSvc: JWT service for token operations
//   - tokenTTL: Token time-to-live in minutes
//   - refetchTTL: Auth provider refetch interval in minutes
//
// Returns:
//   - *Service: Configured service container with all dependencies wired
func NewServices(db db.DB, cache cache.Service, enc encrypt.Service, jwtSvc jwt.Service, tokenTTL int64, refetchTTL int64) *Service {
	pstr := project.NewStore(db)
	psvc := project.NewService(pstr)
	cstr := client.NewStore(db)

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
	csvc := client.NewService(cstr, psvc, apSvc, userSvc)
	authSvc := auth.NewService(apSvc, csvc, cache, jwtSvc, enc, userSvc, tokenTTL, refetchTTL)
	polstr := policy.NewStore()
	polSvc := policy.NewService(polstr)

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
