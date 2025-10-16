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

func TestNewRemoveDeletedResourceFromUser(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewRemoveDeletedResourceFromUser(userSvc)

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_user", policy.ID())
	assert.Equal(t, "Remove deleted resources from user specified in user policy", policy.Name())
	assert.NotNil(t, policy.userSvc)
}

func TestRemoveDeletedResourceFromUser_ID(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewRemoveDeletedResourceFromUser(userSvc)

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_user", policy.ID())
}

func TestRemoveDeletedResourceFromUser_Name(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewRemoveDeletedResourceFromUser(userSvc)

	assert.Equal(t, "Remove deleted resources from user specified in user policy", policy.Name())
}

func TestRemoveDeletedResourceFromUser_HandleEvent(t *testing.T) {
	userSvc := &services.MockUserService{}

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

	userSvc.On("RemoveResourceFromAll", ctx, resource.Key).Return(nil)

	policy := NewRemoveDeletedResourceFromUser(userSvc)
	policy.HandleEvent(event)

	userSvc.AssertExpectations(t)
}

func TestRemoveDeletedResourceFromUser_HandleEvent_Error(t *testing.T) {
	userSvc := &services.MockUserService{}

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

	userSvc.On("RemoveResourceFromAll", ctx, resource.Key).Return(errors.New("remove error"))

	policy := NewRemoveDeletedResourceFromUser(userSvc)
	policy.HandleEvent(event)

	userSvc.AssertExpectations(t)
}

func TestRemoveDeletedResourceFromUser_PolicyDef(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewRemoveDeletedResourceFromUser(userSvc)

	policyDef := policy.PolicyDef()

	assert.Equal(t, "@policy/system/remove_deleted_resources_from_user", policyDef.Id)
	assert.Equal(t, "Remove deleted resources from user specified in user policy", policyDef.Name)
	assert.Equal(t, "This policy removes the deleted resource from all users.", policyDef.Description)
	assert.Empty(t, policyDef.Definition.Arguments)
}