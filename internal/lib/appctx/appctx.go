package appctx

import (
	"context"
	"fmt"
	"todo_api/internal/domain/model"
)

type contextKey string

const userKey contextKey = "user"

func SetUser(parents context.Context, user model.User) context.Context {
	return context.WithValue(parents, userKey, user)
}

func GetUser(ctx context.Context) (*model.User, error) {
	v := ctx.Value(userKey)

	user, ok := v.(model.User)
	if !ok {
		return nil, fmt.Errorf("user is not set")
	}

	return &user, nil
}
