package role

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

type service struct {
	store Store
	e     utils.Emitter[utils.Event[sdk.Role], sdk.Role]
}

func NewService(store Store) Service {
	return &service{
		store: store,
		e:     utils.NewEmitter[utils.Event[sdk.Role]](),
	}
}
func (s *service) Create(ctx context.Context, role *sdk.Role) error {
	return s.store.Create(ctx, role)
}

func (s *service) Update(ctx context.Context, role *sdk.Role) error {
	err := s.store.Update(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	s.Emit(newEvent(ctx, goiamuniverse.EventRoleUpdated, *role, middlewares.GetMetadata(ctx)))
	return nil
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetAll(ctx context.Context, query sdk.RoleQuery) (*sdk.RoleList, error) {
	query.ProjectIds = middlewares.GetProjects(ctx)
	return s.store.GetAll(ctx, query)
}

func (s *service) AddResource(ctx context.Context, roleId string, resource sdk.Resources) error {
	role, err := s.GetById(ctx, roleId)
	if err != nil {
		return err
	}
	if len(role.Resources) == 0 {
		role.Resources = map[string]sdk.Resources{}
	}
	role.Resources[resource.Key] = resource
	return s.Update(ctx, role)
}

func (s service) Emit(event utils.Event[sdk.Role]) {
	if event == nil {
		return
	}
	s.e.Emit(event)
}

func (s service) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Role], sdk.Role]) {
	s.e.Subscribe(eventName, subscriber)
}

type event struct {
	name     goiamuniverse.Event
	payload  sdk.Role
	metadata sdk.Metadata
	ctx      context.Context
}

func (e event) Name() goiamuniverse.Event {
	return e.name
}

func (e event) Payload() sdk.Role {
	return e.payload
}

func (e event) Metadata() sdk.Metadata {
	return e.metadata
}

func (e event) Context() context.Context {
	return e.ctx
}

func newEvent(ctx context.Context, name goiamuniverse.Event, payload sdk.Role, metadata sdk.Metadata) utils.Event[sdk.Role] {
	return event{ctx: ctx, name: name, payload: payload, metadata: metadata}
}
