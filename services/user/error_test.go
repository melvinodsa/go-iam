package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUserErrors tests the custom error variables
func TestUserErrors(t *testing.T) {
	// Test ErrorUserNotFound
	assert.NotNil(t, ErrorUserNotFound)
	assert.Equal(t, "user not found", ErrorUserNotFound.Error())
	assert.True(t, errors.Is(ErrorUserNotFound, ErrorUserNotFound))
}

// TestErrorTypes tests that the error is of type error
func TestErrorTypes(t *testing.T) {

	err := ErrorUserNotFound
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), ErrorUserNotFound)
}
