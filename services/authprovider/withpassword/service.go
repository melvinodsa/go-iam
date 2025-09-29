package withpassword

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	Signup(ctx context.Context, email, password, projectID string) error
	Login(ctx context.Context, email, password, projectID string) (*sdk.WithPasswordUser, error)
	UpdatePassword(ctx context.Context, email, projectID, oldPassword, newPassword string) error
}
