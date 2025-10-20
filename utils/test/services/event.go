package services

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/mock"
)

// MockEvent implements utils.Event interface for testing
type MockEvent[T any] struct {
	mock.Mock
}

func (m *MockEvent[T]) Name() goiamuniverse.Event {
	args := m.Called()
	return args.Get(0).(goiamuniverse.Event)
}

func (m *MockEvent[T]) Payload() T {
	args := m.Called()
	return args.Get(0).(T)
}

func (m *MockEvent[T]) Metadata() sdk.Metadata {
	args := m.Called()
	return args.Get(0).(sdk.Metadata)
}

func (m *MockEvent[T]) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}
