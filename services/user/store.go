package user

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	Create(ctx context.Context, user *sdk.User) error
	Update(ctx context.Context, user *sdk.User) error
	GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error)
	GetById(ctx context.Context, id string) (*sdk.User, error)
	GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error)
	GetAll(ctx context.Context, query sdk.UserQuery) ([]sdk.User, error)
}
