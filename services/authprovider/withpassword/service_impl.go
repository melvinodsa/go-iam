package withpassword

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) Signup(ctx context.Context, email, password, projectID string) error {
	return s.store.CreateUser(ctx, email, projectID, password)
}

func (s *service) Login(ctx context.Context, email, password, projectID string) (*sdk.WithPasswordUser, error) {
	user, err := s.store.GetUserByUsername(ctx, email, projectID, password)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}
func (s *service) UpdatePassword(ctx context.Context, email, projectID, oldPassword, newPassword string) error {
	_, err := s.store.GetUserByUsername(ctx, email, projectID, oldPassword)
	if err != nil {
		return err
	}
	return s.store.UpdateUserPassword(ctx, email, projectID, newPassword)
}
