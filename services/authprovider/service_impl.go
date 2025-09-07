package authprovider

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/authprovider/github"
	"github.com/melvinodsa/go-iam/services/authprovider/google"
	"github.com/melvinodsa/go-iam/services/authprovider/microsoft"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils"
)

type service struct {
	s Store
	p project.Service
}

func NewService(s Store, p project.Service) Service {
	return &service{
		s: s,
		p: p,
	}
}

func (s service) GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error) {
	params.ProjectIds = middlewares.GetProjects(ctx)
	return s.s.GetAll(ctx, params)
}

func (s service) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.AuthProvider, error) {
	ap, err := s.s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if dontCheckProjects {
		return ap, nil
	}

	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[ap.ProjectId]; !ok {
		return nil, ErrAuthProviderNotFound
	}
	return ap, nil
}
func (s service) Create(ctx context.Context, provider *sdk.AuthProvider) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[provider.ProjectId]; !ok {
		return sdk.ErrProjectNotFound
	}
	return s.s.Create(ctx, provider)
}
func (s service) Update(ctx context.Context, provider *sdk.AuthProvider) error {
	// check if the project exists
	projectIdsMap := utils.Reduce(middlewares.GetProjects(ctx), func(ini map[string]bool, p string) map[string]bool { ini[p] = true; return ini }, map[string]bool{})
	if _, ok := projectIdsMap[provider.ProjectId]; !ok {
		return sdk.ErrProjectNotFound
	}
	return s.s.Update(ctx, provider)
}

func (s service) GetProvider(ctx context.Context, v sdk.AuthProvider) (sdk.ServiceProvider, error) {
	switch v.Provider {
	case sdk.AuthProviderTypeGoogle:
		return google.NewAuthProvider(v), nil
	case sdk.AuthProviderTypeMicrosoft:
		return microsoft.NewAuthProvider(v), nil
	case sdk.AuthProviderTypeGitHub:
		return github.NewAuthProvider(v), nil
	default:
		return nil, fmt.Errorf("unknown auth provider: %s", v.Provider)
	}
}
