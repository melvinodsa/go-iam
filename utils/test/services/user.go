package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user *sdk.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Update(ctx context.Context, user *sdk.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error) {
	args := m.Called(ctx, email, projectId)
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockUserService) GetById(ctx context.Context, id string) (*sdk.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockUserService) GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error) {
	args := m.Called(ctx, phone, projectId)
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sdk.UserList), args.Error(1)
}

func (m *MockUserService) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	args := m.Called(ctx, userId, roleId)
	return args.Error(0)
}

func (m *MockUserService) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	args := m.Called(ctx, userId, roleId)
	return args.Error(0)
}

func (m *MockUserService) AddResourceToUser(ctx context.Context, userId string, request sdk.AddUserResourceRequest) error {
	args := m.Called(ctx, userId, request)
	return args.Error(0)
}

func (m *MockUserService) AddPolicyToUser(ctx context.Context, userId string, policies map[string]sdk.UserPolicy) error {
	args := m.Called(ctx, userId, policies)
	return args.Error(0)
}

func (m *MockUserService) RemovePolicyFromUser(ctx context.Context, userId string, policyIds []string) error {
	args := m.Called(ctx, userId, policyIds)
	return args.Error(0)
}

func (m *MockUserService) RemoveResourceFromAll(ctx context.Context, resourceKey string) error {
	args := m.Called(ctx, resourceKey)
	return args.Error(0)
}

func (m *MockUserService) TransferOwnership(ctx context.Context, userId, newOwnerId string) error {
	args := m.Called(ctx, userId, newOwnerId)
	return args.Error(0)
}

func (m *MockUserService) HandleEvent(event utils.Event[sdk.Role]) {
	m.Called(event)
}

func (m *MockUserService) Emit(event utils.Event[sdk.User]) {
	m.Called(event)
}

func (m *MockUserService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.User], sdk.User]) {
	m.Called(eventName, subscriber)
}

func (m *MockUserService) CopyUserResources(ctx context.Context, sourceUserId, targetUserId string) error {
	args := m.Called(ctx, sourceUserId, targetUserId)
	return args.Error(0)
}
