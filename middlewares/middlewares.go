package middlewares

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

func GetProjects(ctx context.Context) []string {
	val, ok := ctx.Value(sdk.ProjectsTypeVal).([]string)
	if !ok {
		return nil
	}
	return val
}

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

func GetMetadata(ctx context.Context) sdk.Metadata {
	return sdk.Metadata{
		User:       GetUser(ctx),
		ProjectIds: GetProjects(ctx),
	}
}

func AddMetadata(ctx context.Context, metadata sdk.Metadata) context.Context {
	return context.WithValue(context.WithValue(ctx, sdk.ProjectsTypeVal, metadata.ProjectIds), sdk.UserTypeVal, metadata.User)
}
