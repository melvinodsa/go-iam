package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of auth.Service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeChallenge string) (string, error) {
	args := m.Called(ctx, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeChallenge)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error) {
	args := m.Called(ctx, code, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthRedirectResponse), args.Error(1)
}

func (m *MockAuthService) ClientCallback(ctx context.Context, code, codeChallenge, clientId, clietSecret string) (*sdk.AuthVerifyCodeResponse, error) {
	args := m.Called(ctx, code, codeChallenge, clientId, clietSecret)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthVerifyCodeResponse), args.Error(1)
}

func (m *MockAuthService) GetIdentity(ctx context.Context, token string, forceFetch bool) (*sdk.User, error) {
	args := m.Called(ctx, token, forceFetch)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockAuthService) HandleEvent(event utils.Event[sdk.Client]) {
	m.Called(event)
}

func (m *MockAuthService) ClientCredentials(ctx context.Context, clientId, clientSecret string) (*sdk.AuthVerifyCodeResponse, error) {
	args := m.Called(ctx, clientId, clientSecret)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthVerifyCodeResponse), args.Error(1)
}
