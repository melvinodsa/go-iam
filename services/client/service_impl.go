package client

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/middlewares/projects"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils"
)

type service struct {
	s Store
	p project.Service
	e utils.Emitter[utils.Event[sdk.Client], sdk.Client]
}

func NewService(s Store, p project.Service) Service {
	return service{s: s, p: p, e: utils.NewEmitter[utils.Event[sdk.Client]]()}
}

func (s service) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	queryParams.ProjectIds = projects.GetProjects(ctx)
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

	projectIdsMap := utils.Reduce(projects.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[cl.ProjectId]; !ok {
		return nil, ErrClientNotFound
	}
	return nil, fmt.Errorf("client not found")
}
func (s service) Create(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(projects.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
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
	s.Emit(newEvent(sdk.EventClientCreated, *client))
	return nil
}
func (s service) Update(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(projects.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
	}
	err := s.s.Update(ctx, client)
	if err != nil {
		return fmt.Errorf("error while updating client: %w", err)
	}
	s.Emit(newEvent(sdk.EventClientUpdated, *client))
	return nil
}

func (s service) Emit(event utils.Event[sdk.Client]) {
	if event == nil {
		return
	}
	s.e.Emit(event)
}

func (s service) Subscribe(eventName string, subscriber utils.Subscriber[utils.Event[sdk.Client], sdk.Client]) {
	s.e.Subscribe(eventName, subscriber)
}

type event struct {
	name    string
	payload sdk.Client
}

func (e event) Name() string {
	return e.name
}

func (e event) Payload() sdk.Client {
	return e.payload
}

func newEvent(name string, payload sdk.Client) utils.Event[sdk.Client] {
	return event{name: name, payload: payload}
}
