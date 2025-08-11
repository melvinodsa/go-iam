package auth

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl string) (string, error)
	Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error)
	ClientCallback(ctx context.Context, code string) (*sdk.AuthVerifyCodeResponse, error)
	GetIdentity(ctx context.Context, accessToken string) (*sdk.User, error)
	ClientCredentials(ctx context.Context, clientId, clientSecret string) (*sdk.ClientCredentialsDataResponse, error)

}
