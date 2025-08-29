package middlewares

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type projectType struct{}

var projects = projectType{}

type userType struct{}

var userValue = userType{}

func GetProjects(ctx context.Context) []string {
	return ctx.Value(projects).([]string)
}

func GetUser(ctx context.Context) *sdk.User {
	user := ctx.Value(userValue)
	if user == nil {
		return nil
	}
	authUser, ok := user.(*sdk.User)
	if !ok {
		return nil
	}
	return authUser
}

func GetMetadata(ctx context.Context) sdk.Metadata {
	return sdk.Metadata{
		User:       GetUser(ctx),
		ProjectIds: GetProjects(ctx),
	}
}

func AddMetadata(ctx context.Context, metadata sdk.Metadata) context.Context {
	return context.WithValue(context.WithValue(ctx, projects, metadata.ProjectIds), userValue, metadata.User)
}
