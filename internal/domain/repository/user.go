package repository

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type UserRepository interface {
	Get(model.UserIdentifier) (*model.User, apperr.AppErr)

	Update(*model.User) apperr.AppErr
}
