// Package providers manages dependency injection and initialization for the Go IAM system.
// It orchestrates the setup of services, database connections, caching, middleware,
// and event subscriptions. This package serves as the central dependency injection
// container that wires together all components of the IAM system.
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
	"github.com/melvinodsa/go-iam/services/policy/system"
	"github.com/melvinodsa/go-iam/utils"
	goiamclient "github.com/melvinodsa/go-iam/utils/goiamclient"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

// Provider encapsulates all the core dependencies and services for the Go IAM system.
// It serves as the main dependency injection container that holds services, database
// connections, cache, middlewares, and client configurations.
type Provider struct {
	S          *Service              // Service container with all business logic services
	D          db.DB                 // Database connection interface
	C          cache.Service         // Cache service (Redis or mock)
	PM         *projects.Middlewares // Project-based access control middleware
	AM         *auth.Middlewares     // Authentication middleware
	AuthClient *sdk.Client           // Go IAM system client configuration
}

// InjectDefaultProviders creates and configures a complete Provider instance with all dependencies.
// This is the main entry point for initializing the Go IAM system. It sets up database connections,
// caching, services, middlewares, and event subscriptions based on the provided configuration.
//
// The function performs the following initialization steps:
// 1. Establishes database connection and runs migrations
// 2. Configures caching (Redis or mock based on configuration)
// 3. Initializes encryption and JWT services
// 4. Creates all business logic services with proper dependencies
// 5. Sets up authentication and project middlewares
// 6. Configures event subscriptions for cross-service communication
// 7. Creates default project if it doesn't exist
//
// Parameters:
//   - cnf: Application configuration containing all service settings
//
// Returns:
//   - *Provider: Fully configured provider with all dependencies
//   - error: Error if any initialization step fails
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
	am, err := auth.NewMiddlewares(svcs.Auth, svcs.Clients)
	if err != nil {
		return nil, err
	}
	authClient, err := goiamclient.GetGoIamClient(svcs.Clients)
	if err != nil {
		return nil, err
	}

	pvd := &Provider{
		S:          svcs,
		D:          d,
		C:          cS,
		PM:         pm,
		AM:         am,
		AuthClient: authClient,
	}

	// subscribe to user update events
	svcs.User.Subscribe(goiamuniverse.EventUserUpdated, svcs.AuthSync)

	// subscribe to client events for checking auth client
	svcs.Clients.Subscribe(goiamuniverse.EventClientCreated, pvd)
	svcs.Clients.Subscribe(goiamuniverse.EventClientUpdated, pvd)
	svcs.Clients.Subscribe(goiamuniverse.EventClientCreated, svcs.Auth)
	svcs.Clients.Subscribe(goiamuniverse.EventClientUpdated, svcs.Auth)

	// subscribe to resource events for updating downstream dependencies
	svcs.Resources.Subscribe(goiamuniverse.EventResourceDeleted, system.NewRemoveDeletedResourceFromRole(svcs.Role))
	svcs.Resources.Subscribe(goiamuniverse.EventResourceDeleted, system.NewRemoveDeletedResourceFromUser(svcs.User))

	// creating default project if it doesn't exist
	err = CheckAndAddDefaultProject(svcs.Projects)
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

// Handle creates a Fiber middleware that stores the Provider instance in the request context.
// This allows handlers and other middleware to access all services and dependencies
// throughout the request lifecycle.
//
// Usage:
//
//	app.Use(providers.Handle(provider))
//
// Parameters:
//   - p: Provider instance to store in context
//
// Returns:
//   - func(c *fiber.Ctx) error: Fiber middleware function
func Handle(p *Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(providerKey, p)
		return c.Next()
	}
}

// GetProviders retrieves the Provider instance from the Fiber context.
// This function should be called from handlers that need access to services
// and dependencies. The provider must have been previously stored using Handle middleware.
//
// Parameters:
//   - c: Fiber context containing the provider
//
// Returns:
//   - *Provider: Provider instance stored in the context
//
// Panics if no provider is found in the context.
func GetProviders(c *fiber.Ctx) *Provider {
	return c.Locals(providerKey).(*Provider)
}

// HandleEvent implements the event handler interface for client-related events.
// This method is automatically called when clients are created or updated,
// allowing the provider to update its authentication client configuration
// when Go IAM clients change.
//
// Event handling:
// - Listens for EventClientCreated and EventClientUpdated events
// - Updates AuthClient when Go IAM clients are modified
// - Propagates changes to authentication middleware
//
// Parameters:
//   - e: Event containing client information
func (p *Provider) HandleEvent(e utils.Event[sdk.Client]) {
	if e.Name() != goiamuniverse.EventClientCreated && e.Name() != goiamuniverse.EventClientUpdated {
		return
	}
	if !e.Payload().GoIamClient {
		return
	}
	var err error
	p.AuthClient, err = goiamclient.GetGoIamClient(p.S.Clients)
	if err != nil {
		log.Errorw("failed to get Go IAM client", "error", err)
		return
	}
	p.AM.AuthClient = p.AuthClient
}
