// Package middlewares provides utility functions for extracting and managing
// authentication and authorization metadata from request contexts.
// It handles user information, project access, and metadata operations
// used throughout the Go IAM system.
package middlewares

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

// GetProjects extracts the list of project IDs from the given context.
// This function retrieves project IDs that were previously stored in the context
// by the projects middleware, typically from the X-Project-Ids header.
//
// Parameters:
//   - ctx: The context containing project information
//
// Returns:
//   - []string: Slice of project IDs if found, nil if not present or invalid type
func GetProjects(ctx context.Context) []string {
	val, ok := ctx.Value(sdk.ProjectsTypeVal).([]string)
	if !ok {
		return nil
	}
	return val
}

// GetUser extracts the authenticated user from the given context.
// This function retrieves user information that was previously stored in the context
// by the authentication middleware after successful token validation.
//
// Parameters:
//   - ctx: The context containing user information
//
// Returns:
//   - *sdk.User: Pointer to the authenticated user if found, nil if not present or invalid type
func GetUser(ctx context.Context) *sdk.User {
	user := ctx.Value(sdk.UserTypeVal)
	if user == nil {
		return nil
	}
	authUser, ok := user.(*sdk.User)
	if !ok {
		return nil
	}
	return authUser
}

// GetMetadata creates a consolidated metadata object containing both user and project information
// from the context. This function combines the results of GetUser and GetProjects for
// convenient access to all authentication and authorization metadata.
//
// Parameters:
//   - ctx: The context containing user and project information
//
// Returns:
//   - sdk.Metadata: Metadata object containing user and project IDs
func GetMetadata(ctx context.Context) sdk.Metadata {
	return sdk.Metadata{
		User:       GetUser(ctx),
		ProjectIds: GetProjects(ctx),
	}
}

// AddMetadata stores user and project information into the context.
// This function is typically used by services or handlers to inject metadata
// into the context for downstream processing.
//
// Parameters:
//   - ctx: The context to store metadata in
//   - metadata: The metadata containing user and project information to store
//
// Returns:
//   - context.Context: New context with the metadata stored
func AddMetadata(ctx context.Context, metadata sdk.Metadata) context.Context {
	return context.WithValue(context.WithValue(ctx, sdk.ProjectsTypeVal, metadata.ProjectIds), sdk.UserTypeVal, metadata.User)
}
