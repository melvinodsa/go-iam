package client

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/services/user"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type service struct {
	s      Store
	p      project.Service
	e      utils.Emitter[utils.Event[sdk.Client], sdk.Client]
	authP  authprovider.Service
	usrSvc user.Service
}

func NewService(s Store, p project.Service, authP authprovider.Service, usrSvc user.Service) Service {
	return service{
		s:      s,
		p:      p,
		e:      utils.NewEmitter[utils.Event[sdk.Client]](),
		authP:  authP,
		usrSvc: usrSvc,
	}
}

func (s service) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	queryParams.ProjectIds = middlewares.GetProjects(ctx)
	return s.s.GetAll(ctx, queryParams)
}

func (s service) GetGoIamClients(ctx context.Context, params sdk.ClientQueryParams) ([]sdk.Client, error) {
	params.GoIamClient = true
	params.SortByUpdatedAt = true
	providers, err := s.s.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (s service) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error) {
	cl, err := s.s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if dontCheckProjects {
		return cl, nil
	}

	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[cl.ProjectId]; !ok {
		return nil, sdk.ErrClientNotFound
	}
	return cl, nil
}
func (s service) Create(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
	}
	if client.DefaultAuthProviderId != "" {
		// verifying if auth provider exists
		_, err := s.authP.Get(ctx, client.DefaultAuthProviderId, true)
		if err != nil {
			return fmt.Errorf("failed to get auth provider: %w", err)
		}
	}
	// create a random string secret
	sec, err := generateRandomSecret(32)
	if err != nil {
		return fmt.Errorf("error while creating client secret: %w", err)
	}
	client.Secret = sec
	err = s.s.Create(ctx, client)
	if err != nil {
		return fmt.Errorf("error while creating client: %w", err)
	}
	hashedSec, err := hashSecret(sec)
	if err != nil {
		return fmt.Errorf("error hashing client secret: %w", err)
	}
	client.Secret = hashedSec
	err = s.createAndLinkServiceAccountUser(ctx, client, client.ServiceAccountEmail)
	if err != nil {
		return fmt.Errorf("failed to create service account user: %w", err)
	}
	client.Secret = sec // return the plain secret only during creation
	s.Emit(newEvent(ctx, goiamuniverse.EventClientCreated, *client, middlewares.GetMetadata(ctx)))
	return nil
}
func (s service) Update(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
	}
	err := s.s.Update(ctx, client)
	if err != nil {
		return fmt.Errorf("error while updating client: %w", err)
	}
	s.Emit(newEvent(ctx, goiamuniverse.EventClientUpdated, *client, middlewares.GetMetadata(ctx)))
	return nil
}

func (s service) VerifySecret(plainSecret, hashedSecret string) error {
	if plainSecret == "" {
		return fmt.Errorf("plain secret cannot be empty")
	}

	// hashSecret function is defined in services/client/helpers.go
	// It hashes the secret using SHA256 and encodes to base64
	hashedPlain, err := hashSecret(plainSecret)
	if err != nil {
		return fmt.Errorf("failed to hash secret for verification: %w", err)
	}

	if hashedPlain != hashedSecret {
		return fmt.Errorf("secret verification failed: invalid secret")
	}

	return nil
}

func (s service) RegenerateSecret(ctx context.Context, clientId string) (*sdk.Client, error) {
	client, err := s.Get(ctx, clientId, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching client: %w", err)
	}

	// Generate a new random secret
	newSecret, err := generateRandomSecret(32)
	if err != nil {
		return nil, fmt.Errorf("error generating new client secret: %w", err)
	}

	hashedSec, err := hashSecret(newSecret)
	if err != nil {
		return nil, fmt.Errorf("error hashing client secret: %w", err)
	}
	client.Secret = hashedSec
	err = s.s.Update(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error updating client with new secret: %w", err)
	}
	client.Secret = newSecret // return the plain secret only during regeneration

	return client, nil
}

func (s service) createAndLinkServiceAccountUser(ctx context.Context, client *sdk.Client, email string) error {
	user, err := s.createServiceAccountUser(ctx, client, email)
	if err != nil {
		return fmt.Errorf("failed to create service account user: %w", err)
	}

	client.LinkedUserId = user.Id
	err = s.Update(ctx, client)
	if err != nil {
		log.Debugw("failed to update client with linked user", "client_id", client.Id, "user_id", user.Id, "error", err)
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

func (s service) createServiceAccountUser(ctx context.Context, client *sdk.Client, email string) (*sdk.User, error) {
	creator := middlewares.GetUser(ctx)
	// Implementation for creating a service account user
	user := &sdk.User{
		Email:          email,
		ProjectId:      client.ProjectId,
		LinkedClientId: client.Id,
		Name:           fmt.Sprintf("Service Account of %s", client.Name),
		CreatedBy:      creator.Email,
	}

	err := s.usrSvc.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account user: %w", err)
	}

	return user, nil
}

func (s service) Emit(event utils.Event[sdk.Client]) {
	if event == nil {
		return
	}
	s.e.Emit(event)
}

func (s service) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Client], sdk.Client]) {
	s.e.Subscribe(eventName, subscriber)
}

type event struct {
	name     goiamuniverse.Event
	payload  sdk.Client
	metadata sdk.Metadata
	ctx      context.Context
}

func (e event) Name() goiamuniverse.Event {
	return e.name
}

func (e event) Payload() sdk.Client {
	return e.payload
}

func (e event) Metadata() sdk.Metadata {
	return e.metadata
}

func (e event) Context() context.Context {
	return e.ctx
}

func newEvent(ctx context.Context, name goiamuniverse.Event, payload sdk.Client, metadata sdk.Metadata) utils.Event[sdk.Client] {
	return event{ctx: ctx, name: name, payload: payload, metadata: metadata}
}
