package context

import (
	"context"

	"someblocks/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithCurrentUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func CurrentUser(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
