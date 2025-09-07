package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/mock"
)

type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) Create(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}
func (m *MockRoleService) Update(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}
func (m *MockRoleService) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.Role), args.Error(1)
}
func (m *MockRoleService) GetAll(ctx context.Context, query sdk.RoleQuery) (*sdk.RoleList, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sdk.RoleList), args.Error(1)
}
func (m *MockRoleService) AddResource(ctx context.Context, roleId string, resource sdk.Resources) error {
	args := m.Called(ctx, roleId, resource)
	return args.Error(0)
}
func (m *MockRoleService) Emit(event utils.Event[sdk.Role]) {
	m.Called(event)
}

func (m *MockRoleService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Role], sdk.Role]) {
	m.Called(eventName, subscriber)
}
