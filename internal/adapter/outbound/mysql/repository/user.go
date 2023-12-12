package repository

import (
	"errors"
	"todo_api/internal/adapter/outbound/mysql/model"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Get(id domain.UserIdentifier) (*domain.User, apperr.AppErr) {
	var row *model.User
	if err := r.db.Preload("Company").First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	user, aerr := model.MarshalUser(row)
	if aerr != nil {
		return nil, aerr
	}
	return user, nil
}

func (r *UserRepository) Create(user *domain.User) (*domain.UserIdentifier, apperr.AppErr) {
	row := model.UnmarshalUser(user)
	if err := r.db.Create(&row).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	id := domain.UserIdentifier(row.ID)
	return &id, nil
}

func (r *UserRepository) Update(user *domain.User) apperr.AppErr {
	row := model.UnmarshalUser(user)
	if err := r.db.Save(&row).Error; err != nil {
		return apperr.NewInternalServerError().Wrap(err)
	}
	return nil
}
