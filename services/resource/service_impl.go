package resource

import (
	"context"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

type service struct {
	s Store
	e utils.Emitter[utils.Event[sdk.Resource], sdk.Resource]
}

func NewService(s Store) Service {
	return service{s: s,
		e: utils.NewEmitter[utils.Event[sdk.Resource]]()}
}

func (s service) Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error) {
	query.ProjectIds = middlewares.GetProjects(ctx)
	return s.s.Search(ctx, query)
}

func (s service) Get(ctx context.Context, id string) (*sdk.Resource, error) {
	return s.s.Get(ctx, id)
}

func (s service) Create(ctx context.Context, resource *sdk.Resource) error {
	_, err := s.s.Create(ctx, resource)
	s.Emit(newEvent(ctx, utils.EventResourceCreated, *resource, middlewares.GetMetadata(ctx)))
	return err
}

func (s service) Update(ctx context.Context, resource *sdk.Resource) error {
	return s.s.Update(ctx, resource)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.s.Delete(ctx, id)
}

func (s service) Emit(event utils.Event[sdk.Resource]) {
	if event == nil {
		return
	}
	s.e.Emit(event)
}

func (s service) Subscribe(eventName string, subscriber utils.Subscriber[utils.Event[sdk.Resource], sdk.Resource]) {
	s.e.Subscribe(eventName, subscriber)
}

type event struct {
	name     string
	payload  sdk.Resource
	metadata sdk.Metadata
	ctx      context.Context
}

func (e event) Name() string {
	return e.name
}

func (e event) Payload() sdk.Resource {
	return e.payload
}

func (e event) Metadata() sdk.Metadata {
	return e.metadata
}

func (e event) Context() context.Context {
	return e.ctx
}

func newEvent(ctx context.Context, name string, payload sdk.Resource, metadata sdk.Metadata) utils.Event[sdk.Resource] {
	return event{ctx: ctx, name: name, payload: payload, metadata: metadata}
}
