package system

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewPolicyCheck(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	assert.NotNil(t, pc.userSvc)
}

func TestPolicyCheck_RunCheck_Success_PolicyExists(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	ctx := context.Background()
	userId := "user123"
	policyId := "@policy/system/test_policy"

	// Mock user with the policy
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			policyId: {Name: "Test Policy"},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	user, exists, err := pc.RunCheck(ctx, policyId, userId)

	// Verify
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, testUser, user)
	userSvc.AssertExpectations(t)
}

func TestPolicyCheck_RunCheck_Success_PolicyNotExists(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	ctx := context.Background()
	userId := "user123"
	policyId := "@policy/system/test_policy"

	// Mock user without the policy
	testUser := &sdk.User{
		Id:       userId,
		Policies: map[string]sdk.UserPolicy{}, // Empty policies
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	user, exists, err := pc.RunCheck(ctx, policyId, userId)

	// Verify
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Equal(t, testUser, user)
	userSvc.AssertExpectations(t)
}

func TestPolicyCheck_RunCheck_UserServiceError(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	ctx := context.Background()
	userId := "user123"
	policyId := "@policy/system/test_policy"

	// Mock user service to return error
	userSvc.On("GetById", ctx, userId).Return(nil, errors.New("user service error"))

	// Execute
	user, exists, err := pc.RunCheck(ctx, policyId, userId)

	// Verify
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Nil(t, user)
	assert.Equal(t, "user service error", err.Error())
	userSvc.AssertExpectations(t)
}

func TestPolicyCheck_RunCheck_PolicyExists_WithOtherPolicies(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	ctx := context.Background()
	userId := "user123"
	policyId := "@policy/system/test_policy"

	// Mock user with multiple policies including the target one
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/other_policy": {Name: "Other Policy"},
			policyId:                      {Name: "Test Policy"},
			"@policy/system/third_policy": {Name: "Third Policy"},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	user, exists, err := pc.RunCheck(ctx, policyId, userId)

	// Verify
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, testUser, user)
	userSvc.AssertExpectations(t)
}

func TestPolicyCheck_RunCheck_DifferentPolicyId(t *testing.T) {
	userSvc := &MockUserService{}
	pc := NewPolicyCheck(userSvc)

	ctx := context.Background()
	userId := "user123"
	policyId := "@policy/system/test_policy"
	differentPolicyId := "@policy/system/different_policy"

	// Mock user with a different policy
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			differentPolicyId: {Name: "Different Policy"},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	user, exists, err := pc.RunCheck(ctx, policyId, userId)

	// Verify
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Equal(t, testUser, user)
	userSvc.AssertExpectations(t)
}
