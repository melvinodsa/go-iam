package auth

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl string, redis string) (string, error)
	Redirect(ctx context.Context, code, state string, redis string) (*sdk.AuthRedirectResponse, error)
	ClientCallback(ctx context.Context, code string, redis string) (*sdk.AuthVerifyCodeResponse, error)
	GetIdentity(ctx context.Context, accessToken string, redis string) (*sdk.User, error)
}
