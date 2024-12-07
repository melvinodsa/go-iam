package client

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	s Store
}

func NewService(s Store) Service {
	return service{s: s}
}

func (s service) GetAll(ctx context.Context) ([]sdk.Client, error) {
	return s.s.GetAll(ctx)
}
func (s service) Get(ctx context.Context, id string) (*sdk.Client, error) {
	return s.s.Get(ctx, id)
}
func (s service) Create(ctx context.Context, client *sdk.Client) error {
	// create a random string secret
	sec, err := generateRandomSecret(32)
	if err != nil {
		return fmt.Errorf("error while creating client secret: %w", err)
	}
	client.Secret = sec
	return s.s.Create(ctx, client)
}
func (s service) Update(ctx context.Context, client *sdk.Client) error {
	return s.s.Update(ctx, client)
}
