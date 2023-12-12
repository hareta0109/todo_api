package repository

import (
	"errors"
	"todo_api/internal/adapter/outbound/mysql/model"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) Get(id domain.UserIdentifier) (*domain.Auth, apperr.AppErr) {
	var row *model.Auth
	if err := r.db.Preload("Company").First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	auth, aerr := model.MarshalAuth(row)
	if aerr != nil {
		return nil, aerr
	}
	return auth, nil
}

func (r *AuthRepository) Create(auth *domain.Auth) (*domain.UserIdentifier, apperr.AppErr) {
	row := model.UnmarshalAuth(auth)
	if err := r.db.Create(&row).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	id := domain.UserIdentifier(row.ID)
	return &id, nil
}

func (r *AuthRepository) Update(auth *domain.Auth) apperr.AppErr {
	row := model.UnmarshalAuth(auth)
	if err := r.db.Save(&row).Error; err != nil {
		return apperr.NewInternalServerError().Wrap(err)
	}
	return nil
}
