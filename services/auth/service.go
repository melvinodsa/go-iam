package auth

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

type Service interface {
	GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeVerifier string) (string, error)
	Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error)
	ClientCallback(ctx context.Context, code, codeVerifier, clientId, clietSecret string) (*sdk.AuthVerifyCodeResponse, error)
	GetIdentity(ctx context.Context, accessToken string) (*sdk.User, error)
	HandleEvent(event utils.Event[sdk.Client])
}
