package services

import (
	"github.com/stretchr/testify/mock"
)

// MockJWTService implements jwt.Service interface for testing
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(claims map[string]interface{}, expiryTimeInSeconds int64) (string, error) {
	args := m.Called(claims, expiryTimeInSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
