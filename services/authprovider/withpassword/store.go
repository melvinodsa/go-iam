package withpassword

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetUserByUsername(ctx context.Context, email, projectID, password string) (*sdk.WithPasswordUser, error)
	CreateUser(ctx context.Context, email string, projectID string, password string) error
	UpdateUserPassword(ctx context.Context, email, projectID, newPassword string) error
}
