package services

import (
	"github.com/stretchr/testify/mock"
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
