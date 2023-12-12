package repository

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type AuthRepository interface {
	Get(model.UserIdentifier) (*model.Auth, apperr.AppErr)
	Create(*model.Auth) (*model.UserIdentifier, apperr.AppErr)
	Update(*model.Auth) apperr.AppErr
}
