package user

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

type Service interface {
	Create(ctx context.Context, user *sdk.User) error
	Update(ctx context.Context, user *sdk.User) error
	GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error)
	GetById(ctx context.Context, id string) (*sdk.User, error)
	GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error)
	GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error)
	AddRoleToUser(ctx context.Context, userId, roleId string) error
	RemoveRoleFromUser(ctx context.Context, userId, roleId string) error
	AddResourceToUser(ctx context.Context, userId string, request sdk.AddUserResourceRequest) error
	AddPolicyToUser(ctx context.Context, userId string, policies map[string]sdk.UserPolicy) error
	RemovePolicyFromUser(ctx context.Context, userId string, policyIds []string) error
	HandleEvent(event utils.Event[sdk.Role])
	utils.Emitter[utils.Event[sdk.User], sdk.User]
}
