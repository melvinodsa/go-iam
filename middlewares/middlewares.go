package middlewares

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

func GetProjects(ctx context.Context) []string {
	return ctx.Value("projects").([]string)
}

func GetUser(ctx context.Context) *sdk.User {
	user := ctx.Value("user")
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
