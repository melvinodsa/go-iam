package authprovider

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockEncryptService implements encrypt.Service interface for testing
type MockEncryptService struct {
	mock.Mock
}

func (m *MockEncryptService) Encrypt(rawMessage string) (string, error) {
	args := m.Called(rawMessage)
	return args.String(0), args.Error(1)
}

func (m *MockEncryptService) Decrypt(encryptedMessage string) (string, error) {
	args := m.Called(encryptedMessage)
	return args.String(0), args.Error(1)
}

func TestStore_decryptSecrets(t *testing.T) {
	tests := []struct {
		name          string
		provider      *models.AuthProvider
		encSetup      func(*MockEncryptService)
		expectedError error
		validateCall  func(*testing.T, *models.AuthProvider)
	}{
		{
			name: "success_decrypt_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client ID",
						Value:    "client123",
						Key:      "@GOOGLE/CLIENT_ID",
						IsSecret: false,
					},
					{
						Label:    "Client Secret",
						Value:    "encrypted_secret",
						Key:      "@GOOGLE/CLIENT_SECRET",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Decrypt", "encrypted_secret").Return("decrypted_secret", nil)
			},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "client123", provider.Params[0].Value)
				assert.Equal(t, "decrypted_secret", provider.Params[1].Value)
			},
		},
		{
			name: "success_no_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client ID",
						Value:    "client123",
						Key:      "@GOOGLE/CLIENT_ID",
						IsSecret: false,
					},
				},
			},
			encSetup:      func(m *MockEncryptService) {},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "client123", provider.Params[0].Value)
			},
		},
		{
			name: "success_empty_params",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{},
			},
			encSetup:      func(m *MockEncryptService) {},
			expectedError: nil,
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
		{
			name: "error_decryption_failure",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client Secret",
						Value:    "encrypted_secret",
						Key:      "@GOOGLE/CLIENT_SECRET",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Decrypt", "encrypted_secret").Return("", errors.New("decryption failed"))
			},
			expectedError: errors.New("error decrypting auth provider secret at 0 : decryption failed"),
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
		{
			name: "success_multiple_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Secret 1",
						Value:    "encrypted_secret1",
						Key:      "@GOOGLE/SECRET1",
						IsSecret: true,
					},
					{
						Label:    "Secret 2",
						Value:    "encrypted_secret2",
						Key:      "@GOOGLE/SECRET2",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Decrypt", "encrypted_secret1").Return("decrypted_secret1", nil)
				m.On("Decrypt", "encrypted_secret2").Return("decrypted_secret2", nil)
			},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "decrypted_secret1", provider.Params[0].Value)
				assert.Equal(t, "decrypted_secret2", provider.Params[1].Value)
			},
		},
		{
			name: "error_decryption_failure_second_param",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Secret 1",
						Value:    "encrypted_secret1",
						Key:      "@GOOGLE/SECRET1",
						IsSecret: true,
					},
					{
						Label:    "Secret 2",
						Value:    "encrypted_secret2",
						Key:      "@GOOGLE/SECRET2",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Decrypt", "encrypted_secret1").Return("decrypted_secret1", nil)
				m.On("Decrypt", "encrypted_secret2").Return("", errors.New("decryption failed"))
			},
			expectedError: errors.New("error decrypting auth provider secret at 1 : decryption failed"),
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEnc := &MockEncryptService{}
			store := store{enc: mockEnc}

			tt.encSetup(mockEnc)

			err := store.decryptSecrets(tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				tt.validateCall(t, tt.provider)
			}

			mockEnc.AssertExpectations(t)
		})
	}
}

func TestStore_encryptSecrets(t *testing.T) {
	tests := []struct {
		name          string
		provider      *models.AuthProvider
		encSetup      func(*MockEncryptService)
		expectedError error
		validateCall  func(*testing.T, *models.AuthProvider)
	}{
		{
			name: "success_encrypt_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client ID",
						Value:    "client123",
						Key:      "@GOOGLE/CLIENT_ID",
						IsSecret: false,
					},
					{
						Label:    "Client Secret",
						Value:    "raw_secret",
						Key:      "@GOOGLE/CLIENT_SECRET",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Encrypt", "raw_secret").Return("encrypted_secret", nil)
			},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "client123", provider.Params[0].Value)
				assert.Equal(t, "encrypted_secret", provider.Params[1].Value)
			},
		},
		{
			name: "success_no_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client ID",
						Value:    "client123",
						Key:      "@GOOGLE/CLIENT_ID",
						IsSecret: false,
					},
				},
			},
			encSetup:      func(m *MockEncryptService) {},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "client123", provider.Params[0].Value)
			},
		},
		{
			name: "success_empty_params",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{},
			},
			encSetup:      func(m *MockEncryptService) {},
			expectedError: nil,
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
		{
			name: "error_encryption_failure",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Client Secret",
						Value:    "raw_secret",
						Key:      "@GOOGLE/CLIENT_SECRET",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Encrypt", "raw_secret").Return("", errors.New("encryption failed"))
			},
			expectedError: errors.New("error encrypting auth provider secret at 0 : encryption failed"),
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
		{
			name: "success_multiple_secrets",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Secret 1",
						Value:    "raw_secret1",
						Key:      "@GOOGLE/SECRET1",
						IsSecret: true,
					},
					{
						Label:    "Public Param",
						Value:    "public_value",
						Key:      "@GOOGLE/PUBLIC",
						IsSecret: false,
					},
					{
						Label:    "Secret 2",
						Value:    "raw_secret2",
						Key:      "@GOOGLE/SECRET2",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Encrypt", "raw_secret1").Return("encrypted_secret1", nil)
				m.On("Encrypt", "raw_secret2").Return("encrypted_secret2", nil)
			},
			expectedError: nil,
			validateCall: func(t *testing.T, provider *models.AuthProvider) {
				assert.Equal(t, "encrypted_secret1", provider.Params[0].Value)
				assert.Equal(t, "public_value", provider.Params[1].Value)
				assert.Equal(t, "encrypted_secret2", provider.Params[2].Value)
			},
		},
		{
			name: "error_encryption_failure_first_param",
			provider: &models.AuthProvider{
				Params: []models.AuthProviderParam{
					{
						Label:    "Secret 1",
						Value:    "raw_secret1",
						Key:      "@GOOGLE/SECRET1",
						IsSecret: true,
					},
					{
						Label:    "Secret 2",
						Value:    "raw_secret2",
						Key:      "@GOOGLE/SECRET2",
						IsSecret: true,
					},
				},
			},
			encSetup: func(m *MockEncryptService) {
				m.On("Encrypt", "raw_secret1").Return("", errors.New("encryption failed"))
			},
			expectedError: errors.New("error encrypting auth provider secret at 0 : encryption failed"),
			validateCall:  func(t *testing.T, provider *models.AuthProvider) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEnc := &MockEncryptService{}
			store := store{enc: mockEnc}

			tt.encSetup(mockEnc)

			err := store.encryptSecrets(tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				tt.validateCall(t, tt.provider)
			}

			mockEnc.AssertExpectations(t)
		})
	}
}

func createTestAuthProvider(id, name, projectId string, providers []models.AuthProviderParam) models.AuthProvider {
	return models.AuthProvider{
		Id:     id,
		Name:   name,
		Params: providers,
	}
}

func createTestAuthProviderSdk(id, name, projectId string, providers []models.AuthProviderParam) sdk.AuthProvider {
	return sdk.AuthProvider{
		Id:     id,
		Name:   name,
		Params: fromProviderParamsModelToSdk(providers),
	}
}

func TestStore_Get_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record := createTestAuthProvider("1", "test", "test-project", []models.AuthProviderParam{})

	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{
		{Key: md.IdKey, Value: record.Id},
		{Key: md.NameKey, Value: record.Name},
		{Key: md.ProjectIdKey, Value: record.ProjectId},
		{Key: md.ParamsKey, Value: record.Params},
	}, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: record.Id}}, mock.Anything).Return(mockResult)

	mockEnc := &MockEncryptService{}
	store := NewStore(mockEnc, mockDB)

	// Execute
	result, err := store.Get(ctx, record.Id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, record.Id, result.Id)
	assert.Equal(t, record.Name, result.Name)
	assert.Equal(t, record.ProjectId, result.ProjectId)
	assert.Equal(t, fromModelToSdk(&record).Params, result.Params)
}

func TestStore_Get_NotFound(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	// Mock FindOne to return no result
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: "unknown"}}, mock.Anything).Return(mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil))

	mockEnc := &MockEncryptService{}
	store := NewStore(mockEnc, mockDB)

	// Execute
	result, err := store.Get(ctx, "unknown")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStore_Get_Error(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	// Mock FindOne to return an error
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: "unknown"}}, mock.Anything).Return(mongo.NewSingleResultFromDocument(bson.D{}, errors.New("custom error"), nil))

	mockEnc := &MockEncryptService{}
	store := NewStore(mockEnc, mockDB)

	// Execute
	result, err := store.Get(ctx, "unknown")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStore_Get_DecryptError(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record := createTestAuthProvider("1", "test", "test-project", []models.AuthProviderParam{
		{Key: "secret1", Value: "encrypted_secret1", IsSecret: true},
	})

	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{
		{Key: md.IdKey, Value: record.Id},
		{Key: md.NameKey, Value: record.Name},
		{Key: md.ProjectIdKey, Value: record.ProjectId},
		{Key: md.ParamsKey, Value: record.Params},
	}, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: record.Id}}, mock.Anything).Return(mockResult)

	mockEnc := &MockEncryptService{}
	mockEnc.On("Decrypt", "encrypted_secret1").Return("", errors.New("decryption failed"))
	store := NewStore(mockEnc, mockDB)

	// Execute
	result, err := store.Get(ctx, record.Id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStore_GetAll_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	// Create test data
	record1 := createTestAuthProvider("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_id", Value: "client1", IsSecret: false},
		{Key: "client_secret", Value: "encrypted_secret1", IsSecret: true},
	})
	record2 := createTestAuthProvider("2", "Provider 2", "project1", []models.AuthProviderParam{
		{Key: "api_key", Value: "key123", IsSecret: false},
	})
	record3 := createTestAuthProvider("3", "Provider 3", "project2", []models.AuthProviderParam{})

	// Create cursor from documents
	providers := []models.AuthProvider{record1, record2, record3}
	documents := make([]interface{}, len(providers))
	for i, p := range providers {
		documents[i] = p
	}
	cursor, _ := mongo.NewCursorFromDocuments(documents, nil, nil)

	// Mock Find to return cursor
	filter := bson.D{}
	filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: []string{"project1", "project2"}}}})
	mockDB.On("Find", ctx, md, filter, mock.Anything).Return(cursor, nil)

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Decrypt", "encrypted_secret1").Return("decrypted_secret1", nil)
	store := NewStore(mockEnc, mockDB)

	// Execute
	params := sdk.AuthProviderQueryParams{
		ProjectIds: []string{"project1", "project2"},
	}
	result, err := store.GetAll(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Validate first provider
	assert.Equal(t, "1", result[0].Id)
	assert.Equal(t, "Provider 1", result[0].Name)
	assert.Len(t, result[0].Params, 2)
	assert.Equal(t, "client_id", result[0].Params[0].Key)
	assert.Equal(t, "client1", result[0].Params[0].Value)
	assert.False(t, result[0].Params[0].IsSecret)
	assert.Equal(t, "client_secret", result[0].Params[1].Key)
	assert.Equal(t, "decrypted_secret1", result[0].Params[1].Value) // Should be decrypted
	assert.True(t, result[0].Params[1].IsSecret)

	// Validate second provider
	assert.Equal(t, "2", result[1].Id)
	assert.Equal(t, "Provider 2", result[1].Name)
	assert.Len(t, result[1].Params, 1)

	// Validate third provider
	assert.Equal(t, "3", result[2].Id)
	assert.Equal(t, "Provider 3", result[2].Name)
	assert.Empty(t, result[2].Params)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_GetAll_FindError(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	// Mock Find to return an error
	filter := bson.D{}
	filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: []string{"project1", "project2"}}}})
	mockDB.On("Find", ctx, md, filter, mock.Anything).Return(nil, errors.New("find error"))

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	store := NewStore(mockEnc, mockDB)

	// Execute
	params := sdk.AuthProviderQueryParams{
		ProjectIds: []string{"project1", "project2"},
	}
	result, err := store.GetAll(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_GetAll_DecryptError(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record1 := createTestAuthProvider("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_id", Value: "client1", IsSecret: false},
		{Key: "client_secret", Value: "encrypted_secret1", IsSecret: true},
	})

	cursor, _ := mongo.NewCursorFromDocuments([]interface{}{record1}, nil, nil)

	// Mock Find to return an error
	filter := bson.D{}
	filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: bson.D{{Key: "$in", Value: []string{"project1", "project2"}}}})
	mockDB.On("Find", ctx, md, filter, mock.Anything).Return(cursor, nil)

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Decrypt", "encrypted_secret1").Return("", errors.New("decryption failed"))
	store := NewStore(mockEnc, mockDB)

	// Execute
	params := sdk.AuthProviderQueryParams{
		ProjectIds: []string{"project1", "project2"},
	}
	result, err := store.GetAll(ctx, params)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_Create_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	// md := models.GetAuthProviderModel()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock InsertOne to return no error
	mockDB.On("InsertOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "1"}, nil)

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("encrypted_secret1", nil)
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Create(ctx, &record)

	// Assert
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_Create_Error(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock InsertOne to return an error
	mockDB.On("InsertOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "1"}, errors.New("insert error"))

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("encrypted_secret1", nil)
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Create(ctx, &record)

	// Assert
	assert.Error(t, err)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_Create_EncryptError(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock InsertOne to return no error
	mockDB.On("InsertOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "1"}, nil)

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("", errors.New("encryption failed"))
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Create(ctx, &record)

	// Assert
	assert.Error(t, err)

	mockEnc.AssertExpectations(t)
}

func TestStore_Update_Success(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock UpdateOne to return no error
	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{
		{Key: md.IdKey, Value: record.Id},
		{Key: md.NameKey, Value: record.Name},
		{Key: md.ProjectIdKey, Value: record.ProjectId},
		{Key: md.ParamsKey, Value: record.Params},
	}, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: record.Id}}, mock.Anything).Return(mockResult)
	mockDB.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1}, nil)

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("encrypted_secret1", nil)
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Update(ctx, &record)

	// Assert
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_Update_Error(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{
		{Key: md.IdKey, Value: record.Id},
		{Key: md.NameKey, Value: record.Name},
		{Key: md.ProjectIdKey, Value: record.ProjectId},
		{Key: md.ParamsKey, Value: record.Params},
	}, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: record.Id}}, mock.Anything).Return(mockResult)

	// Mock UpdateOne to return an error
	mockDB.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 0}, errors.New("update error"))

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("encrypted_secret1", nil)
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Update(ctx, &record)

	// Assert
	assert.Error(t, err)

	mockDB.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
}

func TestStore_Update_EncryptError(t *testing.T) {
	mockDB := test.SetupMockDB()
	ctx := context.Background()
	md := models.GetAuthProviderModel()

	record := createTestAuthProviderSdk("1", "Provider 1", "project1", []models.AuthProviderParam{
		{Key: "client_secret", Value: "secret", IsSecret: true},
	})

	// Mock FindOne to return existing migration (no error)
	mockResult := mongo.NewSingleResultFromDocument(bson.D{
		{Key: md.IdKey, Value: record.Id},
		{Key: md.NameKey, Value: record.Name},
		{Key: md.ProjectIdKey, Value: record.ProjectId},
		{Key: md.ParamsKey, Value: record.Params},
	}, nil, nil)
	mockDB.On("FindOne", ctx, md, bson.D{{Key: md.IdKey, Value: record.Id}}, mock.Anything).Return(mockResult)

	// Mock UpdateOne to return an error
	mockDB.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 0}, errors.New("update error"))

	// Mock encryption service
	mockEnc := &MockEncryptService{}
	mockEnc.On("Encrypt", "secret").Return("", errors.New("encryption failed"))
	store := NewStore(mockEnc, mockDB)

	// Execute
	err := store.Update(ctx, &record)

	// Assert
	assert.Error(t, err)

	mockEnc.AssertExpectations(t)
}
