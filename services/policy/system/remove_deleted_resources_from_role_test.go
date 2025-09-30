package system

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
)


func TestNewRemoveDeletedResourceFromRole(t *testing.T) {
	roleSvc := &services.MockRoleService{}
	policy := NewRemoveDeletedResourceFromRole(roleSvc)

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_role", policy.ID())
	assert.Equal(t, "Remove deleted resources from role specified in user policy", policy.Name())
	assert.NotNil(t, policy.roleSvc)
}

func TestRemoveDeletedResourceFromRole_ID(t *testing.T) {
	roleSvc := &services.MockRoleService{}
	policy := NewRemoveDeletedResourceFromRole(roleSvc)

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_role", policy.ID())
}

func TestRemoveDeletedResourceFromRole_Name(t *testing.T) {
	roleSvc := &services.MockRoleService{}
	policy := NewRemoveDeletedResourceFromRole(roleSvc)

	assert.Equal(t, "Remove deleted resources from role specified in user policy", policy.Name())
}

func TestRemoveDeletedResourceFromRole_HandleEvent(t *testing.T) {
	roleSvc := &services.MockRoleService{}

	ctx := context.Background()
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceDeleted,
		resource,
		sdk.Metadata{},
	)

	roleSvc.On("RemoveResourceFromAll", ctx, resource.Key).Return(nil)

	policy := NewRemoveDeletedResourceFromRole(roleSvc)
	policy.HandleEvent(event)

	roleSvc.AssertExpectations(t)
}

func TestRemoveDeletedResourceFromRole_HandleEvent_Error(t *testing.T) {
	roleSvc := &services.MockRoleService{}

	ctx := context.Background()
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceDeleted,
		resource,
		sdk.Metadata{},
	)

	roleSvc.On("RemoveResourceFromAll", ctx, resource.Key).Return(errors.New("remove error"))

	policy := NewRemoveDeletedResourceFromRole(roleSvc)
	policy.HandleEvent(event)

	roleSvc.AssertExpectations(t)
}

func TestRemoveDeletedResourceFromRole_PolicyDef(t *testing.T) {
	roleSvc := &services.MockRoleService{}
	policy := NewRemoveDeletedResourceFromRole(roleSvc)

	policyDef := policy.PolicyDef()

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_role", policyDef.Id)
	assert.Equal(t, "Remove deleted resources from role specified in user policy", policyDef.Name)
	assert.Equal(t, "This policy removes the deleted resource from all roles.", policyDef.Description)
	assert.Empty(t, policyDef.Definition.Arguments)
}