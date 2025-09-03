package authprovider

import (
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestFromModelToSdk(t *testing.T) {
	now := time.Now()
	model := &models.AuthProvider{
		Id:       "ap1",
		Name:     "Test Provider",
		Icon:     "test-icon",
		Provider: models.AuthProviderType("GOOGLE"),
		Params: []models.AuthProviderParam{
			{
				Label:    "Client ID",
				Value:    "client123",
				Key:      "@GOOGLE/CLIENT_ID",
				IsSecret: false,
			},
		},
		ProjectId: "project1",
		Enabled:   true,
		CreatedAt: &now,
		UpdatedAt: &now,
		CreatedBy: "user1",
		UpdatedBy: "user1",
	}

	result := fromModelToSdk(model)

	assert.NotNil(t, result)
	assert.Equal(t, "ap1", result.Id)
	assert.Equal(t, "Test Provider", result.Name)
	assert.Equal(t, "test-icon", result.Icon)
	assert.Equal(t, sdk.AuthProviderTypeGoogle, result.Provider)
	assert.Equal(t, "project1", result.ProjectId)
	assert.True(t, result.Enabled)
	assert.Equal(t, &now, result.CreatedAt)
	assert.Equal(t, &now, result.UpdatedAt)
	assert.Equal(t, "user1", result.CreatedBy)
	assert.Equal(t, "user1", result.UpdatedBy)
	assert.Len(t, result.Params, 1)
	assert.Equal(t, "Client ID", result.Params[0].Label)
	assert.Equal(t, "client123", result.Params[0].Value)
	assert.Equal(t, "@GOOGLE/CLIENT_ID", result.Params[0].Key)
	assert.False(t, result.Params[0].IsSecret)
}

func TestFromModelListToSdk(t *testing.T) {
	now := time.Now()
	models := []models.AuthProvider{
		{
			Id:        "ap1",
			Name:      "Provider 1",
			Provider:  models.AuthProviderType("GOOGLE"),
			Params:    []models.AuthProviderParam{},
			Enabled:   true,
			CreatedAt: &now,
		},
		{
			Id:        "ap2",
			Name:      "Provider 2",
			Provider:  models.AuthProviderType("GOOGLE"),
			Params:    []models.AuthProviderParam{},
			Enabled:   false,
			CreatedAt: &now,
		},
	}

	result := fromModelListToSdk(models)

	assert.Len(t, result, 2)
	assert.Equal(t, "ap1", result[0].Id)
	assert.Equal(t, "Provider 1", result[0].Name)
	assert.True(t, result[0].Enabled)
	assert.Equal(t, "ap2", result[1].Id)
	assert.Equal(t, "Provider 2", result[1].Name)
	assert.False(t, result[1].Enabled)
}

func TestFromProviderParamsModelToSdk(t *testing.T) {
	tests := []struct {
		name     string
		params   []models.AuthProviderParam
		expected []sdk.AuthProviderParam
	}{
		{
			name:     "empty_params",
			params:   []models.AuthProviderParam{},
			expected: []sdk.AuthProviderParam(nil),
		},
		{
			name: "single_param",
			params: []models.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
			},
			expected: []sdk.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
			},
		},
		{
			name: "multiple_params",
			params: []models.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
				{
					Label:    "Client Secret",
					Value:    "secret123",
					Key:      "@GOOGLE/CLIENT_SECRET",
					IsSecret: true,
				},
			},
			expected: []sdk.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
				{
					Label:    "Client Secret",
					Value:    "secret123",
					Key:      "@GOOGLE/CLIENT_SECRET",
					IsSecret: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromProviderParamsModelToSdk(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFromSdkToModel(t *testing.T) {
	now := time.Now()
	sdk := sdk.AuthProvider{
		Id:       "ap1",
		Name:     "Test Provider",
		Icon:     "test-icon",
		Provider: sdk.AuthProviderTypeGoogle,
		Params: []sdk.AuthProviderParam{
			{
				Label:    "Client ID",
				Value:    "client123",
				Key:      "@GOOGLE/CLIENT_ID",
				IsSecret: false,
			},
		},
		ProjectId: "project1",
		Enabled:   true,
		CreatedAt: &now,
		UpdatedAt: &now,
		CreatedBy: "user1",
		UpdatedBy: "user1",
	}

	result := fromSdkToModel(sdk)

	assert.Equal(t, "ap1", result.Id)
	assert.Equal(t, "Test Provider", result.Name)
	assert.Equal(t, "test-icon", result.Icon)
	assert.Equal(t, models.AuthProviderType("GOOGLE"), result.Provider)
	assert.Equal(t, "project1", result.ProjectId)
	assert.True(t, result.Enabled)
	assert.Equal(t, &now, result.CreatedAt)
	assert.Equal(t, &now, result.UpdatedAt)
	assert.Equal(t, "user1", result.CreatedBy)
	assert.Equal(t, "user1", result.UpdatedBy)
	assert.Len(t, result.Params, 1)
	assert.Equal(t, "Client ID", result.Params[0].Label)
	assert.Equal(t, "client123", result.Params[0].Value)
	assert.Equal(t, "@GOOGLE/CLIENT_ID", result.Params[0].Key)
	assert.False(t, result.Params[0].IsSecret)
}

func TestFromProviderParamsSdkToModel(t *testing.T) {
	tests := []struct {
		name     string
		params   []sdk.AuthProviderParam
		expected []models.AuthProviderParam
	}{
		{
			name:     "empty_params",
			params:   []sdk.AuthProviderParam{},
			expected: []models.AuthProviderParam(nil),
		},
		{
			name: "single_param",
			params: []sdk.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
			},
			expected: []models.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
			},
		},
		{
			name: "multiple_params_with_secrets",
			params: []sdk.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
				{
					Label:    "Client Secret",
					Value:    "secret123",
					Key:      "@GOOGLE/CLIENT_SECRET",
					IsSecret: true,
				},
			},
			expected: []models.AuthProviderParam{
				{
					Label:    "Client ID",
					Value:    "client123",
					Key:      "@GOOGLE/CLIENT_ID",
					IsSecret: false,
				},
				{
					Label:    "Client Secret",
					Value:    "secret123",
					Key:      "@GOOGLE/CLIENT_SECRET",
					IsSecret: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromProviderParamsSdkToModel(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}
