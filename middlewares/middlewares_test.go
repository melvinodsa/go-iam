package middlewares

import (
	"context"
	"fmt"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func createTestUser() *sdk.User {
	return &sdk.User{
		Id:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
	}
}

func createTestUser2() *sdk.User {
	return &sdk.User{
		Id:    "test-user-id-2",
		Name:  "Test User 2",
		Email: "test2@example.com",
	}
}

func createTestMetadata() sdk.Metadata {
	return sdk.Metadata{
		User:       createTestUser(),
		ProjectIds: []string{"project1", "project2", "project3"},
	}
}

func TestGetProjects_WithProjects(t *testing.T) {
	// Create context with projects
	projectList := []string{"project1", "project2", "project3"}
	ctx := context.WithValue(context.Background(), projects, projectList)

	result := GetProjects(ctx)

	assert.Equal(t, projectList, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "project1", result[0])
	assert.Equal(t, "project2", result[1])
	assert.Equal(t, "project3", result[2])
}

func TestGetProjects_EmptyProjects(t *testing.T) {
	// Create context with empty projects
	projectList := []string{}
	ctx := context.WithValue(context.Background(), projects, projectList)

	result := GetProjects(ctx)

	assert.Equal(t, projectList, result)
	assert.Len(t, result, 0)
}

func TestGetProjects_SingleProject(t *testing.T) {
	// Create context with single project
	projectList := []string{"single-project"}
	ctx := context.WithValue(context.Background(), projects, projectList)

	result := GetProjects(ctx)

	assert.Equal(t, projectList, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "single-project", result[0])
}

func TestGetProjects_PanicsOnWrongType(t *testing.T) {
	// Create context with wrong type (should panic on type assertion)
	ctx := context.WithValue(context.Background(), projects, "not-a-slice")

	assert.Panics(t, func() {
		GetProjects(ctx)
	})
}

func TestGetProjects_PanicsOnMissingValue(t *testing.T) {
	// Create context without projects value (should panic on type assertion)
	ctx := context.Background()

	assert.Panics(t, func() {
		GetProjects(ctx)
	})
}

func TestGetUser_WithUser(t *testing.T) {
	// Create context with user
	testUser := createTestUser()
	ctx := context.WithValue(context.Background(), userValue, testUser)

	result := GetUser(ctx)

	assert.NotNil(t, result)
	assert.Equal(t, testUser, result)
	assert.Equal(t, "test-user-id", result.Id)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestGetUser_NoUser(t *testing.T) {
	// Create context without user
	ctx := context.Background()

	result := GetUser(ctx)

	assert.Nil(t, result)
}

func TestGetUser_NilUser(t *testing.T) {
	// Create context with explicit nil user
	ctx := context.WithValue(context.Background(), userValue, nil)

	result := GetUser(ctx)

	assert.Nil(t, result)
}

func TestGetUser_WrongType(t *testing.T) {
	// Create context with wrong type for user
	ctx := context.WithValue(context.Background(), userValue, "not-a-user")

	result := GetUser(ctx)

	assert.Nil(t, result)
}

func TestGetUser_EmptyUser(t *testing.T) {
	// Create context with empty user struct
	emptyUser := &sdk.User{}
	ctx := context.WithValue(context.Background(), userValue, emptyUser)

	result := GetUser(ctx)

	assert.NotNil(t, result)
	assert.Equal(t, emptyUser, result)
	assert.Equal(t, "", result.Id)
	assert.Equal(t, "", result.Name)
	assert.Equal(t, "", result.Email)
}

func TestGetMetadata_WithFullMetadata(t *testing.T) {
	// Create context with both user and projects
	testUser := createTestUser()
	projectList := []string{"project1", "project2"}
	ctx := context.WithValue(
		context.WithValue(context.Background(), userValue, testUser),
		projects, projectList,
	)

	result := GetMetadata(ctx)

	assert.NotNil(t, result)
	assert.Equal(t, testUser, result.User)
	assert.Equal(t, projectList, result.ProjectIds)
	assert.Len(t, result.ProjectIds, 2)
}

func TestGetMetadata_WithUserOnly(t *testing.T) {
	// Create context with user but without projects (should panic on GetProjects)
	testUser := createTestUser()
	ctx := context.WithValue(context.Background(), userValue, testUser)

	assert.Panics(t, func() {
		GetMetadata(ctx)
	})
}

func TestGetMetadata_WithProjectsOnly(t *testing.T) {
	// Create context with projects but without user
	projectList := []string{"project1", "project2"}
	ctx := context.WithValue(context.Background(), projects, projectList)

	result := GetMetadata(ctx)

	assert.NotNil(t, result)
	assert.Nil(t, result.User)
	assert.Equal(t, projectList, result.ProjectIds)
}

func TestGetMetadata_EmptyContext(t *testing.T) {
	// Create empty context (should panic on GetProjects)
	ctx := context.Background()

	assert.Panics(t, func() {
		GetMetadata(ctx)
	})
}

func TestGetMetadata_EmptyValues(t *testing.T) {
	// Create context with empty user and projects
	projectList := []string{}
	ctx := context.WithValue(
		context.WithValue(context.Background(), userValue, nil),
		projects, projectList,
	)

	result := GetMetadata(ctx)

	assert.NotNil(t, result)
	assert.Nil(t, result.User)
	assert.Equal(t, projectList, result.ProjectIds)
	assert.Len(t, result.ProjectIds, 0)
}

func TestAddMetadata_WithFullMetadata(t *testing.T) {
	// Test adding complete metadata
	testMetadata := createTestMetadata()
	baseCtx := context.Background()

	resultCtx := AddMetadata(baseCtx, testMetadata)

	// Verify the context contains the expected values
	retrievedUser := GetUser(resultCtx)
	retrievedProjects := GetProjects(resultCtx)

	assert.Equal(t, testMetadata.User, retrievedUser)
	assert.Equal(t, testMetadata.ProjectIds, retrievedProjects)
}

func TestAddMetadata_WithNilUser(t *testing.T) {
	// Test adding metadata with nil user
	testMetadata := sdk.Metadata{
		User:       nil,
		ProjectIds: []string{"project1", "project2"},
	}
	baseCtx := context.Background()

	resultCtx := AddMetadata(baseCtx, testMetadata)

	// Verify the context contains the expected values
	retrievedUser := GetUser(resultCtx)
	retrievedProjects := GetProjects(resultCtx)

	assert.Nil(t, retrievedUser)
	assert.Equal(t, testMetadata.ProjectIds, retrievedProjects)
}

func TestAddMetadata_WithEmptyProjects(t *testing.T) {
	// Test adding metadata with empty projects
	testMetadata := sdk.Metadata{
		User:       createTestUser(),
		ProjectIds: []string{},
	}
	baseCtx := context.Background()

	resultCtx := AddMetadata(baseCtx, testMetadata)

	// Verify the context contains the expected values
	retrievedUser := GetUser(resultCtx)
	retrievedProjects := GetProjects(resultCtx)

	assert.Equal(t, testMetadata.User, retrievedUser)
	assert.Equal(t, testMetadata.ProjectIds, retrievedProjects)
	assert.Len(t, retrievedProjects, 0)
}

func TestAddMetadata_EmptyMetadata(t *testing.T) {
	// Test adding completely empty metadata
	testMetadata := sdk.Metadata{}
	baseCtx := context.Background()

	resultCtx := AddMetadata(baseCtx, testMetadata)

	// Verify the context contains the expected values
	retrievedUser := GetUser(resultCtx)
	retrievedProjects := GetProjects(resultCtx)

	assert.Nil(t, retrievedUser)
	assert.Nil(t, retrievedProjects) // This should be nil, not empty slice
}

func TestAddMetadata_OverwriteExisting(t *testing.T) {
	// Test overwriting existing context values
	originalUser := &sdk.User{Id: "original-user", Name: "Original User"}
	originalProjects := []string{"original-project"}

	// Create context with existing values
	existingCtx := context.WithValue(
		context.WithValue(context.Background(), userValue, originalUser),
		projects, originalProjects,
	)

	// Add new metadata that should overwrite
	newMetadata := createTestMetadata()
	resultCtx := AddMetadata(existingCtx, newMetadata)

	// Verify the new values overwrote the old ones
	retrievedUser := GetUser(resultCtx)
	retrievedProjects := GetProjects(resultCtx)

	assert.Equal(t, newMetadata.User, retrievedUser)
	assert.Equal(t, newMetadata.ProjectIds, retrievedProjects)
	assert.NotEqual(t, originalUser, retrievedUser)
	assert.NotEqual(t, originalProjects, retrievedProjects)
}

func TestAddMetadata_ContextChaining(t *testing.T) {
	// Test that context chaining works correctly
	baseCtx := context.Background()

	// Add metadata
	testMetadata := createTestMetadata()
	resultCtx := AddMetadata(baseCtx, testMetadata)

	// Verify we can retrieve the full metadata back
	retrievedMetadata := GetMetadata(resultCtx)

	assert.Equal(t, testMetadata.User.Id, retrievedMetadata.User.Id)
	assert.Equal(t, testMetadata.User.Name, retrievedMetadata.User.Name)
	assert.Equal(t, testMetadata.User.Email, retrievedMetadata.User.Email)
	assert.Equal(t, testMetadata.ProjectIds, retrievedMetadata.ProjectIds)
}

func TestContextKeys_AreUnique(t *testing.T) {
	// Test that our context keys are properly isolated
	ctx := context.Background()

	// Add some values using our keys
	ctx = context.WithValue(ctx, projects, []string{"test-project"})
	ctx = context.WithValue(ctx, userValue, createTestUser())

	// Try to add conflicting values with different keys
	ctx = context.WithValue(ctx, projects, []string{"conflicting-value"})
	ctx = context.WithValue(ctx, userValue, createTestUser2())

	// Verify our functions still work correctly
	retrievedProjects := GetProjects(ctx)
	retrievedUser := GetUser(ctx)

	assert.NotEqual(t, []string{"test-project"}, retrievedProjects)

	assert.Equal(t, []string{"conflicting-value"}, retrievedProjects)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, "test-user-id-2", retrievedUser.Id)
}

// Benchmark tests
func BenchmarkGetProjects(b *testing.B) {
	projectList := []string{"project1", "project2", "project3", "project4", "project5"}
	ctx := context.WithValue(context.Background(), projects, projectList)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetProjects(ctx)
	}
}

func BenchmarkGetUser(b *testing.B) {
	testUser := createTestUser()
	ctx := context.WithValue(context.Background(), userValue, testUser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUser(ctx)
	}
}

func BenchmarkGetMetadata(b *testing.B) {
	testUser := createTestUser()
	projectList := []string{"project1", "project2", "project3"}
	ctx := context.WithValue(
		context.WithValue(context.Background(), userValue, testUser),
		projects, projectList,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetMetadata(ctx)
	}
}

func BenchmarkAddMetadata(b *testing.B) {
	testMetadata := createTestMetadata()
	baseCtx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddMetadata(baseCtx, testMetadata)
	}
}

// Edge case tests
func TestComplexScenarios(t *testing.T) {
	t.Run("Multiple context operations", func(t *testing.T) {
		// Simulate a complex scenario with multiple context operations
		ctx := context.Background()

		// Add initial metadata
		metadata1 := sdk.Metadata{
			User:       &sdk.User{Id: "user1", Name: "User 1"},
			ProjectIds: []string{"proj1"},
		}
		ctx = AddMetadata(ctx, metadata1)

		// Verify initial state
		assert.Equal(t, "user1", GetUser(ctx).Id)
		assert.Equal(t, []string{"proj1"}, GetProjects(ctx))

		// Update with new metadata
		metadata2 := sdk.Metadata{
			User:       &sdk.User{Id: "user2", Name: "User 2"},
			ProjectIds: []string{"proj1", "proj2"},
		}
		ctx = AddMetadata(ctx, metadata2)

		// Verify updated state
		assert.Equal(t, "user2", GetUser(ctx).Id)
		assert.Equal(t, []string{"proj1", "proj2"}, GetProjects(ctx))

		// Get full metadata and verify consistency
		finalMetadata := GetMetadata(ctx)
		assert.Equal(t, "user2", finalMetadata.User.Id)
		assert.Equal(t, []string{"proj1", "proj2"}, finalMetadata.ProjectIds)
	})

	t.Run("Large project lists", func(t *testing.T) {
		// Test with large number of projects
		largeProjectList := make([]string, 100)
		for i := 0; i < 100; i++ {
			largeProjectList[i] = fmt.Sprintf("project-%d", i)
		}

		metadata := sdk.Metadata{
			User:       createTestUser(),
			ProjectIds: largeProjectList,
		}

		ctx := AddMetadata(context.Background(), metadata)
		retrievedProjects := GetProjects(ctx)

		assert.Len(t, retrievedProjects, 100)
		assert.Equal(t, "project-0", retrievedProjects[0])
		assert.Equal(t, "project-99", retrievedProjects[99])
	})
}
