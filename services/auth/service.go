package auth

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

type Service interface {
	GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeChallenge string) (string, error)
	Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error)
	ClientCallback(ctx context.Context, code, codeChallenge, clientId, clietSecret string) (*sdk.AuthVerifyCodeResponse, error)
	GetIdentity(ctx context.Context, accessToken string) (*sdk.User, error)
	ClientCredentials(ctx context.Context, clientId, clientSecret string) (*sdk.ClientCredentialsDataResponse, error)
	HandleEvent(event utils.Event[sdk.Client])
}
