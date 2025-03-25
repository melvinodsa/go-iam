package role

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	Create(ctx context.Context, role *sdk.Role) error
	Update(ctx context.Context, role *sdk.Role) error
	GetById(ctx context.Context, id string) (*sdk.Role, error)
	GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error)
	AddRoleToUser(ctx context.Context, userId, roleId string) error
	RemoveRoleFromUser(ctx context.Context, userId, roleId string) error
}
