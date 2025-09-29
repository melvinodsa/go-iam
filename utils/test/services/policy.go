package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/mock"
)

type MockPolicyService struct {
	mock.Mock
}

func (m *MockPolicyService) GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.PolicyList), args.Error(1)
}