package client

import (
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestFromModelListToSdk(t *testing.T) {
	t.Run("empty_list", func(t *testing.T) {
		modelClients := []models.Client{}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})
	
	t.Run("single_client", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(time.Hour)
		
		modelClients := []models.Client{
			{
				Id:                    "client1",
				Name:                  "Test Client",
				Description:           "A test client",
				Secret:                "secret123",
				Tags:                  []string{"tag1", "tag2"},
				RedirectURLs:          []string{"https://example.com/callback"},
				DefaultAuthProviderId: "provider1",
				GoIamClient:           true,
				ProjectId:             "project1",
				Scopes:                []string{"read", "write"},
				Enabled:               true,
				CreatedAt:             &createdAt,
				CreatedBy:             "user1",
				UpdatedAt:             &updatedAt,
				UpdatedBy:             "user2",
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		
		client := result[0]
		assert.Equal(t, "client1", client.Id)
		assert.Equal(t, "Test Client", client.Name)
		assert.Equal(t, "A test client", client.Description)
		assert.Equal(t, "secret123", client.Secret)
		assert.Equal(t, []string{"tag1", "tag2"}, client.Tags)
		assert.Equal(t, []string{"https://example.com/callback"}, client.RedirectURLs)
		assert.Equal(t, "provider1", client.DefaultAuthProviderId)
		assert.True(t, client.GoIamClient)
		assert.Equal(t, "project1", client.ProjectId)
		assert.Equal(t, []string{"read", "write"}, client.Scopes)
		assert.True(t, client.Enabled)
		assert.Equal(t, &createdAt, client.CreatedAt)
		assert.Equal(t, "user1", client.CreatedBy)
		assert.Equal(t, &updatedAt, client.UpdatedAt)
		assert.Equal(t, "user2", client.UpdatedBy)
	})
	
	t.Run("multiple_clients", func(t *testing.T) {
		createdAt1 := time.Now()
		updatedAt1 := time.Now().Add(time.Hour)
		createdAt2 := time.Now().Add(time.Minute)
		updatedAt2 := time.Now().Add(2 * time.Hour)
		
		modelClients := []models.Client{
			{
				Id:                    "client1",
				Name:                  "First Client",
				Description:           "First test client",
				Secret:                "secret1",
				Tags:                  []string{"tag1"},
				RedirectURLs:          []string{"https://first.com/callback"},
				DefaultAuthProviderId: "provider1",
				GoIamClient:           true,
				ProjectId:             "project1",
				Scopes:                []string{"read"},
				Enabled:               true,
				CreatedAt:             &createdAt1,
				CreatedBy:             "user1",
				UpdatedAt:             &updatedAt1,
				UpdatedBy:             "user1",
			},
			{
				Id:                    "client2",
				Name:                  "Second Client",
				Description:           "Second test client",
				Secret:                "secret2",
				Tags:                  []string{"tag2", "tag3"},
				RedirectURLs:          []string{"https://second.com/callback", "https://second.com/alt"},
				DefaultAuthProviderId: "provider2",
				GoIamClient:           false,
				ProjectId:             "project2",
				Scopes:                []string{"read", "write", "admin"},
				Enabled:               false,
				CreatedAt:             &createdAt2,
				CreatedBy:             "user2",
				UpdatedAt:             &updatedAt2,
				UpdatedBy:             "user3",
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result))
		
		// Check first client
		client1 := result[0]
		assert.Equal(t, "client1", client1.Id)
		assert.Equal(t, "First Client", client1.Name)
		assert.Equal(t, "First test client", client1.Description)
		assert.Equal(t, "secret1", client1.Secret)
		assert.Equal(t, []string{"tag1"}, client1.Tags)
		assert.Equal(t, []string{"https://first.com/callback"}, client1.RedirectURLs)
		assert.Equal(t, "provider1", client1.DefaultAuthProviderId)
		assert.True(t, client1.GoIamClient)
		assert.Equal(t, "project1", client1.ProjectId)
		assert.Equal(t, []string{"read"}, client1.Scopes)
		assert.True(t, client1.Enabled)
		assert.Equal(t, &createdAt1, client1.CreatedAt)
		assert.Equal(t, "user1", client1.CreatedBy)
		assert.Equal(t, &updatedAt1, client1.UpdatedAt)
		assert.Equal(t, "user1", client1.UpdatedBy)
		
		// Check second client
		client2 := result[1]
		assert.Equal(t, "client2", client2.Id)
		assert.Equal(t, "Second Client", client2.Name)
		assert.Equal(t, "Second test client", client2.Description)
		assert.Equal(t, "secret2", client2.Secret)
		assert.Equal(t, []string{"tag2", "tag3"}, client2.Tags)
		assert.Equal(t, []string{"https://second.com/callback", "https://second.com/alt"}, client2.RedirectURLs)
		assert.Equal(t, "provider2", client2.DefaultAuthProviderId)
		assert.False(t, client2.GoIamClient)
		assert.Equal(t, "project2", client2.ProjectId)
		assert.Equal(t, []string{"read", "write", "admin"}, client2.Scopes)
		assert.False(t, client2.Enabled)
		assert.Equal(t, &createdAt2, client2.CreatedAt)
		assert.Equal(t, "user2", client2.CreatedBy)
		assert.Equal(t, &updatedAt2, client2.UpdatedAt)
		assert.Equal(t, "user3", client2.UpdatedBy)
	})
	
	t.Run("clients_with_nil_timestamps", func(t *testing.T) {
		modelClients := []models.Client{
			{
				Id:          "client1",
				Name:        "Client with nil timestamps",
				Description: "Test client",
				ProjectId:   "project1",
				CreatedAt:   nil,
				UpdatedAt:   nil,
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		
		client := result[0]
		assert.Equal(t, "client1", client.Id)
		assert.Equal(t, "Client with nil timestamps", client.Name)
		assert.Nil(t, client.CreatedAt)
		assert.Nil(t, client.UpdatedAt)
	})
	
	t.Run("clients_with_empty_arrays", func(t *testing.T) {
		modelClients := []models.Client{
			{
				Id:           "client1",
				Name:         "Client with empty arrays",
				Description:  "Test client",
				ProjectId:    "project1",
				Tags:         []string{},
				RedirectURLs: []string{},
				Scopes:       []string{},
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		
		client := result[0]
		assert.Equal(t, "client1", client.Id)
		assert.Equal(t, "Client with empty arrays", client.Name)
		assert.Empty(t, client.Tags)
		assert.Empty(t, client.RedirectURLs)
		assert.Empty(t, client.Scopes)
	})
	
	t.Run("clients_with_nil_arrays", func(t *testing.T) {
		modelClients := []models.Client{
			{
				Id:           "client1",
				Name:         "Client with nil arrays",
				Description:  "Test client",
				ProjectId:    "project1",
				Tags:         nil,
				RedirectURLs: nil,
				Scopes:       nil,
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		
		client := result[0]
		assert.Equal(t, "client1", client.Id)
		assert.Equal(t, "Client with nil arrays", client.Name)
		assert.Nil(t, client.Tags)
		assert.Nil(t, client.RedirectURLs)
		assert.Nil(t, client.Scopes)
	})
	
	t.Run("clients_with_default_values", func(t *testing.T) {
		modelClients := []models.Client{
			{
				Id:                    "",
				Name:                  "",
				Description:           "",
				Secret:                "",
				DefaultAuthProviderId: "",
				GoIamClient:           false,
				ProjectId:             "",
				Enabled:               false,
				CreatedBy:             "",
				UpdatedBy:             "",
			},
		}
		
		result := fromModelListToSdk(modelClients)
		
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		
		client := result[0]
		assert.Empty(t, client.Id)
		assert.Empty(t, client.Name)
		assert.Empty(t, client.Description)
		assert.Empty(t, client.Secret)
		assert.Empty(t, client.DefaultAuthProviderId)
		assert.False(t, client.GoIamClient)
		assert.Empty(t, client.ProjectId)
		assert.False(t, client.Enabled)
		assert.Empty(t, client.CreatedBy)
		assert.Empty(t, client.UpdatedBy)
	})
}

func TestFromModelToSdk(t *testing.T) {
	t.Run("complete_client", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(time.Hour)
		
		modelClient := &models.Client{
			Id:                    "client1",
			Name:                  "Test Client",
			Description:           "A test client",
			Secret:                "secret123",
			Tags:                  []string{"tag1", "tag2"},
			RedirectURLs:          []string{"https://example.com/callback"},
			DefaultAuthProviderId: "provider1",
			GoIamClient:           true,
			ProjectId:             "project1",
			Scopes:                []string{"read", "write"},
			Enabled:               true,
			CreatedAt:             &createdAt,
			CreatedBy:             "user1",
			UpdatedAt:             &updatedAt,
			UpdatedBy:             "user2",
		}
		
		result := fromModelToSdk(modelClient)
		
		assert.NotNil(t, result)
		assert.Equal(t, "client1", result.Id)
		assert.Equal(t, "Test Client", result.Name)
		assert.Equal(t, "A test client", result.Description)
		assert.Equal(t, "secret123", result.Secret)
		assert.Equal(t, []string{"tag1", "tag2"}, result.Tags)
		assert.Equal(t, []string{"https://example.com/callback"}, result.RedirectURLs)
		assert.Equal(t, "provider1", result.DefaultAuthProviderId)
		assert.True(t, result.GoIamClient)
		assert.Equal(t, "project1", result.ProjectId)
		assert.Equal(t, []string{"read", "write"}, result.Scopes)
		assert.True(t, result.Enabled)
		assert.Equal(t, &createdAt, result.CreatedAt)
		assert.Equal(t, "user1", result.CreatedBy)
		assert.Equal(t, &updatedAt, result.UpdatedAt)
		assert.Equal(t, "user2", result.UpdatedBy)
	})
	
	t.Run("nil_client", func(t *testing.T) {
		var modelClient *models.Client = nil
		
		// This should panic or cause issues, but let's test the function behavior
		// In a real scenario, this would be handled by validation before calling the function
		assert.Panics(t, func() {
			fromModelToSdk(modelClient)
		})
	})
	
	t.Run("empty_client", func(t *testing.T) {
		modelClient := &models.Client{}
		
		result := fromModelToSdk(modelClient)
		
		assert.NotNil(t, result)
		assert.Empty(t, result.Id)
		assert.Empty(t, result.Name)
		assert.Empty(t, result.Description)
		assert.Empty(t, result.Secret)
		assert.Nil(t, result.Tags)
		assert.Nil(t, result.RedirectURLs)
		assert.Empty(t, result.DefaultAuthProviderId)
		assert.False(t, result.GoIamClient)
		assert.Empty(t, result.ProjectId)
		assert.Nil(t, result.Scopes)
		assert.False(t, result.Enabled)
		assert.Nil(t, result.CreatedAt)
		assert.Empty(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedAt)
		assert.Empty(t, result.UpdatedBy)
	})
}

func TestFromSdkToModel(t *testing.T) {
	t.Run("complete_client", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(time.Hour)
		
		sdkClient := sdk.Client{
			Id:                    "client1",
			Name:                  "Test Client",
			Description:           "A test client",
			Secret:                "secret123",
			Tags:                  []string{"tag1", "tag2"},
			RedirectURLs:          []string{"https://example.com/callback"},
			DefaultAuthProviderId: "provider1",
			GoIamClient:           true,
			ProjectId:             "project1",
			Scopes:                []string{"read", "write"},
			Enabled:               true,
			CreatedAt:             &createdAt,
			CreatedBy:             "user1",
			UpdatedAt:             &updatedAt,
			UpdatedBy:             "user2",
		}
		
		result := fromSdkToModel(sdkClient)
		
		assert.Equal(t, "client1", result.Id)
		assert.Equal(t, "Test Client", result.Name)
		assert.Equal(t, "A test client", result.Description)
		assert.Equal(t, "secret123", result.Secret)
		assert.Equal(t, []string{"tag1", "tag2"}, result.Tags)
		assert.Equal(t, []string{"https://example.com/callback"}, result.RedirectURLs)
		assert.Equal(t, "provider1", result.DefaultAuthProviderId)
		assert.True(t, result.GoIamClient)
		assert.Equal(t, "project1", result.ProjectId)
		assert.Equal(t, []string{"read", "write"}, result.Scopes)
		assert.True(t, result.Enabled)
		assert.Equal(t, &createdAt, result.CreatedAt)
		assert.Equal(t, "user1", result.CreatedBy)
		assert.Equal(t, &updatedAt, result.UpdatedAt)
		assert.Equal(t, "user2", result.UpdatedBy)
	})
	
	t.Run("empty_client", func(t *testing.T) {
		sdkClient := sdk.Client{}
		
		result := fromSdkToModel(sdkClient)
		
		assert.Empty(t, result.Id)
		assert.Empty(t, result.Name)
		assert.Empty(t, result.Description)
		assert.Empty(t, result.Secret)
		assert.Nil(t, result.Tags)
		assert.Nil(t, result.RedirectURLs)
		assert.Empty(t, result.DefaultAuthProviderId)
		assert.False(t, result.GoIamClient)
		assert.Empty(t, result.ProjectId)
		assert.Nil(t, result.Scopes)
		assert.False(t, result.Enabled)
		assert.Nil(t, result.CreatedAt)
		assert.Empty(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedAt)
		assert.Empty(t, result.UpdatedBy)
	})
}

func TestHashSecret(t *testing.T) {
	t.Run("hash_simple_secret", func(t *testing.T) {
		secret := "mysecret123"
		
		hashedSecret, err := hashSecret(secret)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedSecret)
		assert.NotEqual(t, secret, hashedSecret)
		// SHA256 hash encoded in base64 should be 44 characters long
		assert.Equal(t, 44, len(hashedSecret))
	})
	
	t.Run("hash_empty_secret", func(t *testing.T) {
		secret := ""
		
		hashedSecret, err := hashSecret(secret)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedSecret)
		assert.Equal(t, 44, len(hashedSecret))
	})
	
	t.Run("hash_same_secret_produces_same_hash", func(t *testing.T) {
		secret := "consistent_secret"
		
		hash1, err1 := hashSecret(secret)
		hash2, err2 := hashSecret(secret)
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, hash1, hash2)
	})
	
	t.Run("hash_different_secrets_produce_different_hashes", func(t *testing.T) {
		secret1 := "secret1"
		secret2 := "secret2"
		
		hash1, err1 := hashSecret(secret1)
		hash2, err2 := hashSecret(secret2)
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})
	
	t.Run("hash_special_characters", func(t *testing.T) {
		secret := "special!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/~`"
		
		hashedSecret, err := hashSecret(secret)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedSecret)
		assert.Equal(t, 44, len(hashedSecret))
	})
}

func TestGenerateRandomSecret(t *testing.T) {
	t.Run("generate_secret_with_specific_length", func(t *testing.T) {
		length := 32
		
		secret, err := generateRandomSecret(length)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		assert.Equal(t, length, len(secret))
	})
	
	t.Run("generate_secret_with_different_lengths", func(t *testing.T) {
		lengths := []int{8, 16, 24, 32, 64}
		
		for _, length := range lengths {
			secret, err := generateRandomSecret(length)
			
			assert.NoError(t, err, "Failed for length %d", length)
			assert.NotEmpty(t, secret, "Secret empty for length %d", length)
			assert.Equal(t, length, len(secret), "Wrong length for %d", length)
		}
	})
	
	t.Run("generate_multiple_secrets_are_different", func(t *testing.T) {
		length := 32
		secrets := make([]string, 10)
		
		for i := 0; i < 10; i++ {
			secret, err := generateRandomSecret(length)
			assert.NoError(t, err)
			secrets[i] = secret
		}
		
		// Check that all secrets are different
		for i := 0; i < len(secrets); i++ {
			for j := i + 1; j < len(secrets); j++ {
				assert.NotEqual(t, secrets[i], secrets[j], "Secrets %d and %d are identical", i, j)
			}
		}
	})
	
	t.Run("generate_secret_with_zero_length", func(t *testing.T) {
		length := 0
		
		secret, err := generateRandomSecret(length)
		
		assert.NoError(t, err)
		assert.Empty(t, secret)
		assert.Equal(t, 0, len(secret))
	})
	
	t.Run("generate_secret_with_small_length", func(t *testing.T) {
		length := 1
		
		secret, err := generateRandomSecret(length)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		assert.Equal(t, length, len(secret))
	})
}
