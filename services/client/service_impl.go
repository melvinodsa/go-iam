package client

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils"
)

type service struct {
	s Store
	p project.Service
}

func NewService(s Store, p project.Service) Service {
	return service{s: s, p: p}
}

func (s service) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	queryParams.ProjectIds = utils.Map(middlewares.GetProjects(ctx), func(p sdk.Project) string { return p.Id })
	return s.s.GetAll(ctx, queryParams)
}
func (s service) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error) {
	cl, err := s.s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if dontCheckProjects {
		return cl, nil
	}

	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p sdk.Project) map[string]bool { ini[p.Id] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[cl.ProjectId]; !ok {
		return nil, ErrClientNotFound
	}
	return nil, fmt.Errorf("client not found")
}
func (s service) Create(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p sdk.Project) map[string]bool { ini[p.Id] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
	}
	// create a random string secret
	sec, err := generateRandomSecret(32)
	if err != nil {
		return fmt.Errorf("error while creating client secret: %w", err)
	}
	client.Secret = sec
	return s.s.Create(ctx, client)
}
func (s service) Update(ctx context.Context, client *sdk.Client) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p sdk.Project) map[string]bool { ini[p.Id] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[client.ProjectId]; !ok {
		return project.ErrProjectNotFound
	}
	return s.s.Update(ctx, client)
}
